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
	"time"

	"github.com/mhvn092/movie-go/pkg/exception"
)

// validateInterface validates a struct or slice of structs and returns a list of validation errors.
func validateInterface(s interface{}) []string {
	var errors []string
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Handle slices
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			// Recursively validate each element in the slice
			itemErrors := validateSingleStruct(item, fmt.Sprintf("[%d]", i))
			errors = append(errors, itemErrors...)
		}
		return errors
	}

	// Handle structs
	if val.Kind() != reflect.Struct {
		return []string{"invalid interface passed as input"}
	}

	return validateSingleStruct(val, "")
}

// validateSingleStruct validates a single struct and prefixes field names with context (e.g., "Staffs[0].").
func validateSingleStruct(val reflect.Value, prefix string) []string {
	var errors []string
	typ := val.Type()

	// Loop through struct fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		tag := field.Tag.Get("validate")
		fieldName := prefix + field.Name

		if tag == "" {
			continue // Skip fields without validation tags
		}

		// Handle nested slices
		if value.Kind() == reflect.Slice {
			for j := 0; j < value.Len(); j++ {
				item := value.Index(j)
				if item.Kind() == reflect.Struct {
					// Recursively validate each struct in the slice
					nestedErrors := validateSingleStruct(
						item,
						fmt.Sprintf("%s[%d].", fieldName, j),
					)
					errors = append(errors, nestedErrors...)
				}
			}
			continue
		}

		// Split validation rules
		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)

			if rule == "required" && isEmpty(value) {
				errors = append(errors, fieldName+" is required")
			}

			if rule == "is_string" && value.Type().Kind() != reflect.String {
				errors = append(errors, fieldName+" must be a string")
			}

			if rule == "required" && value.Kind() == reflect.Slice && value.Len() == 0 {
				errors = append(errors, fieldName+" is required and cannot be empty")
			}

			if rule == "is_int" && value.Type().Kind() != reflect.Int {
				errors = append(errors, fieldName+" must be an int")
			}

			if rule == "is_email" && !isValidEmail(value.String()) {
				errors = append(errors, fieldName+" must be a valid email")
			}

			if rule == "is_date_string" && !isValidDate(value.String()) {
				errors = append(
					errors,
					fieldName+" must be a valid date string with the format of 2025-07-01",
				)
			}

			if rule == "is_strong_password" && !isStrongPassword(value.String()) {
				errors = append(errors, fieldName+" you should choose a strong password")
			}

			if rule == "is_phone_number" && !isValidPhoneNumber(value.String()) {
				errors = append(errors, fieldName+" phone number is not valid")
			}

			if strings.HasPrefix(rule, "min_len=") {
				minLen := parseMinLen(rule)
				if len(value.String()) < minLen {
					errors = append(
						errors,
						fieldName+fmt.Sprintf(" must be at least %d characters long", minLen),
					)
				}
			}

			// Custom validation for ProductionYear
			if rule == "is_valid_year" && value.Type().Kind() == reflect.Int {
				if !isValidProductionYear(value.Int()) {
					errors = append(
						errors,
						fieldName+" must be a valid production year between 1888 and 2030",
					)
				}
			}
		}
	}

	return errors
}

// isValidProductionYear checks if the year is within a reasonable range.
func isValidProductionYear(year int64) bool {
	currentYear := time.Now().Year()
	return year >= 1888 && year <= int64(currentYear+5) // Allow up to 5 years in the future
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

func isValidDate(dateString string) bool {
	_, err := time.Parse("2006-01-02", dateString)
	return err == nil
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
