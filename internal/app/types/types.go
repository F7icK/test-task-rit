package types

const (
	ActionType          = "action"
	ConditionType       = "condition"
	ConditionActionType = "condition-action"

	CreateFile             = "create file"
	ChangeFileName         = "change file name"
	RemoveFile             = "remove file"
	GetFileCreationTime    = "get file creation time"
	WritingLineToFile      = "writing line to file"
	CreatedAtTimeCondition = "created at time condition"
)

type DataJSON struct {
	Actions Actions `json:"actions"`
}

type Actions []*Action
type Action struct {
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Result     string   `json:"result"`
	Parameters []string `json:"parameters"`
}
