package desec

// NotFound Not found error.
type NotFound struct {
	Detail string `json:"detail"`
}

func (n NotFound) Error() string {
	return n.Detail
}
