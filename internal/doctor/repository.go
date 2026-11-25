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
	GetDoctorById(ctx context.Context, id string) (Doctor, error)
	SearchDoctorByName(ctx context.Context, name string, limit, offset int) ([]Doctor, error)
	CountDoctors(ctx context.Context, name string) (int, error)
	UpdateDoctor(ctx context.Context, doctor Doctor) (Doctor, error)
	HasLicenseNumber(ctx context.Context, id string, licenseNumber string) (bool, error)
	DeleteDoctor(ctx context.Context, id string)(error)
	DoctorExistsAsPrimaryDoctor(ctx context.Context, name string) (bool, error)
}
