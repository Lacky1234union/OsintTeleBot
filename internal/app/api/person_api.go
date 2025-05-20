package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
)

type PersonService interface {
	RegisterUser(ctx context.Context, user models.Person) error
	FindUserByEmail(ctx context.Context, email string) (models.Person, error)
	FindUserByName(ctx context.Context, name string) (models.Person, error)
	FindUserByPhone(ctx context.Context, phone string) (models.Person, error)
}

type PersonAPI struct {
	Service PersonService
}

func NewPersonAPI(service PersonService) *PersonAPI {
	return &PersonAPI{Service: service}
}

func (api *PersonAPI) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := api.Service.RegisterUser(r.Context(), person)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (api *PersonAPI) FindUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	person, err := api.Service.FindUserByEmail(r.Context(), email)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func (api *PersonAPI) FindUserByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	person, err := api.Service.FindUserByName(r.Context(), name)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func (api *PersonAPI) FindUserByPhoneHandler(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	person, err := api.Service.FindUserByPhone(r.Context(), phone)
	if err != nil {
		writeAPIError(w, err)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func writeAPIError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch err {
	case errs.ErrBadData, errs.ErrNilContext:
		status = http.StatusBadRequest
	case errs.ErrPersonNotFound, errs.ErrEmailNotFound, errs.ErrPhoneNotFound:
		status = http.StatusNotFound
	case errs.ErrAlreadyExists:
		status = http.StatusConflict
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
