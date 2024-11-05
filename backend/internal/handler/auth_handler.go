// internal/handler/auth.go
package handler

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "time"
    "tech-test/backend/internal/domain"
    userInterface "tech-test/backend/internal/service/interfaces/user"
    "tech-test/backend/internal/utils"
    "tech-test/backend/internal/middleware"
    "fmt"
    "errors"
)                       

type AuthHandler struct {
    userService userInterface.UserService
}

func NewAuthHandler(userService userInterface.UserService) *AuthHandler {
    return &AuthHandler{
        userService: userService,
    }
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var loginRequest struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading body: %v", err)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid request body",
            err,
        ))
        return
    }
    log.Printf("Received login request: %s", string(body))

    r.Body = io.NopCloser(bytes.NewBuffer(body))

    if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
        log.Printf("Error decoding login request: %v", err)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid request format",
            err,
        ))
        return
    }

    log.Printf("Attempting login for email: %s", loginRequest.Email)

    user, err := h.userService.GetUserByEmail(r.Context(), loginRequest.Email)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusNotFound,
            domain.ErrCodeNotFound,
            fmt.Sprintf("No account found with email: %s", loginRequest.Email),
            err,
        ))
        return
    }

    
    if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusUnauthorized,
            domain.ErrCodeAuthentication,
            "Incorrect password",
            errors.New("password mismatch"),
        ))
        return
    }

    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusInternalServerError,
            domain.ErrCodeInternal,
            "Error generating token",
            err,
        ))
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
        "token": token,
        "user":  user,
    })
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} domain.APIError "Invalid request body"
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // Read and log the request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading body: %v", err)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid request body",
            err,
        ))
        return
    }
    log.Printf("Received registration request: %s", string(body))

    
    r.Body = io.NopCloser(bytes.NewBuffer(body))

    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Error decoding registration request: %v", err)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid request format",
            err,
        ))
        return
    }

    log.Printf("Parsed registration request: %+v", req)

    dob, err := time.Parse("2006-01-02", req.DOB)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            "Invalid date format for DOB. Use YYYY-MM-DD",
            err,
        ))
        return
    }

    user := &domain.User{
        Email:     req.Email,
        Password:  req.Password,
        FirstName: req.FirstName,
        Surname:   req.Surname,
        DOB:       dob,
    }

    if err := h.userService.Register(r.Context(), user); err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusBadRequest,
            domain.ErrCodeInvalidInput,
            err.Error(),
            err,
        ))
        return
    }

    utils.RespondWithJSON(w, http.StatusCreated, map[string]string{
        "message": "User registered successfully",
    })
}

type LoginResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    User  struct {
        ID        uint   `json:"id" example:"1"`
        Email     string `json:"email" example:"john.doe@example.com"`
        FirstName string `json:"firstName" example:"John"`
        Surname   string `json:"surname" example:"Doe"`
    } `json:"user"`
}

type RegisterRequest struct {
    Email     string `json:"email" binding:"required"`
    Password  string `json:"password" binding:"required"`
    FirstName string `json:"firstName" binding:"required"`
    Surname   string `json:"surname" binding:"required"`
    DOB       string `json:"dob" binding:"required"`
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    log.Printf("GetCurrentUser handler called")
    
    userID, ok := middleware.GetUserIDFromContext(r.Context())
    if !ok {
        log.Printf("No user ID found in context")
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusUnauthorized,
            domain.ErrCodeAuthentication,
            "User not authenticated",
            nil,
        ))
        return
    }

    log.Printf("Found user ID in context: %d", userID)
    user, err := h.userService.GetUserByID(r.Context(), userID)
    if err != nil {
        log.Printf("Error getting user by ID: %v", err)
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusNotFound,
            domain.ErrCodeNotFound,
            "User not found",
            err,
        ))
        return
    }

    
    user.Password = ""

    log.Printf("Successfully retrieved user: %+v", user)
    utils.RespondWithJSON(w, http.StatusOK, user)
}
