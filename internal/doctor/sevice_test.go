package doctor

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/HemanthGowda/hospital-management-go/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestRegisterCreatesDoctor(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	response, err := service.Register(context.Background(), DoctorRequest{
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
	repository := newFakeDoctorRepo()
	service := NewService(repository)
	response, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	_, err2 := service.Register(context.Background(), DoctorRequest{
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

func TestDoctorValuesWithSpacesAareTrimmedAndRegistered(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	response, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath    ",
		Address:       "Mysuru   ",
		Speciality:    "    Dental Specialist",
		LicenseNumber: "   KMC123",
		Department:    "    Oral",
	})
	require.NoError(t, err)
	require.Equal(t, "Dr. Kamath", response.Name)
	require.Equal(t, "Mysuru", response.Address)
	require.Equal(t, "Dental Specialist", response.Speciality)
	require.Equal(t, "kmc123", response.LicenseNumber)
	require.Equal(t, "Oral", response.Department)
}

func TestInvalidDoctorRegisterFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
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
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
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

func TestRegisterEmptyDepartmentFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Prasad",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestRegisterDoctorWithLengthyDepartmentFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Prasad",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    strings.Repeat("s", 90),
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestRegisterDoctorWithLengthyNameFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          strings.Repeat("abc", 80),
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Dental Care",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestRegisterDoctorWithLengthyAddress(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       strings.Repeat("Mysuru", 80),
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Dental Care",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestRegisterDoctorWithLengthySpeciality(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    strings.Repeat("Dental Specialist", 80),
		LicenseNumber: "KMC123",
		Department:    "Dental Care",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestRegisterDoctorWithEmptySpecialityFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "",
		LicenseNumber: "KMC123",
		Department:    "Dental Care",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestRegisterDoctorWithEmptyAdressAndLicenseNumberFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "",
		Speciality:    "Dental",
		LicenseNumber: "",
		Department:    "Dental Care",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 2)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestRegisterDoctorWithMultipleEmptyFieldsFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err1 := service.Register(context.Background(), DoctorRequest{
		Name:          "",
		Address:       "",
		Speciality:    "",
		LicenseNumber: "",
		Department:    "Dental Care",
	})

	_, err2 := service.Register(context.Background(), DoctorRequest{
		Name:          "",
		Address:       "",
		Speciality:    "",
		LicenseNumber: "",
		Department:    "",
	})
	var domainErr1 *domain.Error
	var domainErr2 *domain.Error

	require.ErrorAs(t, err1, &domainErr1)
	require.ErrorAs(t, err2, &domainErr2)
	require.Equal(t, len(domainErr1.Violations), 4)
	require.Equal(t, len(domainErr2.Violations), 5)
	require.Equal(t, domainErr1.Message, "Request validation failed")
	require.Equal(t, domainErr2.Message, "Request validation failed")
}

func TestRegisterDoctorWithMultipleFieldsLengthyFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          strings.Repeat("Doctor", 80),
		Address:       strings.Repeat("Address", 80),
		Speciality:    strings.Repeat("Dental Specialist", 80),
		LicenseNumber: strings.Repeat("abc", 80),
		Department:    strings.Repeat("Dental Care", 80),
	})

	var domainErr1 *domain.Error

	require.ErrorAs(t, err, &domainErr1)
	require.Equal(t, len(domainErr1.Violations), 5)
	require.Equal(t, domainErr1.Message, "Request validation failed")

}

func TestGetDoctorByIdReturnsDoctorDetails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	response, err1 := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	require.NoError(t, err1)

	doctorDetails, err2 := service.GetDoctorById(context.Background(), response.Id)

	require.NoError(t, err2)
	require.Equal(t, "doctor-1", doctorDetails.Id)
	require.Equal(t, "Dr. Kamath", doctorDetails.Name)
	require.Equal(t, "Mysuru", doctorDetails.Address)
	require.Equal(t, "Dental Specialist", doctorDetails.Speciality)
	require.Equal(t, "kmc123", doctorDetails.LicenseNumber)
	require.Equal(t, "Oral", doctorDetails.Department)
}

func TestGetDoctorByIdWithoutExistingDoctor(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.GetDoctorById(context.Background(), "Doctor-1")

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindNotFound, domainErr.Kind)
	require.Equal(t, domainErr.Message, "Doctor Not Found")

}
func TestGetDoctorByIdWithEmptyId(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err1 := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	require.NoError(t, err1)

	_, err2 := service.GetDoctorById(context.Background(), "")

	var domainErr *domain.Error
	require.ErrorAs(t, err2, &domainErr)
	require.Equal(t, domainErr.Message, "Doctor Not Found")
}

func TestSearchDoctorWithName(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err1 := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	_, err2 := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamala",
		Address:       "Mysuru",
		Speciality:    "Heart Specialist",
		LicenseNumber: "KMC345",
		Department:    "Cardiology",
	})
	require.NoError(t, err1)
	require.NoError(t, err2)

	doctors, err3 := service.SearchDoctorByName(context.Background(), "Dr. Kam", 1, 20)

	require.NoError(t, err3)

	require.Equal(t, len(doctors.Doctors), 2)
	require.Equal(t, doctors.Total, 2)

}

