package doctor

import (
	"context"
	"testing"
	"time"

	"github.com/HemanthGowda/hospital-management-go/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestRegisterCreatesDoctor(t *testing.T) {
	repository := newFakeUserRepo()
	service := NewService(repository)

	response, err := service.Register(context.Background(), RegisterDoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	require.NoError(t, err)
	require.Equal(t, "Dr. Kamath", response.Name)
	require.Equal(t, "Mysuru", response.Address)
	require.Equal(t, "Dental Specialist", response.Speciality)
	require.Equal(t, "kmc123", response.LicenseNumber)
	require.Equal(t, "Oral", response.Department)
}

func TestResgisterRejectsExistingLicenseNumber(t *testing.T) {
	repository := newFakeUserRepo()
	service := NewService(repository)
	response, _ := service.Register(context.Background(), RegisterDoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	_, err2 := service.Register(context.Background(), RegisterDoctorRequest{
		Name:          "Dr.Prasad",
		Address:       "Banglore",
		Speciality:    "Heart Specialist",
		LicenseNumber: "KMC123",
		Department:    "Cardiology",
	})

	var domainErr *domain.Error
	require.ErrorAs(t, err2, &domainErr)
	require.Equal(t, domain.KindConflict, domainErr.Kind)
	require.Contains(t, domainErr.Message, "Duplicate licensenumber not allowed")
	require.Contains(t, response.LicenseNumber, "kmc123")

}

func TestInvalidDoctorRegisterFails(t *testing.T) {
	repository := newFakeUserRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), RegisterDoctorRequest{
		Name:          "",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestRegisterWithLengthyLicenseNumberAndMissingNameFails(t *testing.T) {
	repository := newFakeUserRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), RegisterDoctorRequest{
		Name:          "",
		Address:       "Mysuru",
		Speciality:    "Dentist",
		LicenseNumber: "KMCH12345678910ABCDEFGHIJKLMNOPQRSTUVWXYZ1112334455",
		Department:    "Oral",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 2)
	require.Equal(t, domainErr.Message, "Request validation failed")

}

type fakeUserRepo struct {
	byname map[string]Doctor
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byname: map[string]Doctor{}}
}

func (f *fakeUserRepo) ExistsByLisenceNumber(_ context.Context, licenseNumber string) (bool, error) {
	_, ok := f.byname[licenseNumber]
	return ok, nil
}

func (f *fakeUserRepo) Create(_ context.Context, doctor Doctor) (Doctor, error) {
	if _, ok := f.byname[doctor.LicenseNumber]; ok {
		return Doctor{}, ErrDuplicate
	}
	doctor.Id = "doctor-1"

	doctor.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	doctor.UpdatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f.byname[doctor.LicenseNumber] = doctor
	return doctor, nil
}
