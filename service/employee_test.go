package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mekari-test/domain"
	mockdomain "mekari-test/mock"
	"testing"
	"time"
)

func TestEmployeeService_Register(t *testing.T) {
	var (
		employee = domain.Employee{
			Id:        1,
			FirstName: "Muhammad",
			LastName:  "Fadli",
			Email:     "mnurfadd@gmail.com",
			HireDate:  time.Time{},
		}
	)

	testCases := []struct {
		name         string
		expectedErr  bool
		context      context.Context
		mockUserRepo func(mock *mockdomain.MockEmployeeRepository)
	}{
		{
			name:        "Success",
			expectedErr: false,
			context:     context.Background(),
			mockUserRepo: func(mock *mockdomain.MockEmployeeRepository) {
				mock.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:        "Failed",
			expectedErr: true,
			context:     context.Background(),
			mockUserRepo: func(mock *mockdomain.MockEmployeeRepository) {
				mock.EXPECT().Store(gomock.Any(), gomock.Any()).
					Return(errors.New("internal server error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err          error
				mockCtrl     = gomock.NewController(t)
				mockUserRepo = mockdomain.NewMockEmployeeRepository(mockCtrl)
				service      = employeeService{employeeRepo: mockUserRepo}
			)
			defer mockCtrl.Finish()
			tc.mockUserRepo(mockUserRepo)

			err = service.Register(tc.context, employee)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}

func TestEmployeeService_GetById(t *testing.T) {
	var (
		employee = domain.Employee{
			Id:        1,
			FirstName: "Muhammad",
			LastName:  "Fadli",
			Email:     "mnurfadd@gmail.com",
			HireDate:  time.Time{},
		}
	)

	testCases := []struct {
		name         string
		expectedErr  bool
		context      context.Context
		mockUserRepo func(mock *mockdomain.MockEmployeeRepository)
		expected     domain.Employee
	}{
		{
			name:        "Success",
			expectedErr: false,
			context:     context.Background(),
			mockUserRepo: func(mock *mockdomain.MockEmployeeRepository) {
				mock.EXPECT().GetById(gomock.Any(), gomock.Any()).
					Return(employee, nil)
			},
			expected: employee,
		},
		{
			name:        "Failed",
			expectedErr: true,
			context:     context.Background(),
			mockUserRepo: func(mock *mockdomain.MockEmployeeRepository) {
				mock.EXPECT().GetById(gomock.Any(), gomock.Any()).
					Return(domain.Employee{}, errors.New("internal server error"))
			},
			expected: domain.Employee{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err          error
				res          domain.Employee
				mockCtrl     = gomock.NewController(t)
				mockUserRepo = mockdomain.NewMockEmployeeRepository(mockCtrl)
				service      = employeeService{employeeRepo: mockUserRepo}
			)
			defer mockCtrl.Finish()
			tc.mockUserRepo(mockUserRepo)

			res, err = service.GetById(tc.context, employee.Id)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, res)
		})
	}
}
