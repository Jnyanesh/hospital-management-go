package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"golang.org/x/crypto/bcrypt"

	"github.com/HemanthGowda/hospital-management-go/internal/config"
)

var application *Application

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("hospital_management_test"),
		postgres.WithUsername("hospital"),
		postgres.WithPassword("hospital"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres: %v\n", err)
		os.Exit(1)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres dsn: %v\n", err)
		os.Exit(1)
	}

	application, err = New(ctx, config.Config{
		DatabaseURL:       dsn,
		JWTSecret:         "test-secret-key-that-is-long-enough-for-hs256-signing",
		JWTExpiration:     time.Hour,
		JWTIssuer:         "hospital-management-test",
		SeedAdminEmail:    "admin@hospital.local",
		SeedAdminPassword: "ChangeMe123!",
		SeedAdminFullName: "Hospital Admin",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "application: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	application.Close()
	_ = container.Terminate(context.Background())
	os.Exit(code)
}

func TestRegisterLoginAndUpdatePatientHappyPath(t *testing.T) {
	reset(t)
	insertUser(t, "admin@hospital.local", "ChangeMe123!", "Hospital Admin", "ADMIN")
	patientID := insertPatient(t)

	email := "intern+" + uuid.NewString() + "@hospital.local"
	status, body := request(t, http.MethodPost, "/api/v1/users", "", map[string]any{
		"email":    email,
		"password": "InternPass1",
		"fullName": "New Intern",
	})
	require.Equal(t, http.StatusCreated, status)
	require.Equal(t, email, body["email"])
	roles, _ := body["roles"].([]any)
	require.Equal(t, "PATIENT", roles[0])

	token := login(t, "admin@hospital.local", "ChangeMe123!")
	status, body = request(t, http.MethodPut, "/api/v1/patients/"+patientID, token, validPatient())
	require.Equal(t, http.StatusOK, status)
	require.Equal(t, "Ada Lovelace", body["fullName"])
	require.Equal(t, "MRN-100002", body["medicalRecordNumber"])
	require.Equal(t, "1815-12-10", body["dateOfBirth"])
}

func TestRegisterRejectsInvalidPayload(t *testing.T) {
	reset(t)
	status, body := request(t, http.MethodPost, "/api/v1/users", "", map[string]any{
		"email":    "not-an-email",
		"password": "short",
		"fullName": "",
	})
	require.Equal(t, http.StatusBadRequest, status)
	require.Equal(t, "Validation Failed", body["error"])
	require.NotEmpty(t, body["violations"])
}

func TestRegisterValidDoctorSucceeds(t *testing.T) {
	reset(t)
	insertUser(t, "admin@hospital.local", "ChangeMe123!", "Hospital Admin", "ADMIN")

	token := login(t, "admin@hospital.local", "ChangeMe123!")
	status, body := request(t, http.MethodPost, "/api/v1/doctors", token, validDoctor())
	require.Equal(t, http.StatusCreated, status)
	require.Equal(t, "Dr. Kamath", body["name"])
	require.Equal(t, "Dental specialist", body["speciality"])
	require.Equal(t, "Mysore", body["address"])
	require.Equal(t, "kmc123", body["license_number"])
}

func TestUpdatePatientRequiresAuthentication(t *testing.T) {
	reset(t)
	patientID := insertPatient(t)
	status, _ := request(t, http.MethodPut, "/api/v1/patients/"+patientID, "", validPatient())
	require.Equal(t, http.StatusUnauthorized, status)
}

func TestUpdatePatientForbidsPatientRole(t *testing.T) {
	reset(t)
	insertUser(t, "patient@hospital.local", "PatientPass1", "Hospital Patient", "PATIENT")
	patientID := insertPatient(t)
	token := login(t, "patient@hospital.local", "PatientPass1")
	status, body := request(t, http.MethodPut, "/api/v1/patients/"+patientID, token, validPatient())
	require.Equal(t, http.StatusForbidden, status)
	require.Equal(t, "Forbidden", body["error"])
}

func TestLoginRejectsBadPassword(t *testing.T) {
	reset(t)
	insertUser(t, "admin@hospital.local", "ChangeMe123!", "Hospital Admin", "ADMIN")
	status, _ := request(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email":    "admin@hospital.local",
		"password": "WrongPassword1",
	})
	require.Equal(t, http.StatusUnauthorized, status)
}

