package code

const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 101001

	// ErrUserAlreadyExists - 400: User already exists.
	ErrUserAlreadyExists

	// ErrUserPasswordIncorrect - 400: User password is incorrect.
	ErrUserPasswordIncorrect

	// ErrSmsSend - 400: Send sms error.
	ErrSmsSend

	// ErrJWTDeploy - 500: JWT deploy error.
	ErrJWTDeploy

	// ErrCodeNotExist - 400: Sms code incorrect or expired.
	ErrCodeNotExist

	// ErrCodeExpired - 400: Sms code expired.
	ErrCodeExpired

	// ErrJWTReadFiled - 500: JWT read field error.
	ErrJWTReadFiled

	// ErrInvalidPrivKey - 500: that the given private key is invalid.
	ErrInvalidPrivKey

	// ErrFailedTokenCreation - 500: indicates JWT Token failed to create, reason unknown.
	ErrFailedTokenCreation
)
