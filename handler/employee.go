package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"mekari-test/domain"
	"mekari-test/helpers"
	"net/http"
	"strconv"
	"time"
)

type EmployeeHandler struct {
	employeeService domain.EmployeeService
	timeout         time.Duration
}

func NewEmployeeHandler(emp domain.EmployeeService, timeout time.Duration) EmployeeHandler {
	handler := &EmployeeHandler{
		employeeService: emp,
		timeout:         timeout,
	}

	return *handler
}

func (e *EmployeeHandler) GetAllEmployee(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err       error
		employees []domain.Employee
		resp      = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	employees, err = e.employeeService.GetAllEmployee(ctx)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = employees
	json.NewEncoder(w).Encode(resp)
	return
}

func (e *EmployeeHandler) GetEmployeeById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		err      error
		employee = domain.Employee{}
		resp     = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		id int
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	id, err = strconv.Atoi(ps.ByName("id"))
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	employee, err = e.employeeService.GetById(ctx, id)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		switch {
		case errors.Is(err, sql.ErrNoRows):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	resp.Data = employee
	json.NewEncoder(w).Encode(resp)
	return
}
