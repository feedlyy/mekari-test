package domain

import (
	"errors"
	"regexp"
)

func IsValidEmail(email string) bool {
	// Define the regular expression pattern for a valid email address
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regular expression
	regex := regexp.MustCompile(emailPattern)

	return regex.MatchString(email)
}

func (e *Employee) Validate(args ...interface{}) error {
	var err error

	switch {
	case e.FirstName == "":
		err = errors.New("missing required field: first_name")
	case e.LastName == "":
		err = errors.New("missing required field: last_name")
	case e.Email == "":
		err = errors.New("missing required field: email")
	case args[0] == "":
		err = errors.New("missing required field: hire_date")
	}

	if err != nil {
		return err
	}

	if !IsValidEmail(e.Email) {
		return errors.New("please input valid email")
	}

	return nil
}
