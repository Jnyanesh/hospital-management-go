package doctor

import "time"

type Doctor struct {
	Id            string
	Name          string
	Speciality    string
	Address       string
	LicenseNumber string
	Department    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type RegisterDoctorRequest struct {
	Name          string `json:"name"`
	Speciality    string `json:"speciality"`
	Address       string `json:"address"`
	LicenseNumber string `json:"license_number"`
	Department    string `json:"department"`
}

type DoctorResponse struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Speciality    string    `json:"speciality"`
	Address       string    `json:"address"`
	LicenseNumber string    `json:"license_number"`
	Department    string    `json:"department"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
