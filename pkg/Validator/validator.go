package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/mhvn092/movie-go/pkg/exception"
)

// Validation function using reflection
func validateInterface(s interface{}) []string {
	var errors []string
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		errors = []string{"invalid interface passed as input"}
		return errors
	}
	typ := val.Type()

	// Loop through struct fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		tag := field.Tag.Get("validate")

		if tag == "" {
			continue // Skip fields without validation tags
		}

		// Split validation rules
		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)

			if rule == "required" && isEmpty(value) {
				errors = append(errors, field.Name+" is required")
			}

			if rule == "is_string" && field.Type.Kind() != reflect.String {
				errors = append(errors, field.Name+" must be a string")
			}

			if rule == "is_email" && !isValidEmail(value.String()) {
				errors = append(errors, field.Name+" must be a valid email")
			}

			if rule == "is_strong_password" && !isStrongPassword(value.String()) {
				errors = append(errors, field.Name+" you should choose a strong password")
			}

			if rule == "is_phone_number" && !isValidPhoneNumber(value.String()) {
				errors = append(errors, field.Name+" phone number is not valid")
			}

			if strings.HasPrefix(rule, "min_len=") {
				minLen := parseMinLen(rule)
				if len(value.String()) < minLen {
					errors = append(
						errors,
						field.Name+
							fmt.Sprintf(" must be at least %d characters long", minLen),
					)
				}
			}
		}
	}

	return errors
}

// Helper: Check if a value is empty
func isEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return strings.TrimSpace(value.String()) == ""
	case reflect.Ptr:
		return value.IsNil()
	default:
		return false
	}
}

// Helper: Check if a string is a valid email
func isValidEmail(email string) bool {
	regexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regexp.MatchString(email)
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	lowercase := regexp.MustCompile(`[a-z]`)
	uppercase := regexp.MustCompile(`[A-Z]`)
	digit := regexp.MustCompile(`\d`)
	special := regexp.MustCompile(`[@$!%*?&]`)

	return lowercase.MatchString(password) &&
		uppercase.MatchString(password) &&
		digit.MatchString(password) &&
		special.MatchString(password)
}

func isValidPhoneNumber(phoneNumber string) bool {
	regexp := regexp.MustCompile(`^(09\d{9}|\+989\d{9}|0\d{2,3}\d{7,8}|\+98\d{2,3}\d{7,8}|\d{10})$`)
	return regexp.MatchString(phoneNumber)
}

// Helper: Parse min length rule
func parseMinLen(rule string) int {
	var length int
	fmt.Sscanf(rule, "min_len=%d", &length)
	return length
}

func JsonBodyHasErrors(req *http.Request, w http.ResponseWriter, payload interface{}) bool {
	// Read and check the request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		exception.HttpError(err, w, "Failed to read request body", http.StatusInternalServerError)
		return true
	}
	if len(body) == 0 {
		exception.HttpError(
			errors.New("Validation Error"),
			w,
			"Empty request body",
			http.StatusBadRequest,
		)
		return true
	}

	err = json.Unmarshal(body, payload)
	if err != nil {
		exception.HttpError(err, w, "Invalid JSON payload", http.StatusBadRequest)
		return true
	}

	// Validate the payload
	validationErrors := validateInterface(payload)
	if validationErrors != nil {

		exception.HttpError(
			errors.New(strings.Join(validationErrors, " ,")),
			w,
			"Invalid Input sent",
			http.StatusBadRequest,
		)
		return true
	}

	return false
}