func TestRegisterDocterWithoutLoginFails(t *testing.T) {
	reset(t)
	insertUser(t, "admin@hospital.local", "ChangeMe123!", "Hospital Admin", "ADMIN")

	insertDoctor(t, "Dr. Singh", "Dental specialist", "Mysore", "KMC123", "Dental")
	status, _ := request(t, http.MethodPost, "/api/v1/doctors", "", validDoctor())
	require.Equal(t, http.StatusUnauthorized, status)
}

func TestDoctorRegisterByLoggedInPatientFails(t *testing.T) {
	reset(t)
	insertUser(t, "patient@hospital.local", "ChangeMe123!", "Hospital Patient", "PATIENT")

	token := login(t, "patient@hospital.local", "ChangeMe123!")
	insertDoctor(t, "Dr. Singh", "Dental specialist", "Mysore", "KMC123", "Dental")
	status, _ := request(t, http.MethodPost, "/api/v1/doctors", token, validDoctor())
	require.Equal(t, http.StatusForbidden, status)
}

func TestRegisterDoctorForbidsDuplicateLisenceNumber(t *testing.T) {
	reset(t)
	insertUser(t, "admin@hospital.local", "ChangeMe123!", "Hospital Admin", "ADMIN")

	token := login(t, "admin@hospital.local", "ChangeMe123!")
	insertDoctor(t, "Dr. Singh", "Dental specialist", "Mysore", "KMC123", "Dental")
	status, _ := request(t, http.MethodPost, "/api/v1/doctors", token, validDoctor())
	require.Equal(t, http.StatusConflict, status)
}

func reset(t *testing.T) {
	t.Helper()
	_, err := application.Pool.Exec(context.Background(), `TRUNCATE TABLE patients, doctors, user_roles, users RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
}

func insertUser(t *testing.T, email, rawPassword, fullName, role string) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.MinCost)
	require.NoError(t, err)
	id := uuid.NewString()
	_, err = application.Pool.Exec(context.Background(), `
		INSERT INTO users (id, email, password_hash, full_name, enabled)
		VALUES ($1, $2, $3, $4, true)
	`, id, email, string(hash), fullName)
	require.NoError(t, err)
	_, err = application.Pool.Exec(context.Background(), `INSERT INTO user_roles (user_id, role) VALUES ($1, $2)`, id, role)
	require.NoError(t, err)
}

func insertDoctor(t *testing.T, name, speciality, address, license_number, department string) {
	t.Helper()
	id := uuid.NewString()
	_, err := application.Pool.Exec(context.Background(), `
		INSERT INTO doctors (
			id, fullname, speciality, address, licensenumber, department, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, '1800-04-05', '1999-06-07')
	`, id, name, speciality, address, license_number, department)
	require.NoError(t, err)
}

func insertPatient(t *testing.T) string {
	t.Helper()
	id := uuid.NewString()
	_, err := application.Pool.Exec(context.Background(), `
		INSERT INTO patients (
			id, full_name, primary_doctor_name, medical_record_number, date_of_birth, gender, notes
		) VALUES ($1, 'Seed Patient', 'Dr. Seed', 'MRN-100001', '1990-01-01', 'UNKNOWN', 'A seeded patient')
	`, id)
	require.NoError(t, err)
	return id
}

func login(t *testing.T, email, rawPassword string) string {
	t.Helper()
	status, body := request(t, http.MethodPost, "/api/v1/auth/login", "", map[string]any{
		"email":    email,
		"password": rawPassword,
	})
	require.Equal(t, http.StatusOK, status)
	token, _ := body["accessToken"].(string)
	require.NotEmpty(t, token)
	return token
}

func request(t *testing.T, method, path, token string, payload any) (int, map[string]any) {
	t.Helper()
	var reader *bytes.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		require.NoError(t, err)
		reader = bytes.NewReader(encoded)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, "http://localhost"+path, reader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := application.App.Test(req, fiber.TestConfig{Timeout: 30 * time.Second})
	require.NoError(t, err)
	defer response.Body.Close()

	var body map[string]any
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	return response.StatusCode, body
}

func validPatient() map[string]any {
	return map[string]any{
		"fullName":            "Ada Lovelace",
		"primaryDoctorName":   "Dr. Grace Hopper",
		"medicalRecordNumber": "MRN-100002",
		"dateOfBirth":         "1815-12-10",
		"gender":              "FEMALE",
		"phone":               "+15551230001",
		"bloodGroup":          "O+",
		"notes":               "Updated by integration test",
	}
}

func validDoctor() map[string]any {
	return map[string]any{
		"name":           "Dr. Kamath",
		"speciality":     "Dental specialist",
		"address":        "Mysore",
		"license_number": "KMC123",
		"department":     "Dental",
	}
}
