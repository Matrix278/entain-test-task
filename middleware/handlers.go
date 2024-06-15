package middleware

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

type Handler struct {
	repository *repository.Store
}

func NewHandler(repository *repository.Store) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (handler *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.repository.GetAllUsers()
	if err != nil {
		log.Printf("unable to get all users. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, responsemodel.GetAllUsersResponse{
		Users: users,
	})
}

func (handler *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	user, err := handler.repository.GetUser(strfmt.UUID4(userID))
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

func (handler *Handler) GetAllTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	transactions, err := handler.repository.GetAllTransactionsByUserID(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("unable to get transactions. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, responsemodel.GetAllTransactionsByUserIDResponse{
		Transactions: transactions,
	})
}

func (handler *Handler) ProcessRecord(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	var processRecordRequest requestmodel.ProcessRecordRequest
	err := json.NewDecoder(r.Body).Decode(&processRecordRequest)
	if err != nil {
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

	err = handler.repository.ProcessRecord(strfmt.UUID4(userID), processRecordRequest)
	if err != nil {
		log.Printf("unable to process record. %v", err)

		if err.Error() == repository.ErrUserNotFound().Error() {
			StatusUnprocessableEntity(w, repository.ErrUserNotFound().Error())
			return
		}

		if err.Error() == repository.ErrInsufficientBalance().Error() {
			StatusUnprocessableEntity(w, repository.ErrInsufficientBalance().Error())
			return
		}

		if err.Error() == repository.ErrTransactionAlreadyExists().Error() {
			StatusUnprocessableEntity(w, repository.ErrTransactionAlreadyExists().Error())
			return
		}

		StatusInternalServerError(w)
		return
	}

	StatusOK(w, responsemodel.ProcessRecordResponse{
		Message: "OK",
	})
}

func (handler *Handler) CancelLatestOddTransactionRecords(numberOfRecords int) {
	transactions, err := handler.repository.GetLatestOddRecordTransactions(numberOfRecords)
	if err != nil {
		log.Printf("unable to get latest odd records. %v", err)
		return
	}

	if len(transactions) == 0 {
		log.Printf("no records to cancel")
		return
	}

	for _, transaction := range transactions {
		if err = handler.repository.CancelTransactionRecord(transaction); err != nil {
			log.Printf("unable to cancel transaction record. %v", err)
			return
		}
	}

	log.Printf("successfully cancelled %d records", len(transactions))
}
