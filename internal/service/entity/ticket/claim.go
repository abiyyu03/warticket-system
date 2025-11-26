package ticket

type (
	ClaimTicketRequest struct {
		ChairCode  *int   `json:"chair_code"` // book for vip only
		Date       string `json:"date"`
		TicketType string `json:"ticket_type"`
	}

	ClaimTicketResponse struct{}
)
