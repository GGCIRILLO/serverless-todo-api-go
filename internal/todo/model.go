package todo

type STATUS string

const (
	TODO        STATUS = "TODO"
	IN_PROGRESS STATUS = "IN_PROGRESS"
	DONE        STATUS = "DONE"
)

type ItemToDo struct {
	Pk          string  `json:"pk"` // partition key
	Sk          string  `json:"sk"` // sort key
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Status      STATUS  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func ValidateStatus(status string) bool {
	switch status {
	case string(TODO), string(IN_PROGRESS), string(DONE):
		return true
	default:
		return false
	}
}

func ValidateTitle(title string) bool {
	return len(title) > 0
}

func ValidateDescription(description *string) bool {
	if description == nil {
		return true
	}
	return len(*description) > 0
}

func ValidateItem(item ItemToDo) bool {
	return ValidateTitle(item.Title) && ValidateStatus(string(item.Status)) && ValidateDescription(item.Description)
}
