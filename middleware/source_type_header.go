package middleware

import (
	"net/http"

	"github.com/entain-test-task/model"
)

type SourceType string

const (
	SourceTypeGame    SourceType = "game"
	SourceTypeServer  SourceType = "server"
	SourceTypePayment SourceType = "payment"
)

var sourceTypeValues = []string{
	string(SourceTypeGame),
	string(SourceTypeServer),
	string(SourceTypePayment),
}

func ValidateSourceHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if model.IsArrayContainsValue(sourceTypeValues, r.Header.Get("Source-Type")) {
			StatusBadRequest(w, "Source-Type header is not valid")
			return
		}

		next.ServeHTTP(w, r)
	})
}
