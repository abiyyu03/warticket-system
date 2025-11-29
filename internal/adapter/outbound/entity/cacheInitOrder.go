package entity

import "encoding/json"

type (
	CacheInitOrderRequest struct {
		UserID int64 `json:"user_id"`
	}

	CacheInitOrderResponse struct {
		ChairCode []string `json:"chair_code"`
		Date      string   `json:"date"`
		EventID   int64    `json:"event_id"`
		UserID    int64    `json:"user_id,omitempty"`
	}
)

func (u CacheInitOrderRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u CacheInitOrderRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}