func TestSearchDoctorWithLimitSetTo1ReturnsSingleDoctor(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err1 := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	_, err2 := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamala",
		Address:       "Mysuru",
		Speciality:    "Heart Specialist",
		LicenseNumber: "KMC345",
		Department:    "Cardiology",
	})
	require.NoError(t, err1)
	require.NoError(t, err2)

	doctors, err3 := service.SearchDoctorByName(context.Background(), "Dr. Kam", 1, 1)

	require.NoError(t, err3)

	require.Equal(t, len(doctors.Doctors), 1)
	require.Equal(t, doctors.Total, 2)

}
func TestSearchDoctorWithLimitSetTo2Returns2Doctors(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	require.NoError(t, err)

	_, err = service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamala",
		Address:       "Mysuru",
		Speciality:    "Heart Specialist",
		LicenseNumber: "KMC345",
		Department:    "Cardiology",
	})
	require.NoError(t, err)

	_, err = service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Karan",
		Address:       "Banglore",
		Speciality:    "Nephrologist",
		LicenseNumber: "KMC145",
		Department:    "Nephrology",
	})
	require.NoError(t, err)

	doctors, err := service.SearchDoctorByName(context.Background(), "Dr. Ka", 1, 2)

	require.NoError(t, err)

	require.Equal(t, len(doctors.Doctors), 2)
	require.Equal(t, 3, doctors.Total)
}

func TestSearchDoctorWithoutRelevantDoctorreturnEmptyPage(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	_, err := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})
	require.NoError(t, err)

	_, err = service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamala",
		Address:       "Mysuru",
		Speciality:    "Heart Specialist",
		LicenseNumber: "KMC345",
		Department:    "Cardiology",
	})
	require.NoError(t, err)

	_, err = service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Karan",
		Address:       "Banglore",
		Speciality:    "Nephrologist",
		LicenseNumber: "KMC145",
		Department:    "Nephrology",
	})
	require.NoError(t, err)

	doctors, err := service.SearchDoctorByName(context.Background(), "za", 1, 2)

	require.NoError(t, err)

	require.Equal(t, len(doctors.Doctors), 0)
	require.Equal(t, doctors.Total, 0)
}

func TestUpdateDoctorWithExistingLicenseNumber(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	doctor, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamala",
		Address:       "Mysuru",
		Speciality:    "Heart Specialist",
		LicenseNumber: "KMC345",
		Department:    "Cardiology",
	})

	updatedDoctor := DoctorRequest{
		Name:          "Dr. Karan",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC345",
		Department:    "Oral",
	}

	_, err := service.UpdateDoctor(context.Background(), doctor.Id, updatedDoctor)

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domainErr.Message, "License Number already exists")

}

func TestUpdateDoctorWitValidDetails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	doctor, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	updatedDoctor := DoctorRequest{
		Name:          "Dr. Karan",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	}

	response, err := service.UpdateDoctor(context.Background(), doctor.Id, updatedDoctor)

	require.NoError(t, err)
	require.Equal(t, "Dr. Karan", response.Name)
	require.Equal(t, "Mysuru", response.Address)
	require.Equal(t, "Dental Specialist", response.Speciality)
	require.Equal(t, "kmc123", response.LicenseNumber)
	require.Equal(t, "Oral", response.Department)

}

func TestUpdateDoctorWithMissingFields(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	doctor, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	updatedDoctor := DoctorRequest{
		Name:          "Dr. Karan",
		Address:       "",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	}

	_, err := service.UpdateDoctor(context.Background(), doctor.Id, updatedDoctor)

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestUpdateDoctorWithrailingSpaceInValuesSucceeds(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	doctor, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	updatedDoctor := DoctorRequest{
		Name:          "Dr. Karan    ",
		Address:       "Banglore    ",
		Speciality:    "Dental Specialist    ",
		LicenseNumber: "     KMC123",
		Department:    "Oral  ",
	}

	response, err := service.UpdateDoctor(context.Background(), doctor.Id, updatedDoctor)

	require.NoError(t, err)
	require.Equal(t, "Dr. Karan", response.Name)
	require.Equal(t, "Banglore", response.Address)
	require.Equal(t, "Dental Specialist", response.Speciality)
	require.Equal(t, "kmc123", response.LicenseNumber)
	require.Equal(t, "Oral", response.Department)

}

func TestUpdateDoctorWithEdgeCases(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	doctor, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	updatedDoctor := DoctorRequest{
		Name:          strings.Repeat("Doctor", 80),
		Address:       strings.Repeat("Banglore", 160),
		Speciality:    strings.Repeat("Dental Specialist", 70),
		LicenseNumber: strings.Repeat("KMC123", 80),
		Department:    strings.Repeat("Oral", 80),
	}
	_, err := service.UpdateDoctor(context.Background(), doctor.Id, updatedDoctor)

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 5)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestUpdateDoctorWithEmptyFieldsFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	doctor, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	updatedDoctor := DoctorRequest{
		Name:          "",
		Address:       "",
		Speciality:    "",
		LicenseNumber: "",
		Department:    "",
	}
	_, err := service.UpdateDoctor(context.Background(), doctor.Id, updatedDoctor)

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 5)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

