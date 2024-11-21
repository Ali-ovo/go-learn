package code

const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100501

	// ErrUserAlreadyExists - 400: User already exists.
	ErrUserAlreadyExists

	// ErrUserPasswordIncorrect - 400: User password is incorrect.
	ErrUserPasswordIncorrect

	// ErrSmsSend - 400: Send sms error.
	ErrSmsSend
)
