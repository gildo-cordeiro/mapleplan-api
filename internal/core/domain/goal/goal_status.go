package goal

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	UnknownStatus    Status = "unknown"
)

func IsValidStatus(status Status) bool {
	switch status {
	case StatusPending, StatusInProgress, StatusCompleted:
		return true
	default:
		return false
	}
}

func GetAllStatuses() []Status {
	return []Status{StatusPending, StatusInProgress, StatusCompleted}
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
