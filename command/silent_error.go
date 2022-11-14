package command

type SilentError struct{}

func (e *SilentError) Error() string {
	return ""
}

func IsSilentError(err error) bool {
	_, ok := err.(*SilentError)
	return ok
}
