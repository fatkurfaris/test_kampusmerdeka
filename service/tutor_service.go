package service

import (
	"errors"
	"ruang_belajar/models/tutors"
	"ruang_belajar/repository/database"

	"golang.org/x/crypto/bcrypt"
)

type TutorService interface {
	RegisterTutor(input tutors.RegisterUserInput) (tutors.Tutor, error)
	GetTutorByID(ID int) (tutors.Tutor, error)
	Login(input tutors.LogisUserInput) (tutors.Tutor, error)
	UpdateTutor(inputID tutors.GetTutorInput, inputData tutors.CreateTutorInput) (tutors.Tutor, error)
}

type tutorService struct {
	repository database.TutorRepository
}

func NewTutorService(repository database.TutorRepository) *tutorService {
	return &tutorService{repository}
}

func (s *tutorService) RegisterTutor(input tutors.RegisterUserInput) (tutors.Tutor, error) {
	tutor := tutors.Tutor{}
	tutor.Nama = input.Nama
	tutor.Email = input.Email
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return tutor, err
	}
	tutor.Password = string(password)

	newTutor, err := s.repository.Save(tutor)
	if err != nil {
		return newTutor, err
	}

	return newTutor, nil

}

func (s *tutorService) GetTutorByID(ID int) (tutors.Tutor, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found with that ID")
	}
	return user, nil

}

func (s *tutorService) Login(input tutors.LogisUserInput) (tutors.Tutor, error) {
	email := input.Email
	password := input.Password

	tutor, err := s.repository.FindByEmail(email)
	if err != nil {
		return tutor, nil
	}

	if tutor.ID == 0 {
		return tutor, errors.New("no tutor found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(tutor.Password), []byte(password))
	if err != nil {
		return tutor, err
	}

	return tutor, nil

}

func (s *tutorService) UpdateTutor(inputID tutors.GetTutorInput, inputData tutors.CreateTutorInput) (tutors.Tutor, error) {
	tutor, err := s.repository.FindByID(inputID.ID)

	if err != nil {
		return tutor, err
	}

	tutor.Nama = inputData.Nama
	tutor.MasaKerja = inputData.MasaKerja
	tutor.SitusWeb = inputData.SitusWeb
	tutor.Kompetensi = inputData.Kompetensi
	tutor.Pekerjaan = inputData.Pekerjaan
	tutor.TopikDiminati = inputData.TopikDiminati

	updateTutor, err := s.repository.Update(tutor)
	if err != nil {
		return updateTutor, err
	}

	return updateTutor, nil

}
