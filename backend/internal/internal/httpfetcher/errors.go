package httpfetcher

type NoRespError struct {
	cause error
}

func (e *NoRespError) Error() string {
	return "no response from server: " + e.cause.Error()
}

func (e *NoRespError) Unwrap() error {
	return e.cause
}
