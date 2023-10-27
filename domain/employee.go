package domain

import (
	"context"
	"time"
)

type Employee struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Email     string    `json:"email,omitempty"`
	HireDate  time.Time `json:"hire_date,omitempty"`
}

type EmployeeRepository interface {
	Get(ctx context.Context) ([]Employee, error)
	GetById(ctx context.Context, id int) (Employee, error)
	Store(ctx context.Context, employee Employee) error
	Update(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
}

type EmployeeService interface {
	GetAllEmployee(ctx context.Context) ([]Employee, error)
	GetById(ctx context.Context, id int) (Employee, error)
	Register(ctx context.Context, id int) error
	Update(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
}
