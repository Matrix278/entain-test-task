package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/entain-test-task/model"
	requestmodel "github.com/entain-test-task/model/request"
	"github.com/entain-test-task/repository"
	"github.com/entain-test-task/service"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Transaction struct {
	service *service.Transaction
}

func NewTransaction(service *service.Transaction) *Transaction {
	return &Transaction{
		service: service,
	}
}

func (controller *Transaction) GetAllTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	getAllTransactionsByUserIDResponse, err := controller.service.GetAllTransactionsByUserID(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("unable to get transactions. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, getAllTransactionsByUserIDResponse)
}

func (controller *Transaction) ProcessRecord(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	var processRecordRequest requestmodel.ProcessRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&processRecordRequest); err != nil {
		typeErr, ok := err.(*json.UnmarshalTypeError)
		if ok {
			StatusBadRequestWithErrors(w, "validation error", []error{
				errors.New(typeErr.Field + " is not is not within allowed range"),
			})
			return
		}

		log.Printf("unable to decode the request body. %v", err)
		StatusInternalServerError(w)
		return
	}

	validationErrors := model.ValidateRequest(processRecordRequest)
	if validationErrors != nil {
		log.Printf("unable to validate request body. %v", validationErrors)
		StatusBadRequestWithErrors(w, "validation error", validationErrors)
		return
	}

	processRecordResponse, err := controller.service.ProcessRecord(strfmt.UUID4(userID), processRecordRequest)
	if err != nil {
		handleProcessRecordError(w, err)
		return
	}

	StatusOK(w, processRecordResponse)
}

func handleProcessRecordError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case repository.ErrUserNotFound().Error():
		StatusUnprocessableEntity(w, repository.ErrUserNotFound().Error())
	case repository.ErrInsufficientBalance().Error():
		StatusUnprocessableEntity(w, repository.ErrInsufficientBalance().Error())
	case repository.ErrTransactionAlreadyExists().Error():
		StatusUnprocessableEntity(w, repository.ErrTransactionAlreadyExists().Error())
	default:
		log.Printf("%s. %v", "unable to process record", err)
		StatusInternalServerError(w)
	}
}
