package customerror

type CustomError struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
	Tip     string `json:"tip,omitempty"`
}

func New(msg string, code string, tip string) CustomError {
	return CustomError{Message: msg, Code: code, Tip: tip}
}

func NewByErr(err error, code string, tip string) CustomError {
	return CustomError{Message: err.Error(), Code: code, Tip: tip}
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	if custom, ok := err.(CustomError); ok {
		return custom
	}

	return CustomError{
		Message: err.Error(),
		Code:    "UNKNOWN",
		Tip:     "Unknown error",
	}
}

func Unwrap(err error) CustomError {
	if err == nil {
		return CustomError{}
	}

	if custom, ok := err.(CustomError); ok {
		return custom
	}

	return CustomError{
		Message: err.Error(),
		Code:    "UNKNOWN",
		Tip:     "Unknown error",
	}
}

func (c CustomError) Error() string {
	return c.Message
}
