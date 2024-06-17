package enum

type SourceType string

const (
	SourceTypeGame    SourceType = "game"
	SourceTypeServer  SourceType = "server"
	SourceTypePayment SourceType = "payment"
)

var SourceTypes = []string{
	string(SourceTypeGame),
	string(SourceTypeServer),
	string(SourceTypePayment),
}
