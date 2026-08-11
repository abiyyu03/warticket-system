package ticket

type (
	PurchaseRequest struct {
		EventID int64 `json:"event_id"`
	}

	PurchaseResponse struct {
		Status string `json:"status"` //
	}
)
