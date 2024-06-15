package enum

type RecordState string

const (
	RecordStateWin  RecordState = "win"
	RecordStateLose RecordState = "lose"
)

var RecordStates = []string{
	string(RecordStateWin),
	string(RecordStateLose),
}
