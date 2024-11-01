
package errors

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

var (
    ErrInternalServer = &APIError{
        Code:    500,
        Message: "Internal server error",
    }
    ErrInvalidInput = &APIError{
        Code:    400,
        Message: "Invalid input",
    }
    ErrUserNotFound = &APIError{
        Code:    404,
        Message: "User not found",
    }
    ErrEmailTaken = &APIError{
        Code:    409,
        Message: "Email already taken",
    }
)
