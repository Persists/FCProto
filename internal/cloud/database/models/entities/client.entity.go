package entities

import "time"

// ClientEntity is the entity for the client
type ClientEntity struct {
	IpAddr   string `bun:",notnull,pk"` // ip address
	LastSeen time.Time

	BaseEntity `bun:",embed:base_"`
}
