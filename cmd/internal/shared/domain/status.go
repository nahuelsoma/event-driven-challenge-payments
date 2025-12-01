package domain

// Status represents the status of a payment
type Status string

const (
	StatusPending   Status = "PENDING"   // The payment is pending
	StatusReserved  Status = "RESERVED"  // The payment is reserved
	StatusCompleted Status = "COMPLETED" // The payment is completed
	StatusFailed    Status = "FAILED"    // The payment is failed
)
