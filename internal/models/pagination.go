package models

type ListFilter struct {
	Pagination
	//TaskFilter
}

type Pagination struct {
	Limit  int
	Offset int
}

//type TaskFilter struct {
//	Title       string
//	Description string
//	Completed   bool
//}
