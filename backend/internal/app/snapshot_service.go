package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type snapshotService struct {
	repo SnapshotRepository
}

func NewSnapshotService(repo SnapshotRepository) SnapshotService {
	return &snapshotService{repo: repo}
}

func (s *snapshotService) GetSnapshot(ctx context.Context, id uuid.UUID) (*domain.VersionSnapshot, error) {
	return s.repo.Get(ctx, id)
}

func (s *snapshotService) ListSnapshots(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.VersionSnapshot, error) {
	return s.repo.List(ctx, roadmapItemID)
}

func (s *snapshotService) ListSnapshotsByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VersionSnapshot, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *snapshotService) CreateSnapshot(ctx context.Context, roadmapItemID uuid.UUID, data map[string]interface{}, createdBy uuid.UUID) (*domain.VersionSnapshot, error) {
	// Calculate Hash
	hash, err := calculateHash(data)
	if err != nil {
		return nil, err
	}

	snap := &domain.VersionSnapshot{
		ID:            uuid.New(),
		RoadmapItemID: roadmapItemID,
		SnapshotData:  data,
		Hash:          hash,
		CreatedBy:     createdBy,
	}
	if err := s.repo.Create(ctx, snap); err != nil {
		return nil, err
	}
	return snap, nil
}

func calculateHash(data map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:]), nil
}
