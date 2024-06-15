package service

import (
	"log"

	requestmodel "github.com/entain-test-task/model/request"
	responsemodel "github.com/entain-test-task/model/response"
	"github.com/entain-test-task/repository"
	"github.com/go-openapi/strfmt"

	_ "github.com/lib/pq"
)

type Transaction struct {
	repository *repository.Store
}

func NewTransaction(repository *repository.Store) *Transaction {
	return &Transaction{
		repository: repository,
	}
}

func (service *Transaction) GetAllTransactionsByUserID(userID strfmt.UUID4) (*responsemodel.GetAllTransactionsByUserIDResponse, error) {
	transactions, err := service.repository.GetAllTransactionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &responsemodel.GetAllTransactionsByUserIDResponse{
		Transactions: transactions,
	}, nil
}

func (service *Transaction) ProcessRecord(userID strfmt.UUID4, processRecordRequest requestmodel.ProcessRecordRequest) (*responsemodel.ProcessRecordResponse, error) {
	if err := service.repository.ProcessRecord(userID, processRecordRequest); err != nil {
		return nil, err
	}

	return &responsemodel.ProcessRecordResponse{
		Message: "OK",
	}, nil
}

func (service *Transaction) CancelLatestOddTransactionRecords(numberOfRecords int) {
	transactions, err := service.repository.GetLatestOddRecordTransactions(numberOfRecords)
	if err != nil {
		log.Printf("unable to get latest odd records. %v", err)
		return
	}

	if len(transactions) == 0 {
		log.Printf("no records to cancel")
		return
	}

	for _, transaction := range transactions {
		if err = service.repository.CancelTransactionRecord(transaction); err != nil {
			log.Printf("unable to cancel transaction record. %v", err)
			return
		}
	}

	log.Printf("successfully cancelled %d records", len(transactions))
}
