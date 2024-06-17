package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/entain-test-task/configuration"
	"github.com/entain-test-task/model"
	requestmodel "github.com/entain-test-task/model/request"
	"github.com/entain-test-task/service"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
)

type Transaction struct {
	cfg                *configuration.Config
	transactionService service.ITransaction
	validator          *model.Validator
}

func NewTransaction(
	cfg *configuration.Config,
	transactionService service.ITransaction,
) *Transaction {
	return &Transaction{
		cfg:                cfg,
		transactionService: transactionService,
		validator:          model.NewValidator(),
	}
}

func (controller *Transaction) GetAllTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]

	if !strfmt.IsUUID4(userID) {
		StatusBadRequest(w, "user_id is not a valid UUID4")
		return
	}

	getAllTransactionsByUserIDResponse, err := controller.transactionService.GetAllTransactionsByUserID(strfmt.UUID4(userID))
	if err != nil {
		log.Printf("getting all transactions by user ID failed. %v", err)
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
		handleDecodeError(err, w)
		return
	}

	validationErrors := controller.validator.ValidateRequest(processRecordRequest)
	if validationErrors != nil {
		log.Printf("validating request failed. %v", validationErrors)
		StatusBadRequestWithErrors(w, "validation error", validationErrors)
		return
	}

	processRecordResponse, err := controller.transactionService.ProcessRecord(strfmt.UUID4(userID), processRecordRequest)
	if err != nil {
		handleProcessRecordError(w, err)
		return
	}

	StatusOK(w, processRecordResponse)
}

func (controller *Transaction) CancelLatestOddTransactionRecords() {
	controller.transactionService.CancelLatestOddTransactionRecords(controller.cfg.NumberOfLatestRecordsForCancelling)
}

func handleProcessRecordError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case model.ErrUserNotFound().Error():
		StatusUnprocessableEntity(w, model.ErrUserNotFound().Error())
	case model.ErrInsufficientBalance().Error():
		StatusUnprocessableEntity(w, model.ErrInsufficientBalance().Error())
	case model.ErrTransactionAlreadyExists().Error():
		StatusUnprocessableEntity(w, model.ErrTransactionAlreadyExists().Error())
	default:
		log.Printf("%s. %v", "processing record failed", err)
		StatusInternalServerError(w)
	}
}

func handleDecodeError(err error, w http.ResponseWriter) {
	if err == nil {
		return
	}

	handleUnmarshalTypeError(err, w)

	log.Printf("decoding request failed. %v", err)
	StatusInternalServerError(w)
}

func handleUnmarshalTypeError(err error, w http.ResponseWriter) {
	typeErr, ok := err.(*json.UnmarshalTypeError)
	if !ok {
		return
	}

	StatusBadRequestWithErrors(w, "validation error", []error{
		errors.New(typeErr.Field + " is not within allowed range"),
	})
}
