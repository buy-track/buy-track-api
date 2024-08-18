package errors

type NotFoundError struct {
}

func (n NotFoundError) Error() string {
	return "Data not found"
}

type DuplicateError struct {
	Data string
}

func (n DuplicateError) Error() string {
	return "Duplicate for " + n.Data + " data"
}
