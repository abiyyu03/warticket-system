package entity

type (
	DecrTicketQuotaRequest struct {
		EventID  int64 `json:"event_id"`
		Quantity int64 `json:"quantity"`
	}
)
