package entities

import (
	"time"
)

type BaseEntity struct {
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"` // first joined
}
