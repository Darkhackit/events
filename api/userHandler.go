package api

import (
	"encoding/json"
	"github.com/Darkhackit/events/dto"
	"github.com/Darkhackit/events/service"
	"io"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	u, err := uh.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	WriteResponse(w, http.StatusOK, u)
}
func (uh *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := uh.service.LoginUser(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponse(w, http.StatusOK, u)
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request dto.UserRequest
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	user, err := uh.service.CreateUser(ctx, request)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusOK, user)
}

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}
