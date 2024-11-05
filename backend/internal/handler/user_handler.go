package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"tech-test/backend/internal/domain"
	userInterface "tech-test/backend/internal/service/interfaces/user"
	"tech-test/backend/internal/utils"
	"github.com/gorilla/mux"
	"strconv"
	"tech-test/backend/internal/middleware"
)

type UserHandler struct {
	userService userInterface.UserService
}

func NewUserHandler(userService userInterface.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body domain.User true "User details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} domain.APIError
// @Router /api/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			"Invalid request body",
			err,
		))
		return
	}

	if err := h.userService.Register(r.Context(), &user); err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			err.Error(),
			err,
		))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "User created successfully",
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			"Invalid user ID",
			err,
		))
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), uint(userID))
	if err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusNotFound,
			domain.ErrCodeNotFound,
			"User not found",
			err,
		))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			"Invalid user ID",
			err,
		))
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			"Invalid request body",
			err,
		))
		return
	}

	if err := h.userService.UpdateUser(r.Context(), uint(userID), &user); err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusInternalServerError,
			domain.ErrCodeInternal,
			"Failed to update user",
			err,
		))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
    if !ok {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusUnauthorized,
            domain.ErrCodeAuthentication,
            "Invalid authentication",
            nil,
        ))
        return
    }

    user, err := h.userService.GetUserByID(r.Context(), userID)
    if err != nil {
        utils.RespondWithError(w, domain.NewAPIError(
            http.StatusNotFound,
            domain.ErrCodeNotFound,
            "User not found",
            err,
        ))
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			"Invalid user ID",
			err,
		))
		return
	}

	if err := h.userService.DeleteUser(r.Context(), uint(userID)); err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusInternalServerError,
			domain.ErrCodeInternal,
			"Failed to delete user",
			err,
		))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register handler called with method: %s", r.Method)
	
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			"Invalid request body",
			err,
		))
		return
	}
	
	log.Printf("Attempting to register user: %+v", user)

	if err := h.userService.Register(r.Context(), &user); err != nil {
		log.Printf("Error registering user: %v", err)
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			err.Error(),
			err,
		))
		return
	}

	log.Printf("User registered successfully")
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusBadRequest,
			domain.ErrCodeInvalidInput,
			"Invalid request body",
			err,
		))
		return
	}

	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusUnauthorized,
			domain.ErrCodeAuthentication,
			"Invalid credentials",
			err,
		))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user": user,
	})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users in the system
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} domain.User
// @Failure 500 {object} domain.APIError
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllUsers handler called")
	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		log.Printf("Error getting users: %v", err)
		utils.RespondWithError(w, domain.NewAPIError(
			http.StatusInternalServerError,
			domain.ErrCodeInternal,
			"Failed to fetch users",
			err,
		))
		return
	}

	log.Printf("Found %d users", len(users))
	utils.RespondWithJSON(w, http.StatusOK, users)
}
