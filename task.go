package mytasks

import "errors"

type TaskList struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}
type UpdateTaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateTaskInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("update data has no values")
	}
	return nil
}
