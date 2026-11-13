package patient

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/HemanthGowda/hospital-management-go/internal/domain"
)

func TestUpdateReplacesMutableFields(t *testing.T) {
	repository := newFakePatientRepo()
	repository.patients["patient-1"] = Patient{
		ID:                  "patient-1",
		FullName:            "Old Name",
		MedicalRecordNumber: "MRN-100001",
		Gender:              "FEMALE",
	}
	service := NewService(repository)

	response, err := service.Update(context.Background(), "patient-1", UpdatePatientRequest{
		FullName:            str("Ada Lovelace"),
		PrimaryDoctorName:   str("Dr. Grace Hopper"),
		MedicalRecordNumber: str("mrn-100002"),
		DateOfBirth:         str("1815-12-10"),
		Gender:              str("female"),
		Phone:               str("+15551230001"),
		BloodGroup:          str("o+"),
		Notes:               str("Updated notes"),
	})
	require.NoError(t, err)
	require.Equal(t, "Ada Lovelace", response.FullName)
	require.Equal(t, "Dr. Grace Hopper", response.PrimaryDoctorName)
	require.Equal(t, "MRN-100002", response.MedicalRecordNumber)
	require.Equal(t, "1815-12-10", response.DateOfBirth)
	require.Equal(t, "FEMALE", response.Gender)
	require.Equal(t, "O+", response.BloodGroup)
}

func TestUpdateFailsWhenPatientMissing(t *testing.T) {
	service := NewService(newFakePatientRepo())
	_, err := service.Update(context.Background(), "missing", sampleRequest())
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindNotFound, domainErr.Kind)
}

func TestUpdateFailsWhenMedicalRecordNumberConflicts(t *testing.T) {
	repository := newFakePatientRepo()
	repository.patients["patient-1"] = Patient{ID: "patient-1", MedicalRecordNumber: "MRN-100001"}
	repository.patients["patient-2"] = Patient{ID: "patient-2", MedicalRecordNumber: "MRN-100002"}
	service := NewService(repository)

	_, err := service.Update(context.Background(), "patient-1", sampleRequest())
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindConflict, domainErr.Kind)
	require.Contains(t, domainErr.Message, "Medical record number")
}

func TestUpdateFailsWhenDateOfBirthIsInTheFuture(t *testing.T) {
	repository := newFakePatientRepo()
	repository.patients["patient-1"] = Patient{ID: "patient-1", MedicalRecordNumber: "MRN-100001"}
	service := NewService(repository)

	request := sampleRequest()
	future := time.Now().UTC().AddDate(1, 0, 0).Format("2006-01-02")
	request.DateOfBirth = &future

	_, err := service.Update(context.Background(), "patient-1", request)
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindBadRequest, domainErr.Kind)
	require.Contains(t, domainErr.Message, "future")
}

func sampleRequest() UpdatePatientRequest {
	return UpdatePatientRequest{
		FullName:            str("Ada Lovelace"),
		PrimaryDoctorName:   str("Dr. Grace Hopper"),
		MedicalRecordNumber: str("MRN-100002"),
		DateOfBirth:         str("1815-12-10"),
		Gender:              str("FEMALE"),
		Phone:               str("+15551230001"),
		BloodGroup:          str("O+"),
		Notes:               str("Notes"),
	}
}

func str(value string) *string { return &value }

type fakePatientRepo struct {
	patients map[string]Patient
}

func newFakePatientRepo() *fakePatientRepo {
	return &fakePatientRepo{patients: map[string]Patient{}}
}

func (f *fakePatientRepo) FindByID(_ context.Context, id string) (Patient, error) {
	patient, ok := f.patients[id]
	if !ok {
		return Patient{}, ErrNotFound
	}
	return patient, nil
}

func (f *fakePatientRepo) MedicalRecordNumberUsedByOther(_ context.Context, medicalRecordNumber, id string) (bool, error) {
	for _, patient := range f.patients {
		if patient.MedicalRecordNumber == medicalRecordNumber && patient.ID != id {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakePatientRepo) Update(_ context.Context, patient Patient) (Patient, error) {
	if _, ok := f.patients[patient.ID]; !ok {
		return Patient{}, ErrNotFound
	}
	f.patients[patient.ID] = patient
	return patient, nil
}

func (f *fakePatientRepo) ExistsByMedicalRecordNumber(_ context.Context, medicalRecordNumber string) (bool, error) {
	for _, patient := range f.patients {
		if patient.MedicalRecordNumber == medicalRecordNumber {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakePatientRepo) Insert(_ context.Context, patient Patient) (Patient, error) {
	f.patients[patient.ID] = patient
	return patient, nil
}
