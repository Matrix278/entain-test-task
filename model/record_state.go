package model

type RecordState string

const (
	RecordStateWin  RecordState = "win"
	RecordStateLose RecordState = "lose"
)

var recordStateValues = []RecordState{
	RecordStateWin,
	RecordStateLose,
}
