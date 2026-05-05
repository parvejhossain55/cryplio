package dto

import "time"

type DisputeResponse struct {
	DisputeID      string     `json:"dispute_id"`
	TradeID        string     `json:"trade_id"`
	RaisedBy       string     `json:"raised_by"`
	ReasonCode     string     `json:"reason_code"`
	ReasonText     *string    `json:"reason_text,omitempty"`
	EvidenceLinks  []string   `json:"evidence_links"`
	Status         string     `json:"status"`
	AssignedAdmin  *string    `json:"assigned_admin,omitempty"`
	AssignedAt     *time.Time `json:"assigned_at,omitempty"`
	ResolutionType *string    `json:"resolution_type,omitempty"`
	ResolutionNote *string    `json:"resolution_note,omitempty"`
	ResolvedBy     *string    `json:"resolved_by,omitempty"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type RaiseDisputeRequest struct {
	ReasonCode string `json:"reason_code" binding:"required"`
	ReasonText string `json:"reason_text"`
}

type ResolveDisputeRequest struct {
	Resolution string `json:"resolution" binding:"required,oneof=release_to_buyer return_to_seller cancel"`
	Note       string `json:"note" binding:"required"`
}
