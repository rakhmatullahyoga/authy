package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/validator.v2"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Cause   string      `json:"cause,omitempty"`
	User    *User       `json:"user,omitempty"`
	Recipe  *UserUpdate `json:"recipe,omitempty"`
}

var (
	deleteUserService   = deleteUser
	getUserService      = getUser
	registerUserService = registerUser
	updateUserService   = updateUser
)

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Post("/signup", registerUserHandler)
	r.With(BasicAuth).Get("/users/{userID}", getUserHandler)
	r.With(BasicAuth).Patch("/users/{userID}", updateUserHandler)
	r.With(BasicAuth).Post("/close", deleteUserHandler)
	return r
}

func writeResponse(w http.ResponseWriter, res Response, status int) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(res)
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var params UserCredential
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeResponse(w, Response{
			Message: "Account creation failed",
			Cause:   "required user_id and password",
		}, http.StatusBadRequest)
		return
	}

	if err := validator.Validate(params); err != nil {
		writeResponse(w, Response{
			Message: "Account creation failed",
			Cause:   "invalid user_id and password format",
		}, http.StatusBadRequest)
		return
	}

	user, err := registerUserService(r.Context(), params)
	if err != nil {
		writeResponse(w, Response{
			Message: "Account creation failed",
			Cause:   err.Error(),
		}, http.StatusBadRequest)
		return
	}

	writeResponse(w, Response{
		Message: "Account successfully created",
		User:    &user,
	}, http.StatusOK)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	user, err := getUserService(r.Context(), userID)
	if err != nil {
		writeResponse(w, Response{
			Message: err.Error(),
		}, http.StatusNotFound)
		return
	}

	writeResponse(w, Response{
		Message: "User details by user_id",
		User:    &user,
	}, http.StatusOK)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var params UserUpdate
	userID := chi.URLParam(r, "userID")
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeResponse(w, Response{
			Message: "User updation failed",
			Cause:   "required nickname or comment",
		}, http.StatusBadRequest)
		return
	}

	if params.Nickname == "" && params.Comment == "" {
		writeResponse(w, Response{
			Message: "Account creation failed",
			Cause:   "required nickname or comment",
		}, http.StatusBadRequest)
		return
	}

	if params.UserID != "" || params.Password != "" {
		writeResponse(w, Response{
			Message: "Account creation failed",
			Cause:   "not updateable user_id and password",
		}, http.StatusBadRequest)
		return
	}

	if err := validator.Validate(params); err != nil {
		writeResponse(w, Response{
			Message: "Account creation failed",
			Cause:   "invalid user_id and password format",
		}, http.StatusBadRequest)
		return
	}

	_, status, err := updateUserService(r.Context(), userID, params)
	if err != nil {
		writeResponse(w, Response{
			Message: err.Error(),
			Recipe:  &params,
		}, status)
		return
	}

	writeResponse(w, Response{
		Message: "User successfully updated",
	}, http.StatusOK)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	deleteUserService(r.Context())
	writeResponse(w, Response{
		Message: "Account and user successfully removed",
	}, http.StatusOK)
}
