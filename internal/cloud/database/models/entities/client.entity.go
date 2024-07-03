package entities

import "time"

type ClientEntity struct {
	IpAddr     string `bun:",notnull,pk"` // ip address
	LastSeen   time.Time
	NotifyAddr string

	BaseEntity `bun:",embed:base_"`
}
