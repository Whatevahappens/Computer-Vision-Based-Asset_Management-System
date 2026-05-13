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

func RunCVAudit(sessionID string, imageData []byte, filename string, userID string) (*model.AuditSession, error) {
	session, err := repository.FindAuditSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("audit session not found")
	}
	detections, err := callCVService(imageData, filename)
	if err != nil {
		return nil, fmt.Errorf("CV service error: %v", err)
	}

	registeredAssets, err := repository.ListAssetsByLocation(session.LocationID)
	if err != nil {
		return nil, err
	}

	detectedCounts := make(map[string]int)
	for _, d := range detections.Detections {
		detectedCounts[d.ClassName]++
	}

	registeredCounts := make(map[string]int)
	registeredMap := make(map[string]string)
	for _, a := range registeredAssets {
		key := a.AssetName
		registeredCounts[key]++
		registeredMap[key] = a.ID
	}

	for _, det := range detections.Detections {
		findingType := model.Matched
		if registeredCounts[det.ClassName] <= 0 {
			findingType = model.Unregistered
		}
		finding := &model.AuditFinding{
			ID:             uuid.New().String(),
			Type:           findingType,
			Confidence:     det.Confidence,
			Notes:          fmt.Sprintf("Detected %s (confidence: %.2f)", det.ClassName, det.Confidence),
			AuditSessionID: sessionID,
		}
		if assetID, ok := registeredMap[det.ClassName]; ok {
			finding.DetectedAssetID = &assetID
		}
		repository.CreateAuditFinding(finding)

		evidence := &model.AuditEvidence{
			ID:             uuid.New().String(),
			FilePath:       detections.ImagePath,
			CapturedAt:     time.Now(),
			ModelName:      detections.ModelName,
			ModelVersion:   detections.ModelVer,
			AuditFindingID: finding.ID,
		}
		repository.CreateAuditEvidence(evidence)
	}

	for _, a := range registeredAssets {
		if detectedCounts[a.AssetName] <= 0 {
			assetID := a.ID
			finding := &model.AuditFinding{
				ID:              uuid.New().String(),
				Type:            model.Missing,
				Confidence:      0,
				Notes:           fmt.Sprintf("Registered asset '%s' not detected", a.AssetName),
				AuditSessionID:  sessionID,
				ExpectedAssetID: &assetID,
			}
			repository.CreateAuditFinding(finding)
		}
	}

	allCategories := make(map[string]bool)
	for k := range registeredCounts {
		allCategories[k] = true
	}
	for k := range detectedCounts {
		allCategories[k] = true
	}
	for cat := range allCategories {
		reg := registeredCounts[cat]
		det := detectedCounts[cat]
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
	part, err := writer.CreateFormFile("file", filename)
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
