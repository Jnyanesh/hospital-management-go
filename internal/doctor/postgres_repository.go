package doctor

import (
	"context"
	"errors"
	"strings"

	"github.com/HemanthGowda/hospital-management-go/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const doctorColumns = `
	id, fullname, speciality, address, licensenumber,
	department, created_at, updated_at
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

func (r *PostgresRepository) GetDoctorById(ctx context.Context, id string) (Doctor, error) {
	if _, err := uuid.Parse(id); err != nil {
		return Doctor{}, ErrNotFound
	}

	row := r.pool.QueryRow(ctx, `
		SELECT `+doctorColumns+`
		FROM doctors
		WHERE id = $1
	`, id)

	return toDoctorResponse(row)
}

func (r *PostgresRepository) SearchDoctorByName(ctx context.Context, name string, limit, offset int) ([]Doctor, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT `+doctorColumns+`
        FROM doctors
        WHERE fullname ILIKE $1
        ORDER BY fullname ASC
		LIMIT $2 OFFSET $3
    `, "%"+strings.TrimSpace(name)+"%", limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	doctors := make([]Doctor, 0)

	for rows.Next() {
		doctor, err := toDoctorResponse(rows)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (r *PostgresRepository) CountDoctors(ctx context.Context, name string) (int, error) {

	var count int
	err := r.pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM doctors
        WHERE fullname ILIKE $1
		`, "%"+strings.TrimSpace(name)+"%").Scan(&count)

	return count, err

}

func (r *PostgresRepository) ExistsByLisenceNumber(ctx context.Context, licenseNumber string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM doctors WHERE lower(licenseNumber) = lower($1))`, licenseNumber).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) HasLicenseNumber(ctx context.Context, id string, licenseNumber string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM doctors WHERE lower(licenseNumber) = lower($2) AND id != $1)`, id, licenseNumber).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) UpdateDoctor(ctx context.Context, doctor Doctor) (Doctor, error) {
	row := r.pool.QueryRow(ctx, `
		Update doctors
		Set fullname = $2,
		speciality = $3,
		address = $4,
		licensenumber = $5, 
		department = $6,
		updated_at = $7
		WHERE id = $1
		Returning `+doctorColumns+`
		`, doctor.Id, doctor.Name, doctor.Speciality, doctor.Address, doctor.LicenseNumber, doctor.Department, doctor.UpdatedAt)

	updated, err := toDoctorResponse(row)
	if errors.Is(err, ErrNotFound) {
		return Doctor{}, ErrNotFound
	}

	if db.IsUniqueViolation(err) {
		return Doctor{}, ErrDuplicate
	}

	return updated, err
}

type scanner interface {
	Scan(dest ...any) error
}

func toDoctorResponse(row scanner) (Doctor, error) {
	var doctorDetails Doctor

	err := row.Scan(
		&doctorDetails.Id,
		&doctorDetails.Name,
		&doctorDetails.Speciality,
		&doctorDetails.Address,
		&doctorDetails.LicenseNumber,
		&doctorDetails.Department,
		&doctorDetails.CreatedAt,
		&doctorDetails.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Doctor{}, ErrNotFound
	}

	if err != nil {
		return Doctor{}, err
	}

	doctorDetails.CreatedAt = doctorDetails.CreatedAt.UTC()
	doctorDetails.UpdatedAt = doctorDetails.UpdatedAt.UTC()

	return doctorDetails, nil
}
