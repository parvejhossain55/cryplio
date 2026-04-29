package dispute

import (
	"time"

	"github.com/google/uuid"
)

// DisputeStatus represents dispute lifecycle state
type DisputeStatus string

const (
	DisputeStatusPending     DisputeStatus = "pending"
	DisputeStatusAssigned    DisputeStatus = "assigned"
	DisputeStatusUnderReview DisputeStatus = "under_review"
	DisputeStatusResolved    DisputeStatus = "resolved"
	DisputeStatusAppealed    DisputeStatus = "appealed"
	DisputeStatusClosed      DisputeStatus = "closed"
)

// DisputeResolution represents the outcome of a dispute
type DisputeResolution string

const (
	DisputeResolutionReleaseToBuyer DisputeResolution = "release_to_buyer"
	DisputeResolutionReturnToSeller DisputeResolution = "return_to_seller"
	DisputeResolutionPartialSplit   DisputeResolution = "partial_split"
	DisputeResolutionCancel         DisputeResolution = "cancel"
)

// Dispute represents a trade dispute
type Dispute struct {
	DisputeID     uuid.UUID          `db:"dispute_id" json:"dispute_id"`
	TradeID       uuid.UUID          `db:"trade_id" json:"trade_id"`
	RaisedBy      uuid.UUID          `db:"raised_by" json:"raised_by"` // User who raised the dispute
	Reason        *string            `db:"reason" json:"reason,omitempty"`
	EvidenceLinks []string           `db:"evidence_links" json:"evidence_links"` // JSONB: URLs to evidence
	Status        DisputeStatus      `db:"status" json:"status"`
	Resolution    *DisputeResolution `db:"resolution" json:"resolution,omitempty"`
	ResolvedBy    *uuid.UUID         `db:"resolved_by" json:"resolved_by,omitempty"` // Admin ID
	ResolvedAt    *time.Time         `db:"resolved_at" json:"resolved_at,omitempty"`
	CreatedAt     time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `db:"updated_at" json:"updated_at"`
}

// IsOpen checks if the dispute is still open
func (d *Dispute) IsOpen() bool {
	return d.Status == DisputeStatusPending || d.Status == DisputeStatusAssigned || d.Status == DisputeStatusUnderReview
}

// IsResolved checks if the dispute has been resolved
func (d *Dispute) IsResolved() bool {
	return d.Status == DisputeStatusResolved || d.Status == DisputeStatusClosed
}

// Assign assigns the dispute to an admin
func (d *Dispute) Assign(adminID uuid.UUID) {
	d.Status = DisputeStatusAssigned
	// Could store assigned_to field if needed
	_ = adminID
}

// Resolve resolves the dispute with a decision
func (d *Dispute) Resolve(adminID uuid.UUID, resolution DisputeResolution) {
	d.Status = DisputeStatusResolved
	d.Resolution = &resolution
	now := time.Now()
	d.ResolvedBy = &adminID
	d.ResolvedAt = &now
}

// Close closes the dispute (after resolution or appeal)
func (d *Dispute) Close() {
	d.Status = DisputeStatusClosed
}

// DisputeMessage represents a message in dispute chat
type DisputeMessage struct {
	MessageID   uuid.UUID `db:"message_id" json:"message_id"`
	DisputeID   uuid.UUID `db:"dispute_id" json:"dispute_id"`
	SenderID    uuid.UUID `db:"sender_id" json:"sender_id"`
	MessageType string    `db:"message_type" json:"message_type"` // text, image, file
	Content     *string   `db:"content" json:"content,omitempty"`
	FileURL     *string   `db:"file_url" json:"file_url,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// IsText checks if the message is text
func (m *DisputeMessage) IsText() bool {
	return m.MessageType == "text"
}

// IsAttachment checks if the message contains a file attachment
func (m *DisputeMessage) IsAttachment() bool {
	return m.MessageType == "image" || m.MessageType == "file"
}
