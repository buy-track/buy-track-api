package domains

type ListItem[T any] struct {
	Data   []T `json:"data"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type LengthListItem[T any] struct {
	Data       []T
	Limit      int
	Offset     int
	TotalCount int
}
