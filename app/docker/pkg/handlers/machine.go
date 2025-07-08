package handlers

import (
	"context"
	"github.com/uptrace/bun"
)

type MachineService interface {
	Create(ctx context.Context, name string) error
}

type machineService struct {
	MachineService

	db *bun.DB
}

func NewMachineHandler(db *bun.DB) MachineService {
	return &machineService{
		db: db,
	}
}
