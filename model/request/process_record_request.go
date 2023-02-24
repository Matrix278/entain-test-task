package model

import (
	"github.com/entain-test-task/model"
	"github.com/go-openapi/strfmt"
)

type ProcessRecordRequest struct {
	TransactionID strfmt.UUID4      `json:"transaction_id" validate:"required,uuid4"`
	Amount        float64           `json:"amount" validate:"required,min=0"`
	State         model.RecordState `json:"state" validate:"required,record_state`
}
