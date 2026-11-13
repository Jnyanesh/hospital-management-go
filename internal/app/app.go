package app

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/HemanthGowda/hospital-management-go/internal/auth"
	"github.com/HemanthGowda/hospital-management-go/internal/config"
	"github.com/HemanthGowda/hospital-management-go/internal/db"
	"github.com/HemanthGowda/hospital-management-go/internal/doctor"
	"github.com/HemanthGowda/hospital-management-go/internal/httpserver"
	"github.com/HemanthGowda/hospital-management-go/internal/password"
	"github.com/HemanthGowda/hospital-management-go/internal/patient"
	"github.com/HemanthGowda/hospital-management-go/internal/seed"
	"github.com/HemanthGowda/hospital-management-go/internal/user"
)

type Application struct {
	App      *fiber.App
	Pool     *pgxpool.Pool
	Config   config.Config
	users    user.Repository
	patients patient.Repository
	doctors  doctor.Repository
}

func New(ctx context.Context, cfg config.Config) (*Application, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		pool.Close()
		return nil, err
	}

	tokens, err := auth.NewJWTService(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTExpiration)
	if err != nil {
		pool.Close()
		return nil, err
	}

	users := user.NewPostgresRepository(pool)
	patients := patient.NewPostgresRepository(pool)
	doctors := doctor.NewPostgresRepository(pool)
	encoder := password.Bcrypt{}
	userService := user.NewService(users, encoder)
	authService := auth.NewService(users, encoder, tokens)
	patientService := patient.NewService(patients)
	doctorService := doctor.NewService(doctors)

	handler, err := httpserver.NewRouter(httpserver.Dependencies{
		Ping:     pool.Ping,
		Tokens:   tokens,
		Users:    user.NewHandler(userService),
		Auth:     auth.NewHandler(authService),
		Patients: patient.NewHandler(patientService),
		Doctors:  doctor.NewHandler(doctorService),
	})
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &Application{
		App:      handler,
		Pool:     pool,
		Config:   cfg,
		users:    users,
		patients: patients,
		doctors:  doctors,
	}, nil
}

func (a *Application) Seed(ctx context.Context) error {
	return seed.Run(ctx, a.Config, a.users, a.patients, password.Bcrypt{})
}

func (a *Application) Close() {
	if a.Pool != nil {
		a.Pool.Close()
	}
}
