package models

import (
	"github.com/alpha-omega-corp/cloud/app/user/pkg/models"
	"github.com/uptrace/bun"
	"time"
)

type Machine struct {
	bun.BaseModel `bun:"table:machines,alias:mch"`

	Id          int64  `json:"id" bun:",pk,autoincrement"`
	Name        string `json:"name" bun:"name,unique"`
	ContainerID string
	UserID      int64
	User        *models.User `bun:"rel:belongs-to"`
	CreatedAt   time.Time    `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time    `bun:",nullzero,notnull,default:current_timestamp"`
}