func TestDeleteValidDoctorSuceeds(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	doctor, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	repository.primaryDoctorExists = false
	response, err := service.DeleteDoctor(context.Background(), doctor.Id)

	require.NoError(t, err)
	require.Equal(t, response.Deleted, true)
	require.Equal(t, response.Name, "Dr. Kamath")
	require.Equal(t, response.Id, doctor.Id)
}

func TestDeleteDoctorFailsIfExistsAsPrimaryDoctor(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	doctor, _ := service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	repository.primaryDoctorExists = true

	response, err := service.DeleteDoctor(context.Background(), doctor.Id)

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domainErr.Kind, domain.KindForbidden)
	require.Equal(t, response.Deleted, false)
}

func TestDeleteValidDoctorWithoutValidIdFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	repository.primaryDoctorExists = true

	response, err := service.DeleteDoctor(context.Background(), "123")

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domainErr.Message, "Doctor Not Found")
	require.Equal(t, response.Deleted, false)
}

func TestDeleteValidDoctorWithoutIdFails(t *testing.T) {
	repository := newFakeDoctorRepo()
	service := NewService(repository)

	service.Register(context.Background(), DoctorRequest{
		Name:          "Dr. Kamath",
		Address:       "Mysuru",
		Speciality:    "Dental Specialist",
		LicenseNumber: "KMC123",
		Department:    "Oral",
	})

	repository.primaryDoctorExists = true

	response, err := service.DeleteDoctor(context.Background(), "")

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domainErr.Message, "Doctor Not Found")
	require.Equal(t, response.Deleted, false)
}

type fakeDoctorRepo struct {
	doctors             map[string]Doctor
	nextId              int
	primaryDoctorExists bool
	checkError          error
}

func newFakeDoctorRepo() *fakeDoctorRepo {
	return &fakeDoctorRepo{doctors: map[string]Doctor{}}
}

func (f *fakeDoctorRepo) ExistsByLisenceNumber(_ context.Context, licenseNumber string) (bool, error) {
	ok := false
	for id := range f.doctors {
		if f.doctors[id].LicenseNumber == licenseNumber {
			ok = true
			break
		}
	}

	return ok, nil
}

func (f *fakeDoctorRepo) Create(_ context.Context, doctor Doctor) (Doctor, error) {
	f.nextId++
	doctor.Id = fmt.Sprintf("doctor-%d", f.nextId)
	doctor.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	doctor.UpdatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f.doctors[doctor.Id] = doctor

	return doctor, nil
}

func (f *fakeDoctorRepo) GetDoctorById(_ context.Context, id string) (Doctor, error) {
	doctorDetails, ok := f.doctors[id]

	if !ok {
		return Doctor{}, ErrNotFound
	}
	return doctorDetails, nil
}

func (f *fakeDoctorRepo) SearchDoctorByName(_ context.Context, name string, limit, offset int) ([]Doctor, error) {
	doctors := make([]Doctor, 0)

	for _, doctor := range f.doctors {
		if strings.Contains(strings.ToLower(doctor.Name), strings.ToLower(name)) {
			doctors = append(doctors, doctor)
		}
	}

	sort.Slice(doctors, func(i, j int) bool {
		return doctors[i].Name < doctors[j].Name
	})

	if offset >= len(doctors) {
		return []Doctor{}, nil
	}

	end := offset + limit
	if end > len(doctors) {
		end = len(doctors)
	}

	return doctors[offset:end], nil
}

func (f *fakeDoctorRepo) UpdateDoctor(_ context.Context, doctor Doctor) (Doctor, error) {
	if _, ok := f.doctors[doctor.Id]; !ok {
		return Doctor{}, ErrNotFound
	}
	doctor.UpdatedAt = time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	f.doctors[doctor.Id] = doctor

	return doctor, nil
}

func (f *fakeDoctorRepo) HasLicenseNumber(_ context.Context, id string, licenseNumber string) (bool, error) {
	ok := false
	for i := range f.doctors {
		if f.doctors[i].LicenseNumber == licenseNumber && f.doctors[i].Id != id {
			ok = true
			break
		}
	}

	return ok, nil
}

func (f *fakeDoctorRepo) CountDoctors(_ context.Context, name string) (int, error) {
	var count int

	for id := range f.doctors {
		if strings.Contains(strings.ToLower(f.doctors[id].Name), strings.ToLower(name)) {
			count += 1
		}
	}
	return count, nil
}

func (f *fakeDoctorRepo) DeleteDoctor(_ context.Context, id string) error {
	for i := range f.doctors {
		if f.doctors[i].Id == id {
			delete(f.doctors, id)
		}
	}
	return nil
}

func (f *fakeDoctorRepo) DoctorExistsAsPrimaryDoctor(_ context.Context, name string) (bool, error) {

	return f.primaryDoctorExists, f.checkError
}
