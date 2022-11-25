package utils

type ValidationErr struct {
	Err string
}

func (e ValidationErr) Error() string {
	return e.Err
}
