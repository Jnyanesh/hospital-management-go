package doctor

import (
	"context"
	"errors"
)

var (
	ErrNotFound  = errors.New("doctor not found")
	ErrDuplicate = errors.New("doctor already exists")
)

type Repository interface {
	ExistsByLisenceNumber(ctx context.Context, licenseNumber string) (bool, error)
	Create(ctx context.Context, doctor Doctor) (Doctor, error)
}

