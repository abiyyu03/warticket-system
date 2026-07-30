package entity

type (
	SetTicketQuotaRequest struct {
		EventID int64 `json:"event_id"`
		Quota   int64 `json:"quota"`
	}
)
