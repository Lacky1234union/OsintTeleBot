package errs

var (
	ErrNotFound      = New("not found")
	ErrBadData       = New("bad input data")
	ErrAlreadyExists = New("already exists")
)
