package entity

import "encoding/json"

type (
	CacheInitOrderRequest struct {
		UserID   int64  `json:"user_id"`
		EventID  int64  `json:"event_id"`
		Date     string `json:"date"`
		Quantity int64  `json:"quantity"`
	}

	CacheInitOrderResponse struct {
		Date     string `json:"date"`
		EventID  int64  `json:"event_id"`
		Quantity int64  `json:"quantity"`
	}
)

func (u CacheInitOrderRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u CacheInitOrderRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}
