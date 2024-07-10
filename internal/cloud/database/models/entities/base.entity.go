package entities

import (
	"time"
)

// BaseEntity is the base entity for all entities
type BaseEntity struct {
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"` // first joined
}
