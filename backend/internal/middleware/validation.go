package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	Surname   string `json:"surname"`
	DOB       string `json:"dob"`
}

func ValidateRegister(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		
		log.Printf("Received registration request: %s", string(body))
		
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" || req.FirstName == "" || req.Surname == "" || req.DOB == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		errors := make(map[string]string)

		if !validatePassword(req.Password) {
			errors["password"] = "Password must be at least 8 characters and contain uppercase, lowercase, number, and special character"
		}

		if !validateEmail(req.Email) {
			errors["email"] = "Invalid email format"
		}

		if !validateDOB(req.DOB) {
			errors["dob"] = "Invalid date of birth or user must be at least 13 years old"
		}

		if len(errors) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": errors,
			})
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	}
}

func validatePassword(password string) bool {
	hasMinLen := len(password) >= 8
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func validateDOB(dob string) bool {
	parsedDOB, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return false
	}

	minAge := 13
	now := time.Now()
	age := now.Year() - parsedDOB.Year()
	
	return age >= minAge && parsedDOB.Before(now)
}
