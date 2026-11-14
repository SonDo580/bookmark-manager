package response

import (
	"errors"
	"fmt"

	"github.com/SonDo580/bookmark-manager/internal/dto"
	"github.com/go-playground/validator/v10"
)

func InputFieldError(err error) dto.ErrorResponse {
	msg := "Invalid Input"
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return dto.ErrorResponse{Message: msg}
	}

	var errs []string
	for _, fieldErr := range err.(validator.ValidationErrors) {
		msg := fmt.Sprintf("%s is invalid (%s)", fieldErr.Field(), fieldErr.Tag())
		errs = append(errs, msg)
	}

	return dto.ErrorResponse{Message: msg, Details: errs}
}
