package errs

var (
	// General errors
	ErrNotFound      = New("not found")
	ErrBadData       = New("bad input data")
	ErrAlreadyExists = New("already exists")
	ErrNilContext    = New("context is nil")

	// Database errors
	ErrDatabaseConnection = New("database connection error")
	ErrDatabaseQuery      = New("database query error")
	ErrDatabaseScan       = New("database scan error")

	// Person repository specific errors
	ErrPersonNotFound    = New("person not found")
	ErrPersonCreate      = New("failed to create person")
	ErrPersonInvalidData = New("invalid person data")
	ErrEmailNotFound     = New("email not found")
	ErrPhoneNotFound     = New("phone not found")
)
