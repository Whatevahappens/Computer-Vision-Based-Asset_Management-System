package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

var CVServiceURL string

func SetCVServiceURL(url string) {
	CVServiceURL = url
}

func StartAudit(locationID, notes, userID string) (*model.AuditSession, error) {
	session := &model.AuditSession{
		ID:          uuid.New().String(),
		StartedAt:   time.Now(),
		Status:      model.InProgress,
		Notes:       notes,
		LocationID:  locationID,
		PerformedBy: userID,
	}
	if err := repository.CreateAuditSession(session); err != nil {
		return nil, err
	}
	return session, nil
}

var classMap = map[string]string{
	"chair":     "Сандал",
	"monitor":   "Дэлгэц",
	"table":     "Ширээ",
	"processor": "Процессор",
}

func mapClassName(yoloClass string) string {
	if mapped, ok := classMap[yoloClass]; ok {
		return mapped
	}
	return yoloClass
}

func getBaseAssetName(name string) string {
	if idx := strings.LastIndex(name, " #"); idx > 0 {
		return name[:idx]
	}
	return name
}

type ImageInput struct {
	Data     []byte
	Filename string
}

func RunCVAudit(sessionID string, images []ImageInput, userID string) (*model.AuditSession, error) {
	session, err := repository.FindAuditSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("audit session not found")
	}

	maxDetected := make(map[string]int)                 // category → max count across corners
	allDetections := make(map[string][]dto.CVDetection) // for detailed findings
	var lastImagePath string
	var modelName, modelVer string

	for _, img := range images {
		detections, err := callCVService(img.Data, img.Filename)
		if err != nil {
			return nil, fmt.Errorf("CV service error: %v", err)
		}
		lastImagePath = detections.ImagePath
		modelName = detections.ModelName
		modelVer = detections.ModelVer

		cornerCounts := make(map[string]int)
		for _, d := range detections.Detections {
			mapped := mapClassName(d.ClassName)
			cornerCounts[mapped]++
			allDetections[mapped] = append(allDetections[mapped], d)
		}

		for cat, count := range cornerCounts {
			if count > maxDetected[cat] {
				maxDetected[cat] = count
			}
		}
	}

	registeredAssets, err := repository.ListAssetsByLocation(session.LocationID)
	if err != nil {
		return nil, err
	}

	registeredCounts := make(map[string]int)
	registeredByCategory := make(map[string][]model.Asset)
	for _, a := range registeredAssets {
		base := getBaseAssetName(a.AssetName)
		registeredCounts[base]++
		registeredByCategory[base] = append(registeredByCategory[base], a)
	}

	for category, maxCount := range maxDetected {
		regCount := registeredCounts[category]
		dets := allDetections[category]

		sort.Slice(dets, func(i, j int) bool {
			return dets[i].Confidence > dets[j].Confidence
		})
		if len(dets) > maxCount {
			dets = dets[:maxCount]
		}

		for i, det := range dets {
			findingType := model.Matched
			if i >= regCount {
				findingType = model.Unregistered
			}

			finding := &model.AuditFinding{
				ID:             uuid.New().String(),
				Type:           findingType,
				Confidence:     det.Confidence,
				Notes:          fmt.Sprintf("Detected %s (confidence: %.2f)", category, det.Confidence),
				AuditSessionID: sessionID,
			}
			if findingType == model.Matched {
				assets := registeredByCategory[category]
				if i < len(assets) {
					assetID := assets[i].ID
					finding.DetectedAssetID = &assetID
				}
			}
			repository.CreateAuditFinding(finding)

			evidence := &model.AuditEvidence{
				ID:             uuid.New().String(),
				FilePath:       lastImagePath,
				CapturedAt:     time.Now(),
				ModelName:      modelName,
				ModelVersion:   modelVer,
				AuditFindingID: finding.ID,
			}
			repository.CreateAuditEvidence(evidence)
		}
	}

	for category, regCount := range registeredCounts {
		detCount := maxDetected[category]
		if detCount < regCount {
			missingCount := regCount - detCount
			assets := registeredByCategory[category]
			for i := regCount - missingCount; i < regCount; i++ {
				assetID := assets[i].ID
				finding := &model.AuditFinding{
					ID:              uuid.New().String(),
					Type:            model.Missing,
					Confidence:      0,
					Notes:           fmt.Sprintf("'%s' олдсонгүй (%d бүртгэлтэйгээс %d илэрсэн)", category, regCount, detCount),
					AuditSessionID:  sessionID,
					ExpectedAssetID: &assetID,
				}
				repository.CreateAuditFinding(finding)
			}
		}
	}

	allCategories := make(map[string]bool)
	for k := range registeredCounts {
		allCategories[k] = true
	}
	for k := range maxDetected {
		allCategories[k] = true
	}

	for cat := range allCategories {
		reg := registeredCounts[cat]
		det := maxDetected[cat]
		summary := &model.AuditSummary{
			ID:              uuid.New().String(),
			Category:        model.AssetCategory(cat),
			RegisteredCount: reg,
			DetectedCount:   det,
			Difference:      det - reg,
			AuditSessionID:  sessionID,
		}
		repository.CreateAuditSummary(summary)
	}

	now := time.Now()
	session.EndedAt = &now
	session.Status = model.Completed
	repository.UpdateAuditSession(session)

	return repository.FindAuditSessionByID(sessionID)
}

func callCVService(imageData []byte, filename string) (*dto.CVDetectionResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	contentType := "application/octet-stream"
	switch {
	case strings.HasSuffix(strings.ToLower(filename), ".jpg"), strings.HasSuffix(strings.ToLower(filename), ".jpeg"):
		contentType = "image/jpeg"
	case strings.HasSuffix(strings.ToLower(filename), ".png"):
		contentType = "image/png"
	case strings.HasSuffix(strings.ToLower(filename), ".webp"):
		contentType = "image/webp"
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", contentType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	part.Write(imageData)
	writer.Close()

	resp, err := http.Post(CVServiceURL+"/detect", writer.FormDataContentType(), body)
	if err != nil {
		return nil, fmt.Errorf("cannot reach CV service: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("CV service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result dto.CVDetectionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse CV response: %v", err)
	}
	return &result, nil
}
