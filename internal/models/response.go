package models

type IndexResponse struct {
	Data  []*Task `json:"data"`
	Total int64   `json:"total"`
}

type DeleteResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type TaskResponse struct {
	Data *Task `json:"data,omitempty"`
}
