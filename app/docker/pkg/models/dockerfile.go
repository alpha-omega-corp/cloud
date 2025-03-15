package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Dockerfile struct {
	bun.BaseModel `bun:"table:dockerfiles,alias:df"`

	Id        int64     `json:"id" bun:",pk,autoincrement"`
	Name      string    `json:"name" bun:"name,unique"`
	Content   []byte    `json:"content" bun:"content"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
