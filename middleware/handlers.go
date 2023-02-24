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
		log.Printf("Unable to get all users. %v", err)
		StatusInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsemodel.GetAllUsersResponse{
		Users: users,
	})
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	user, err := repository.GetUser(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("Unable to get user. %v", err)
		StatusInternalServerError(w)
		return
	}

	if user == nil {
		StatusUnprocessableEntity(w, "user not found")
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

func GetAllTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	transactions, err := repository.GetAllTransactionsByUserID(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("Unable to get transactions. %v", err)
		StatusInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsemodel.GetAllTransactionsByUserIDResponse{
		Transactions: transactions,
	})
}

func ProcessRecord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var processRecordRequest requestmodel.ProcessRecordRequest

	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&processRecordRequest)
	if err != nil {
		typeErr, ok := err.(*json.UnmarshalTypeError)
		if ok {
			StatusBadRequestWithErrors(w, "validation error", []error{
				errors.New(typeErr.Field + " is not is not within allowed range"),
			})
			return
		}

		log.Printf("Unable to decode the request body. %v", err)
		StatusInternalServerError(w)
		return
	}

	validationErrors := model.ValidateRequest(processRecordRequest)
	if validationErrors != nil {
		log.Printf("Unable to validate request body. %v", validationErrors)
		StatusBadRequestWithErrors(w, "validation error", validationErrors)
		return
	}

	err = repository.ProcessRecord(strfmt.UUID4(userID), processRecordRequest)
	if err != nil {
		log.Printf("Unable to process record. %v", err)

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

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responsemodel.ProcessRecordResponse{
		Message: "OK",
	})
}
