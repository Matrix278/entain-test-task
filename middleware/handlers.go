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

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetAllUsers()
	if err != nil {
		log.Printf("unable to get all users. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, responsemodel.GetAllUsersResponse{
		Users: users,
	})
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	user, err := repository.GetUser(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("unable to get user. %v", err)
		StatusInternalServerError(w)
		return
	}

	if user == nil {
		StatusUnprocessableEntity(w, repository.ErrUserNotFound().Error())
		return
	}

	StatusOK(w, user)
}

func GetAllTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	transactions, err := repository.GetAllTransactionsByUserID(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("unable to get transactions. %v", err)
		StatusInternalServerError(w)
		return
	}

	StatusOK(w, responsemodel.GetAllTransactionsByUserIDResponse{
		Transactions: transactions,
	})
}

func ProcessRecord(w http.ResponseWriter, r *http.Request) {
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

	err = repository.ProcessRecord(strfmt.UUID4(userID), processRecordRequest)
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

func CancelLatestOddTransactionRecords(numberOfRecords int) {
	transactions, err := repository.GetLatestOddRecordTransactions(numberOfRecords)
	if err != nil {
		log.Printf("unable to get latest odd records. %v", err)
		return
	}

	if len(transactions) == 0 {
		log.Printf("no records to cancel")
		return
	}

	for _, transaction := range transactions {
		if err = repository.CancelTransactionRecord(transaction); err != nil {
			log.Printf("unable to cancel transaction record. %v", err)
			return
		}
	}
}
