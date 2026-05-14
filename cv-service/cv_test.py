"""
UC-002: Компьютерийн хараанд суурилсан аудит явуулах
Test cases: TC-006, TC-007, TC-010
"""
import pytest
import numpy as np
import cv2
from fastapi.testclient import TestClient

from main import app, preprocess_image

client = TestClient(app)


# ============================================================
# Fixtures
# ============================================================

@pytest.fixture
def sample_image_bytes():
    """Generate a valid JPEG image for testing."""
    img = np.random.randint(50, 200, (480, 640, 3), dtype=np.uint8)
    _, buf = cv2.imencode(".jpg", img)
    return buf.tobytes()


@pytest.fixture
def dark_image_bytes():
    """Generate a dark JPEG image (avg brightness ~30)."""
    img = np.full((480, 640, 3), 30, dtype=np.uint8)
    # Add some noise so it's not perfectly uniform
    noise = np.random.randint(0, 15, (480, 640, 3), dtype=np.uint8)
    img = cv2.add(img, noise)
    _, buf = cv2.imencode(".jpg", img)
    return buf.tobytes()


# ============================================================
# TC-007: CLAHE зургийн боловсруулалт
# ============================================================

def test_preprocess_dark_image_clahe(dark_image_bytes):
    """TC-007: CLAHE should increase brightness of dark images."""
    dark_img = cv2.imdecode(
        np.frombuffer(dark_image_bytes, np.uint8),
        cv2.IMREAD_COLOR
    )
    avg_before = dark_img.mean()

    result = preprocess_image(dark_image_bytes)
    avg_after = result.mean()

    assert avg_after > avg_before, (
        f"CLAHE should increase brightness: "
        f"{avg_before:.1f} -> {avg_after:.1f}"
    )


def test_preprocess_normal_image_no_crash(sample_image_bytes):
    """TC-007: Normal image should process without error."""
    result = preprocess_image(sample_image_bytes)
    assert result is not None
    assert result.shape[2] == 3  # BGR channels preserved


def test_preprocess_output_shape(sample_image_bytes):
    """TC-007: Output shape should match input shape."""
    original = cv2.imdecode(
        np.frombuffer(sample_image_bytes, np.uint8),
        cv2.IMREAD_COLOR
    )
    result = preprocess_image(sample_image_bytes)
    assert result.shape == original.shape


# ============================================================
# TC-006: Объект илрүүлэлтийн бүтэц
# ============================================================

def test_detect_returns_valid_json(sample_image_bytes):
    """TC-006: /detect endpoint returns valid JSON structure."""
    response = client.post(
        "/detect",
        files={"file": ("test.jpg", sample_image_bytes, "image/jpeg")},
        params={"confidence": 0.25},
    )
    assert response.status_code == 200
    data = response.json()

    # Required top-level fields
    assert "detections" in data
    assert "inference_ms" in data
    assert "model_name" in data
    assert "image_width" in data
    assert "image_height" in data

    assert isinstance(data["detections"], list)
    assert isinstance(data["inference_ms"], (int, float))
    assert data["inference_ms"] >= 0


def test_detect_structure_per_detection(sample_image_bytes):
    """TC-006: Each detection has class_name, confidence, box."""
    response = client.post(
        "/detect",
        files={"file": ("test.jpg", sample_image_bytes, "image/jpeg")},
        params={"confidence": 0.1},  # low threshold to get detections
    )
    data = response.json()

    for det in data["detections"]:
        assert "class_name" in det, "missing class_name"
        assert "confidence" in det, "missing confidence"
        assert "box" in det, "missing box"
        assert 0.0 <= det["confidence"] <= 1.0, (
            f"confidence out of range: {det['confidence']}"
        )
        assert len(det["box"]) == 4, (
            f"box should have 4 coords, got {len(det['box'])}"
        )


def test_detect_png_image():
    """TC-006: PNG format should also be accepted."""
    img = np.random.randint(50, 200, (480, 640, 3), dtype=np.uint8)
    _, buf = cv2.imencode(".png", img)

    response = client.post(
        "/detect",
        files={"file": ("test.png", buf.tobytes(), "image/png")},
        params={"confidence": 0.25},
    )
    assert response.status_code == 200
    assert "detections" in response.json()


def test_detect_confidence_threshold(sample_image_bytes):
    """TC-006: Higher confidence threshold should return
    fewer or equal detections."""
    resp_low = client.post(
        "/detect",
        files={"file": ("test.jpg", sample_image_bytes, "image/jpeg")},
        params={"confidence": 0.1},
    )
    resp_high = client.post(
        "/detect",
        files={"file": ("test.jpg", sample_image_bytes, "image/jpeg")},
        params={"confidence": 0.8},
    )
    low_count = len(resp_low.json()["detections"])
    high_count = len(resp_high.json()["detections"])
    assert high_count <= low_count, (
        f"higher threshold should give fewer detections: "
        f"low={low_count}, high={high_count}"
    )


# ============================================================
# TC-010: Алдаатай оролт
# ============================================================

def test_preprocess_invalid_image_raises():
    """TC-010: Corrupt bytes should raise ValueError."""
    with pytest.raises(ValueError, match="Failed to decode"):
        preprocess_image(b"not_an_image_bytes_at_all")


def test_detect_empty_file_rejected():
    """TC-010: Empty file upload should return error."""
    response = client.post(
        "/detect",
        files={"file": ("empty.jpg", b"", "image/jpeg")},
    )
    assert response.status_code in (400, 422)


def test_detect_no_file_rejected():
    """TC-010: Request without file should return 422."""
    response = client.post("/detect")
    assert response.status_code == 422


# ============================================================
# Model info endpoint
# ============================================================

def test_model_info():
    """Model info endpoint returns model metadata."""
    response = client.get("/model-info")
    assert response.status_code == 200
    data = response.json()
    assert "model_name" in data
    assert "num_classes" in data


def test_health():
    """Health endpoint returns ok."""
    response = client.get("/health")
    assert response.status_code == 200