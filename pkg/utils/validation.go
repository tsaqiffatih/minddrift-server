package utils

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		log.Printf("Validation error in ValidateStruct: %T\n", err)
		return err
	}

	return nil
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	log.Printf("Error type in FormatValidationError: %T\n", err)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := strings.ToLower(fieldError.StructField())
			errorMessage := ""

			switch fieldError.Tag() {
			case "required":
				errorMessage = fieldName + " is required"
			case "email":
				errorMessage = "Invalid email format"
			case "min":
				errorMessage = fieldName + " must be at least " + fieldError.Param() + " characters long"
			case "max":
				errorMessage = fieldName + " must be at most " + fieldError.Param() + " characters long"
			case "oneof":
				errorMessage = fieldName + " must be one of " + fieldError.Param()
			case "password":
				errorMessage = "Invalid password format"
			default:
				errorMessage = fieldName + " is invalid"
			}

			errors[fieldName] = errorMessage
		}
	} else if err != nil {
		// Handle other types of errors, such as *errors.errorString
		errors["error"] = "Internal server error"
	}

	return errors
}

func FormatBindingError(err error) map[string]string {
	errors := make(map[string]string)

	log.Printf("Error type in FormatBinding: %T\n", err)
	log.Println(err.Error())
	errors["binding"] = err.Error()

	return errors
}
