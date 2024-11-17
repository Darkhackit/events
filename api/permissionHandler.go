package api

import (
	"encoding/json"
	"fmt"
	"github.com/Darkhackit/events/dto"
	"github.com/Darkhackit/events/service"
	validator2 "github.com/Darkhackit/events/validator"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type PermissionHandler struct {
	service service.PermissionService
}

func (ph *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var request dto.PermissionRequest
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
		WriteResponse(w, http.StatusBadRequest, errors)
		return
	}
	result, err := ph.service.CreatePermission(ctx, request)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusCreated, result)
}

func (ph *PermissionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	result, err := ph.service.GetPermissions(r.Context())
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteResponse(w, http.StatusOK, result)
}
func (ph *PermissionHandler) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	var request dto.AssignPermissionRequest
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
		WriteResponse(w, http.StatusBadRequest, errors)
		return
	}
	err = ph.service.AssignPermission(ctx, request)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = fmt.Fprintf(w, "Permission assigned")
	if err != nil {
		return
	}
}
