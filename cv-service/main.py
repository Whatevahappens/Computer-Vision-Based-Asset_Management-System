"""
Computer Vision Audit Service
FastAPI + YOLO26 (Ultralytics) object detection endpoint
"""

import os
import time
import uuid
from pathlib import Path

import cv2
import numpy as np
from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

app = FastAPI(
    title="CV Audit Service",
    description="YOLO26-based object detection for asset auditing",
    version="1.0.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# ── Directories ──
UPLOAD_DIR = Path("uploads")
MODEL_DIR = Path("models")
UPLOAD_DIR.mkdir(exist_ok=True)
MODEL_DIR.mkdir(exist_ok=True)

# YOLO Model hierarchy
# tries YOLO26 first, falls back to YOLO11, then YOLOv8
_model = None
_model_name = "unknown"
_model_version = "unknown"


def get_model():
    """Lazy-load the best available YOLO model."""
    global _model, _model_name, _model_version

    if _model is not None:
        return _model

    from ultralytics import YOLO

    # Priority order: custom fine-tuned → YOLO26 → YOLO11 → YOLOv8
    candidates = [
        ("models/asset_detector.pt", "YOLO26-custom", "1.0"),
        ("models/yolo26m.pt", "YOLO26-M", "26.0"),
        ("models/yolo11m.pt", "YOLO11-M", "11.0"),
        ("yolo11n.pt", "YOLO11-N", "11.0"),          
        ("yolov8n.pt", "YOLOv8-N", "8.0"),            
    ]

    for path, name, ver in candidates:
        try:
            _model = YOLO(path)
            _model_name = name
            _model_version = ver
            print(f"✓ Loaded model: {name} ({path})")
            return _model
        except Exception as e:
            print(f"  Skipped {path}: {e}")
            continue

    raise RuntimeError("No YOLO model available. Place a .pt file in ./models/")


# ── Response schemas ──
class Detection(BaseModel):
    class_name: str
    confidence: float
    box: list[int]  # [x1, y1, x2, y2]


class DetectionResponse(BaseModel):
    detections: list[Detection]
    image_path: str
    model_name: str
    model_version: str
    inference_ms: float
    image_width: int
    image_height: int


# ── Image preprocessing (OpenCV) ──
def preprocess_image(image_bytes: bytes) -> np.ndarray:
    """
    Decode, resize to 640×640, apply CLAHE for low-light enhancement.
    """
    nparr = np.frombuffer(image_bytes, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    if img is None:
        raise ValueError("Failed to decode image")

    # CLAHE enhancement for low-light conditions
    lab = cv2.cvtColor(img, cv2.COLOR_BGR2LAB)
    l_channel, a, b = cv2.split(lab)
    clahe = cv2.createCLAHE(clipLimit=2.0, tileGridSize=(8, 8))
    l_channel = clahe.apply(l_channel)
    enhanced = cv2.merge([l_channel, a, b])
    img = cv2.cvtColor(enhanced, cv2.COLOR_LAB2BGR)

    return img


# ── Endpoints ──
@app.get("/health")
def health():
    return {"status": "ok", "model_loaded": _model is not None}


@app.post("/detect", response_model=DetectionResponse)
async def detect_objects(
    file: UploadFile = File(...),
    confidence: float = 0.25,
    iou: float = 0.45,
):
    """
    Upload an image and run YOLO object detection.

    - **file**: Image file (JPEG/PNG)
    - **confidence**: Minimum confidence threshold (default 0.25)
    - **iou**: IoU threshold for NMS (default 0.45)
    """
    # Validate file type
    if file.content_type not in ("image/jpeg", "image/png", "image/webp"):
        raise HTTPException(400, "Only JPEG, PNG, WebP images are accepted")

    # Read and preprocess
    image_bytes = await file.read()
    try:
        img = preprocess_image(image_bytes)
    except ValueError as e:
        raise HTTPException(400, str(e))

    h, w = img.shape[:2]

    # Save uploaded image
    ext = file.filename.rsplit(".", 1)[-1] if "." in file.filename else "jpg"
    saved_name = f"{uuid.uuid4().hex}.{ext}"
    saved_path = UPLOAD_DIR / saved_name
    cv2.imwrite(str(saved_path), img)

    # Run inference
    model = get_model()
    start = time.time()
    results = model.predict(
        source=img,
        conf=confidence,
        iou=iou,
        imgsz=640,
        verbose=False,
    )
    inference_ms = (time.time() - start) * 1000

    # Parse results
    detections: list[Detection] = []
    for result in results:
        boxes = result.boxes
        if boxes is None:
            continue
        for i in range(len(boxes)):
            xyxy = boxes.xyxy[i].cpu().numpy().astype(int).tolist()
            conf = float(boxes.conf[i].cpu().numpy())
            cls_id = int(boxes.cls[i].cpu().numpy())
            cls_name = model.names.get(cls_id, f"class_{cls_id}")
            detections.append(Detection(
                class_name=cls_name,
                confidence=round(conf, 4),
                box=xyxy,
            ))

    # Save annotated image
    annotated = results[0].plot() if results else img
    annotated_name = f"annotated_{saved_name}"
    annotated_path = UPLOAD_DIR / annotated_name
    cv2.imwrite(str(annotated_path), annotated)

    return DetectionResponse(
        detections=detections,
        image_path=str(annotated_path),
        model_name=_model_name,
        model_version=_model_version,
        inference_ms=round(inference_ms, 2),
        image_width=w,
        image_height=h,
    )


@app.post("/detect/batch")
async def detect_batch(files: list[UploadFile] = File(...)):
    """Detect objects in multiple images at once."""
    results = []
    for f in files:
        resp = await detect_objects(f)
        results.append(resp)
    return {"results": results, "count": len(results)}


@app.get("/model-info")
def model_info():
    """Return information about the loaded model."""
    model = get_model()
    return {
        "model_name": _model_name,
        "model_version": _model_version,
        "classes": model.names,
        "num_classes": len(model.names),
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
