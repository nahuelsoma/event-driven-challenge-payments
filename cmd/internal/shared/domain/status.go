package domain

// Status represents the status of a payment
type Status string

const (
	StatusPending   Status = "pending"   // The payment is pending
	StatusReserved  Status = "reserved"  // The payment is reserved
	StatusCompleted Status = "completed" // The payment is completed
	StatusFailed    Status = "failed"    // The payment is failed
)
