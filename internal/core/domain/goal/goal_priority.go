package goal

type Priority string

const (
	Low             Priority = "low"
	Medium          Priority = "medium"
	High            Priority = "high"
	UnknownPriority Priority = "unknown"
)

func IsValidPriority(priority Priority) bool {
	switch priority {
	case Low, Medium, High:
		return true
	default:
		return false
	}
}

func GetAllPriorities() []Priority {
	return []Priority{Low, Medium, High}
}

func PriorityToString(priority Priority) string {
	return string(priority)
}

func StringToPriority(priorityStr string) (Priority, bool) {
	priority := Priority(priorityStr)
	if IsValidPriority(priority) {
		return priority, true
	}
	return UnknownPriority, false
}
