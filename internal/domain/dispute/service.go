package dispute

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	CreateDispute(ctx context.Context, d *Dispute) error
	GetDispute(ctx context.Context, id uuid.UUID) (*Dispute, error)
	AssignDispute(ctx context.Context, id uuid.UUID, adminID uuid.UUID) error
	ResolveDispute(ctx context.Context, id uuid.UUID, adminID uuid.UUID, resolution DisputeResolution, note string) error
	UploadEvidence(ctx context.Context, disputeID, userID uuid.UUID, evidenceURL string) error
	ListDisputes(ctx context.Context, status string, limit, offset int) ([]*Dispute, int, error)
	CountDisputes(ctx context.Context, status string) (int, error)
}

type disputeService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &disputeService{repo: repo}
}

func (s *disputeService) CreateDispute(ctx context.Context, d *Dispute) error {
	if d.TradeID == uuid.Nil {
		return errors.New("trade_id is required")
	}
	if d.RaisedBy == uuid.Nil {
		return errors.New("raised_by is required")
	}
	if d.ReasonCode == "" {
		return errors.New("reason_code is required")
	}

	if d.DisputeID == uuid.Nil {
		d.DisputeID = uuid.New()
	}
	d.Status = DisputeStatusPending

	return s.repo.Create(ctx, d)
}

func (s *disputeService) GetDispute(ctx context.Context, id uuid.UUID) (*Dispute, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *disputeService) AssignDispute(ctx context.Context, id uuid.UUID, adminID uuid.UUID) error {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("dispute not found")
	}

	d.Assign(adminID)
	return s.repo.Update(ctx, d)
}

func (s *disputeService) ResolveDispute(ctx context.Context, id uuid.UUID, adminID uuid.UUID, resolution DisputeResolution, note string) error {
	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("dispute not found")
	}

	if !d.IsOpen() {
		return fmt.Errorf("dispute is already in status: %s", d.Status)
	}

	d.Resolve(adminID, resolution, note)
	return s.repo.Update(ctx, d)
}

func (s *disputeService) ListDisputes(ctx context.Context, status string, limit, offset int) ([]*Dispute, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, status, limit, offset)
}

func (s *disputeService) CountDisputes(ctx context.Context, status string) (int, error) {
	return s.repo.CountDisputes(ctx, status)
}

func (s *disputeService) UploadEvidence(ctx context.Context, disputeID, userID uuid.UUID, evidenceURL string) error {
	d, err := s.repo.GetByID(ctx, disputeID)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("dispute not found")
	}
	if !d.IsOpen() {
		return fmt.Errorf("dispute is already in status: %s", d.Status)
	}
	// Only the dispute raiser or the other party can upload evidence
	// For simplicity, allow any participant (raiser or trade party)
	d.EvidenceLinks = append(d.EvidenceLinks, evidenceURL)
	return s.repo.Update(ctx, d)
}
