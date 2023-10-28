package service

import (
	"context"
	"mekari-test/domain"
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

func (e employeeService) Update(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (e employeeService) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
