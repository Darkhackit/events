package api

import (
	"encoding/json"
	"github.com/Darkhackit/events/dto"
	"github.com/Darkhackit/events/service"
	validator2 "github.com/Darkhackit/events/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type RoleHandler struct {
	service service.RoleService
}

func (rh *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var request dto.RoleRequest
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errors := validator2.TransformValidationErrors(err)
		WriteResponse(w, http.StatusUnprocessableEntity, errors)
		return
	}
	role, err := rh.service.CreateRole(ctx, request)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusCreated, role)
}
func (rh *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID, err := strconv.Atoi(vars["role_id"])
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	err = rh.service.DeleteRole(ctx, roleID)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusNoContent, "successfully deleted role")

}
func (rh *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID := vars["role_id"]
	ctx := r.Context()
	if roleID == "" {
		WriteResponse(w, http.StatusBadRequest, "role_id is required")
		return
	}
	Int, err := strconv.Atoi(roleID)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
	}
	result, err := rh.service.GetRole(ctx, Int)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusOK, result)
}
func (rh *RoleHandler) AssignUserRole(w http.ResponseWriter, r *http.Request) {
	var request dto.UserRoleRequest
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errors := validator2.TransformValidationErrors(err)

		WriteResponse(w, http.StatusUnprocessableEntity, errors)
		return
	}

	err = rh.service.AssignRoleToUser(ctx, request)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusNoContent, "successfully assigned role to user")
}
func (rh *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := rh.service.GetRoles(ctx)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusOK, result)
}
