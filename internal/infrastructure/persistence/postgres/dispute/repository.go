package dispute

import (
	"context"
	"database/sql"
	"fmt"

	"cryplio/internal/domain/dispute"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type disputeRepository struct {
	db *sql.DB
}

func NewDisputeRepository(db *sql.DB) dispute.Repository {
	return &disputeRepository{db: db}
}

func (r *disputeRepository) Create(ctx context.Context, d *dispute.Dispute) error {
	query := `
		INSERT INTO disputes (
			dispute_id, trade_id, raised_by, reason_code, reason_text,
			evidence_links, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
		) RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		d.DisputeID, d.TradeID, d.RaisedBy, d.ReasonCode, d.ReasonText,
		pq.Array(d.EvidenceLinks), d.Status,
	).Scan(&d.CreatedAt, &d.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create dispute: %w", err)
	}
	return nil
}

func (r *disputeRepository) Update(ctx context.Context, d *dispute.Dispute) error {
	query := `
		UPDATE disputes
		SET status = $1, assigned_admin = $2, assigned_at = $3,
		    resolution_type = $4, resolution_note = $5,
		    resolved_by = $6, resolved_at = $7, closed_at = $8,
		    updated_at = NOW()
		WHERE dispute_id = $9
	`
	_, err := r.db.ExecContext(
		ctx, query,
		d.Status, d.AssignedAdmin, d.AssignedAt,
		d.ResolutionType, d.ResolutionNote,
		d.ResolvedBy, d.ResolvedAt, nil, // TODO: handle closed_at if needed
		d.DisputeID,
	)
	if err != nil {
		return fmt.Errorf("update dispute: %w", err)
	}
	return nil
}

func (r *disputeRepository) GetByID(ctx context.Context, id uuid.UUID) (*dispute.Dispute, error) {
	query := `
		SELECT dispute_id, trade_id, raised_by, reason_code, reason_text,
		       evidence_links, status, assigned_admin, assigned_at,
		       resolution_type, resolution_note, resolved_by, resolved_at,
		       created_at, updated_at
		FROM disputes
		WHERE dispute_id = $1
	`
	var d dispute.Dispute
	var evidenceLinks []string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.DisputeID, &d.TradeID, &d.RaisedBy, &d.ReasonCode, &d.ReasonText,
		pq.Array(&evidenceLinks), &d.Status, &d.AssignedAdmin, &d.AssignedAt,
		&d.ResolutionType, &d.ResolutionNote, &d.ResolvedBy, &d.ResolvedAt,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get dispute: %w", err)
	}
	d.EvidenceLinks = evidenceLinks
	return &d, nil
}
