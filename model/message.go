package model

type SendMessageRequest struct {
	Message string `json:"message"`
}

type Update struct {
	CreatedAt int64
	Message   string
}

func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}
