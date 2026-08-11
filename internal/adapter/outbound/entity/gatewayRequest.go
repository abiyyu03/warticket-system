package entity

import (
	"encoding/json"
)

type GatewayRequest struct {
	BaseModel
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Providers string `gorm:"column:providers;not null;index" json:"providers"`
	Request   string `gorm:"column:request;type:text" json:"request"`
	Response  string `gorm:"column:response;type:text" json:"response"`
	// Nullable: kosong kalau belum ada response (timeout / gagal konek).
	ResponseCode   *int   `gorm:"column:response_code;index" json:"response_code,omitempty"`
	RequestHeader  string `gorm:"column:request_header;type:text" json:"request_header"`
	ResponseHeader string `gorm:"column:response_header;type:text" json:"response_header"`
}

func (GatewayRequest) TableName() string {
	return "gateway_requests"
}

func (u GatewayRequest) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u GatewayRequest) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}
