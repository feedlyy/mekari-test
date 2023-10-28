package service

import (
	"context"
	"mekari-test/domain"
	"time"
)

type employeeService struct {
	employeeRepo domain.EmployeeRepository
}

func NewEmployeeService(e domain.EmployeeRepository) domain.EmployeeService {
	return employeeService{employeeRepo: e}
}

func (e employeeService) GetAllEmployee(ctx context.Context) ([]domain.Employee, error) {
	return e.employeeRepo.Get(ctx)
}

func (e employeeService) GetById(ctx context.Context, id int) (domain.Employee, error) {
	return e.employeeRepo.GetById(ctx, id)
}

func (e employeeService) Register(ctx context.Context, employee domain.Employee) error {
	return e.employeeRepo.Store(ctx, employee)
}

func (e employeeService) Update(ctx context.Context, employee domain.Employee) error {
	var (
		emp = domain.Employee{}
		err error
	)
	emp, err = e.employeeRepo.GetById(ctx, employee.Id)
	if err != nil {
		return err
	}

	switch {
	case employee.FirstName != "":
		emp.FirstName = employee.FirstName
	case employee.LastName != "":
		emp.LastName = employee.LastName
	case employee.Email != "":
		emp.Email = employee.Email
	case employee.HireDate != (time.Time{}):
		emp.HireDate = employee.HireDate
	}

	return e.employeeRepo.Update(ctx, emp)
}

func (e employeeService) Delete(ctx context.Context, id int) error {
	return e.employeeRepo.Delete(ctx, id)
}
