package servicenow

type ChangeTask struct {
	taskNumber      string
	expectedStart   string
	expectedEnd     string
	taskDescription string
	url             string
}
