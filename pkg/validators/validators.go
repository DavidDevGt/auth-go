package validators

import (
	"errors"
	"regexp"
	"strings"
)

var (
	emailRegex     = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	uppercaseRegex = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex = regexp.MustCompile(`[a-z]`)
	numberRegex    = regexp.MustCompile(`[0-9]`)
	specialRegex   = regexp.MustCompile(`[\W_]`)
)

func ValidateRegister(name, email, password string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	if len(name) < 2 {
		return errors.New("name is too short")
	}
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	if err := ValidatePassword(password, name, email); err != nil {
		return err
	}
	return nil
}

func ValidatePassword(password, name, email string) error {
	pwd := strings.TrimSpace(password)
	if pwd == "" {
		return errors.New("password is required")
	}
	if len(pwd) < 9 {
		return errors.New("password must be at least 9 characters long")
	}
	if strings.Contains(pwd, " ") {
		return errors.New("password cannot contain spaces")
	}
	if !lowercaseRegex.MatchString(pwd) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !uppercaseRegex.MatchString(pwd) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !numberRegex.MatchString(pwd) {
		return errors.New("password must contain at least one digit")
	}
	if !specialRegex.MatchString(pwd) {
		return errors.New("password must contain at least one special character")
	}
	if name != "" && strings.Contains(strings.ToLower(pwd), strings.ToLower(name)) {
		return errors.New("password should not contain your name")
	}
	if email != "" {
		localPart := strings.ToLower(strings.Split(email, "@")[0])
		if strings.Contains(strings.ToLower(pwd), localPart) {
			return errors.New("password should not contain your email username")
		}
	}
	return nil
}

func ValidateLogin(email, password string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	if strings.TrimSpace(password) == "" {
		return errors.New("password is required")
	}
	return nil
}

func ValidateRefreshToken(token string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("refresh_token is required")
	}
	return nil
}
