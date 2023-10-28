package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

func (e *EmployeeHandler) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		err      error
		employee = domain.Employee{
			FirstName: r.PostFormValue("first_name"),
			LastName:  r.PostFormValue("last_name"),
			Email:     r.PostFormValue("email"),
		}
		hireDate = r.PostFormValue("hire_date")
		resp     = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		date time.Time
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// validate request
	err = employee.Validate(hireDate)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	date, err = time.Parse("2006-01-02", hireDate)
	if err != nil {
		resp.Status = helpers.FailMsg
		errMsg := fmt.Sprintf("invalid parsing date, err:%v", err.Error())
		resp.Data = errMsg

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	employee.HireDate = date

	err = e.employeeService.Register(ctx, employee)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(resp)
	return
}

func (e *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		err  error
		id   int
		resp = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
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

	err = e.employeeService.Delete(ctx, id)
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

	json.NewEncoder(w).Encode(resp)
	return
}

func (e *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		err  error
		id   int
		resp = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		employee = domain.Employee{
			FirstName: r.PostFormValue("first_name"),
			LastName:  r.PostFormValue("last_name"),
			Email:     r.PostFormValue("email"),
		}
		hireDate = r.PostFormValue("hire_date")
		date     time.Time
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// parse id into int
	id, err = strconv.Atoi(ps.ByName("id"))
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	employee.Id = id

	if hireDate != "" {
		date, err = time.Parse("2006-01-02", hireDate)
		if err != nil {
			resp.Status = helpers.FailMsg
			errMsg := fmt.Sprintf("invalid parsing date, err:%v", err.Error())
			resp.Data = errMsg

			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		employee.HireDate = date
	}

	if employee.Email != "" {
		valid := domain.IsValidEmail(employee.Email)

		if !valid {
			resp.Status = helpers.FailMsg
			resp.Data = "please input valid email"

			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	err = e.employeeService.Update(ctx, employee)
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

	json.NewEncoder(w).Encode(resp)
	return
}
