package errs

type RespondError struct {
	StatusCode int
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

func NewRespondError(StatusCode, Code int, Msg string) error {
	return &RespondError{}
}

func (e *RespondError) Error() string {
	return e.Msg
}
