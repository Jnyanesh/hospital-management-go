package doctor

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/HemanthGowda/hospital-management-go/internal/domain"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Register(ctx context.Context, request RegisterDoctorRequest) (DoctorResponse, error) {
	if violation := validateCreateRequest(request); len(violation) > 0 {
		return DoctorResponse{}, domain.ValidationFailed(violation)

	}
	licenseNumber := strings.ToLower(strings.TrimSpace(request.LicenseNumber))
	fullName := strings.TrimSpace(request.Name)
	speciality := strings.TrimSpace(request.Speciality)
	department := strings.TrimSpace(request.Department)
	address := strings.TrimSpace(request.Address)
	createdat := time.Now()
	upadatedat := time.Now()

	exists, err := s.repository.ExistsByLisenceNumber(ctx, licenseNumber)
	if err != nil {
		return DoctorResponse{}, err
	}
	if exists {
		return DoctorResponse{}, domain.Conflict("Duplicate licensenumber not allowed")
	}

	created, err := s.repository.Create(ctx, Doctor{
		Name:          fullName,
		Speciality:    speciality,
		Address:       address,
		LicenseNumber: licenseNumber,
		Department:    department,
		CreatedAt:     createdat,
		UpdatedAt:     upadatedat,
	})

	if err != nil {
		return DoctorResponse{}, err
	}

	return toResponse(created), nil
}

func validateCreateRequest(request RegisterDoctorRequest) []domain.Violation {
	var violations []domain.Violation

	switch {
	case strings.TrimSpace(request.Name) == "":
		violations = append(violations, domain.Violation{Field: "name", Message: "Full name is required"})
	case utf8.RuneCountInString(strings.TrimSpace(request.Name)) < 2 || utf8.RuneCountInString(strings.TrimSpace(request.Name)) > 160:
		violations = append(violations, domain.Violation{Field: "name", Message: "Full name must be between 2 and 160 characters"})
	}

	switch {
	case strings.TrimSpace(request.LicenseNumber) == "":
		violations = append(violations, domain.Violation{Field: "lisence_number", Message: "License Number is required"})
	case utf8.RuneCountInString(strings.TrimSpace(request.LicenseNumber)) < 2 || utf8.RuneCountInString(strings.TrimSpace(request.LicenseNumber)) > 40:
		violations = append(violations, domain.Violation{Field: "lisence_number", Message: "License number must be between 2 and 40 characters"})
	}

	switch {
	case strings.TrimSpace(request.Speciality) == "":
		violations = append(violations, domain.Violation{Field: "speciality", Message: "Speciality is required"})
	case utf8.RuneCountInString(strings.TrimSpace(request.Speciality)) < 2 || utf8.RuneCountInString(strings.TrimSpace(request.Speciality)) > 80:
		violations = append(violations, domain.Violation{Field: "speciality", Message: "Speciality must be between 2 and 80 characters"})
	}

	switch {
	case strings.TrimSpace(request.Department) == "":
		violations = append(violations, domain.Violation{Field: "department", Message: "Deparment is required"})
	case utf8.RuneCountInString(strings.TrimSpace(request.Department)) < 2 || utf8.RuneCountInString(strings.TrimSpace(request.Department)) > 80:
		violations = append(violations, domain.Violation{Field: "department", Message: "department  must be between 2 and 80 characters"})
	}

	return violations
}

func toResponse(doctor Doctor) DoctorResponse {
	return DoctorResponse{
		Id:            doctor.Id,
		Name:          doctor.Name,
		Speciality:    doctor.Speciality,
		Address:       doctor.Address,
		LicenseNumber: doctor.LicenseNumber,
		Department:    doctor.Department,
		CreatedAt:     doctor.CreatedAt.UTC(),
		UpdatedAt:     doctor.UpdatedAt.UTC(),
	}
}
