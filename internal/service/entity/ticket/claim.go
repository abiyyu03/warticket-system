package ticket

type (
	PurchaseRequest struct {
		TxID string `json:"tx_id"`
	}

	PurchaseResponse struct {
		Status string `json:"status"` //
	}
)
