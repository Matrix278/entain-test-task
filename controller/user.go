package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/entain-test-task/model"
	requestmodel "github.com/entain-test-task/model/request"
	responsemodel "github.com/entain-test-task/model/response"
	"github.com/entain-test-task/repository"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type User struct {
	repository *repository.Store
}

func NewUser(repository *repository.Store) *User {
	return &User{
		repository: repository,
	}
}

func (controller *User) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := controller.repository.GetAllUsers()
	if err != nil {
		log.Printf("unable to get all users. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, responsemodel.GetAllUsersResponse{
		Users: users,
	})
}

func (controller *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	user, err := controller.repository.GetUser(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("unable to get user. %v", err)

		if err.Error() == repository.ErrUserNotFound().Error() {
			StatusUnprocessableEntity(w, repository.ErrUserNotFound().Error())
			return
		}

		StatusInternalServerError(w)
		return
	}

	StatusOK(w, user)
}

func (controller *User) GetAllTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	transactions, err := controller.repository.GetAllTransactionsByUserID(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("unable to get transactions. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, responsemodel.GetAllTransactionsByUserIDResponse{
		Transactions: transactions,
	})
}

func (controller *User) ProcessRecord(w http.ResponseWriter, r *http.Request) {
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

	if err := controller.repository.ProcessRecord(strfmt.UUID4(userID), processRecordRequest); err != nil {
		handleProcessRecordError(w, err)
		return
	}

	StatusOK(w, responsemodel.ProcessRecordResponse{
		Message: "OK",
	})
}

func (controller *User) CancelLatestOddTransactionRecords(numberOfRecords int) {
	transactions, err := controller.repository.GetLatestOddRecordTransactions(numberOfRecords)
	if err != nil {
		log.Printf("unable to get latest odd records. %v", err)
		return
	}

	if len(transactions) == 0 {
		log.Printf("no records to cancel")
		return
	}

	for _, transaction := range transactions {
		if err = controller.repository.CancelTransactionRecord(transaction); err != nil {
			log.Printf("unable to cancel transaction record. %v", err)
			return
		}
	}

	log.Printf("successfully cancelled %d records", len(transactions))
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
