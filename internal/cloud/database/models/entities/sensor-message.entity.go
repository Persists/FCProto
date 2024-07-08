package entities

import (
	"time"

	"github.com/google/uuid"
)

type SensorMessageEntity struct {
	ID uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`

	Timestamp    time.Time
	Content      string
	ClientIpAddr string
	Client       *ClientEntity `bun:"rel:belongs-to,join:client_ip_addr=ip_addr"`

	BaseEntity `bun:",embed:base_"`
}
