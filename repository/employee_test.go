package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"mekari-test/domain"
	"testing"
	"time"
)

func TestGetEmployeeById(t *testing.T) {
	var (
		returns = domain.Employee{
			Id:        1,
			FirstName: "Muhammad",
			LastName:  "Fadli",
			Email:     "mnurfadd@gmail.com",
			HireDate:  time.Time{},
		}
	)
	testCases := []struct {
		name        string
		expectedErr bool
		context     context.Context
		doMockDB    func(mock sqlmock.Sqlmock)
		expected    domain.Employee
		input       int
	}{
		{
			name:        "Success",
			expectedErr: false,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("^SELECT(.+)").ExpectQuery().WithArgs(returns.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "hire_date"}).
						AddRow(1, "Muhammad", "Fadli", "mnurfadd@gmail.com", time.Time{}))
			},
			input:    returns.Id,
			expected: returns,
		},
		{
			name:        "Failed",
			expectedErr: true,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("^SELECT(.+)").ExpectQuery().WithArgs(returns.Id).
					WillReturnError(errors.New("internal server error"))
			},
			expected: domain.Employee{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err    error
				res    domain.Employee
				db     *sql.DB
				mock   sqlmock.Sqlmock
				mockDB *sqlx.DB
			)
			db, mock, err = sqlmock.New()
			if err != nil {
				panic(err)
			}
			defer db.Close()
			tc.doMockDB(mock)

			mockDB = sqlx.NewDb(db, "postgres")

			repoEmployee := NewEmployeeRepository(mockDB)
			res, err = repoEmployee.GetById(tc.context, tc.input)
			assert.Equal(t, tc.expectedErr, err != nil)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestInsertEmployee(t *testing.T) {
	var (
		returns = domain.Employee{
			Id:        1,
			FirstName: "Muhammad",
			LastName:  "Fadli",
			Email:     "mnurfadd@gmail.com",
			HireDate:  time.Time{},
		}
	)

	testCases := []struct {
		name        string
		expectedErr bool
		context     context.Context
		doMockDB    func(mock sqlmock.Sqlmock)
		input       domain.Employee
	}{
		{
			name:        "Success",
			expectedErr: false,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^INSERT INTO(.+)").
					WithArgs(returns.FirstName, returns.LastName, returns.Email, returns.HireDate).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: returns,
		},
		{
			name:        "Failed",
			expectedErr: true,
			context:     context.Background(),
			doMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("^INSERT INTO(.+)").
					WillReturnError(errors.New("internal server error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				err    error
				db     *sql.DB
				mock   sqlmock.Sqlmock
				mockDB *sqlx.DB
			)
			db, mock, err = sqlmock.New()
			if err != nil {
				panic(err)
			}
			defer db.Close()
			tc.doMockDB(mock)

			mockDB = sqlx.NewDb(db, "postgres")

			repoEmployee := NewEmployeeRepository(mockDB)
			err = repoEmployee.Store(tc.context, tc.input)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}
