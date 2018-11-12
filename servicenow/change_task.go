package servicenow

type ChangeTask struct {
	changeNumber    string
	taskNumber      string
	expectedStart   string
	expectedEnd     string
	taskDescription string
	url             string
}
