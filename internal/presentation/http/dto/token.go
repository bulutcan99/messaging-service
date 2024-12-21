package dto

type (
	TokenRequest struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
	}

	TokenResponse struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		TokenType        string `json:"token_type"`
	}

	ErrorResponse struct {
		Error            string          `json:"error,omitempty"`
		ErrorDescription string          `json:"error_description,omitempty"`
		Errors           []DetailedError `json:"errors,omitempty"`
	}

	DetailedError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

var (
	Token400Error = ErrorResponse{
		Error:            "invalid_request",
		ErrorDescription: "Invalid username or password",
	}

	Token401Error = ErrorResponse{
		Errors: []DetailedError{
			{Code: "unauthorized", Message: "Invalid credentials"},
		},
	}

	Token500Error = ErrorResponse{
		Errors: []DetailedError{
			{Code: "internal_error", Message: "Internal server error"},
		},
	}
)
