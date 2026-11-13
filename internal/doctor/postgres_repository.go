package doctor

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HemanthGowda/hospital-management-go/internal/db"
)

const doctorColumns = `
	id::text, full_name, speciality, address, licensenumber,
	department, craeted_at, updated_at
`

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, doctor Doctor) (Doctor, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Doctor{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	doctor.Id = uuid.NewString()
	err = tx.QueryRow(ctx, `
		INSERT INTO doctors (id, fullname, speciality, address, licensenumber, department, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`, doctor.Id, doctor.Name, doctor.Speciality, doctor.Address, doctor.LicenseNumber, doctor.Department, doctor.CreatedAt, doctor.UpdatedAt).Scan(&doctor.CreatedAt, &doctor.UpdatedAt)
	if db.IsUniqueViolation(err) {
		return Doctor{}, ErrDuplicate
	}
	if err != nil {
		return Doctor{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Doctor{}, err
	}
	return doctor, nil
}

func (r *PostgresRepository) ExistsByLisenceNumber(ctx context.Context, licenseNumber string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM doctors WHERE lower(licenseNumber) = lower($1))`, licenseNumber).Scan(&exists)
	return exists, err
}
