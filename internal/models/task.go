package models

type Task struct {
	ID int `gorm:"primaryKey" json:"id"`
	TaskRequest
}

func (Task) TableName() string {
	return "tasks"
}

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Deleted     bool   `json:"-"`
}
