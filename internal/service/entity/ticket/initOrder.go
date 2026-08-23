package ticket

import "go-projects/hexagonal-example/internal/adapter/outbound/entity"

type (
	InitOrderRequest struct {
		EventID  int64 `json:"event_id"`
		Quantity int64 `json:"quantity"`
	}

	InitOrderResponse struct {
		TxID     string `json:"tx_id"`
		Date     string `json:"date"`
		EventID  int64  `json:"event_id"`
		Quantity int64  `json:"quantity"`
		Price    int64  `json:"price"`
	}
)

// ToObSetCache menyusun payload cache reservasi. tx_id & date diisi server-side
// (date = tanggal event), harga dari event supaya tidak bisa dimanipulasi client.
func (r InitOrderRequest) ToObSetCache(txID, date string, userId, price int64) entity.CacheInitOrderRequest {
	return entity.CacheInitOrderRequest{
		TxID:     txID,
		EventID:  r.EventID,
		Date:     date,
		Quantity: r.Quantity,
		UserID:   userId,
		Price:    price,
	}
}
