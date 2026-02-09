package goal

type Status string

const (
	NotStartedStatus Status = "not-started"
	InProgressStatus Status = "in-progress"
	CompletedStatus  Status = "completed"
	UnknownStatus    Status = "unknown"
)

func IsValidStatus(status Status) bool {
	switch status {
	case NotStartedStatus, InProgressStatus, CompletedStatus:
		return true
	default:
		return false
	}
}

func GetAllStatuses() []Status {
	return []Status{NotStartedStatus, InProgressStatus, CompletedStatus}
}

func StatusToString(status Status) string {
	return string(status)
}

func StringToStatus(statusStr string) (Status, bool) {
	status := Status(statusStr)
	if IsValidStatus(status) {
		return status, true
	}
	return UnknownStatus, false
}
