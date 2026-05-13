# YOLO Models Directory

Place your trained model weights here:

- `asset_detector.pt` — Custom fine-tuned model for asset detection (preferred)
- `yolo26m.pt` — YOLO26-M pre-trained weights
- `yolo11m.pt` — YOLO11-M pre-trained weights

If no model is found, the service will auto-download `yolo11n.pt` from Ultralytics.

## Training Your Own Model

```bash
# Using Google Colab or local GPU
from ultralytics import YOLO

model = YOLO("yolo11n.pt")  # or yolo26n.pt when available
results = model.train(
    data="path/to/your/dataset.yaml",
    epochs=100,
    imgsz=640,
    batch=16,
    name="asset_detector"
)
```

After training, copy `runs/detect/asset_detector/weights/best.pt` to this directory
as `asset_detector.pt`.
