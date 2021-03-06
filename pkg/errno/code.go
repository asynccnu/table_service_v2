package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	ErrNoTable    = &Errno{Code: 20000, Message: "Error get table."}
	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase   = &Errno{Code: 20002, Message: "Database error."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}

	// user errors
	ErrEncrypt              = &Errno{Code: 20101, Message: "Error occurred while encrypting the user password."}
	ErrUserNotFound         = &Errno{Code: 20102, Message: "The user was not found."}
	ErrAuthorizationInvalid = &Errno{Code: 20103, Message: "The Authorization header was invalid."}
	ErrPasswordIncorrect    = &Errno{Code: 20104, Message: "The password was incorrect."}

	// table errors
	ErrWeekConvert   = &Errno{Code: 20201, Message: "Week Integer can't be convert to week String"}
	ErrDeleteTable   = &Errno{Code: 20202, Message: "Error occurred while deleting user's table"}
	ErrDeleteXKTable = &Errno{Code: 20203, Message: "Can not delete tables gotten from XK"}
)
