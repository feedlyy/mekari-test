package repository

import (
	"context"
	sql2 "database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mekari-test/domain"
)

type employeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(db *sqlx.DB) domain.EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

func (e *employeeRepository) Get(ctx context.Context) ([]domain.Employee, error) {
	var (
		err  error
		res  []domain.Employee
		sql  string
		stmt *sqlx.Stmt
		rows *sqlx.Rows
	)
	sql, _, err = sq.Select("id", "first_name", "last_name",
		"email", "hire_date").From("employees").ToSql()
	if err != nil {
		logrus.Errorf("Employee - Repository|err when generate sql, err:%v", err)
		return nil, err
	}

	stmt, err = e.db.PreparexContext(ctx, sql)
	if err != nil {
		logrus.Errorf("Employee - Repository|err when init prepare statement, err:%v", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err = stmt.QueryxContext(ctx)
	for rows.Next() {
		var employee = domain.Employee{}

		err = rows.Scan(&employee.Id, &employee.FirstName, &employee.LastName,
			&employee.Email, &employee.HireDate)
		if err != nil {
			logrus.Errorf("Employee - Repository|err when scan struct, err:%v", err)
			return nil, err
		}

		res = append(res, employee)
	}

	return res, nil
}

func (e *employeeRepository) GetById(ctx context.Context, id int) (domain.Employee, error) {
	var (
		err  error
		res  = domain.Employee{}
		sql  string
		stmt *sqlx.Stmt
		row  *sqlx.Row
	)
	sql, _, err = sq.Select("id", "first_name", "last_name",
		"email", "hire_date").From("employees").Where(sq.And{
		sq.Eq{"id": "id"},
	}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Employee - Repository|err when generate sql, err:%v", err)
		return domain.Employee{}, err
	}

	stmt, err = e.db.PreparexContext(ctx, sql)
	if err != nil {
		logrus.Errorf("Employee - Repository|err when init prepare statement, err:%v", err)
		return domain.Employee{}, err
	}
	defer stmt.Close()

	row = stmt.QueryRowxContext(ctx, id)
	err = row.Scan(&res.Id, &res.FirstName, &res.LastName,
		&res.Email, &res.HireDate)
	if err != nil && err != sql2.ErrNoRows {
		logrus.Errorf("Employee - Repository|err when scan, err:%v", err)
		return domain.Employee{}, err
	}

	if err == sql2.ErrNoRows {
		logrus.Errorf("Employee - Repository|data not found, err:%v", err)
		return domain.Employee{}, err
	}

	return res, nil
}

func (e *employeeRepository) Store(ctx context.Context, employee domain.Employee) error {
	var (
		err error
		sql string
	)
	sql, _, err = sq.Insert("employees").Columns("first_name", "last_name",
		"email", "hire_date").
		Values("first_name", "last_name", "email", "hire_date").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Employees - Repository|err when generate sql, err:%v", err)
		return err
	}

	_, err = e.db.ExecContext(ctx, sql, employee.FirstName, employee.LastName, employee.Email,
		employee.HireDate)
	if err != nil {
		logrus.Errorf("Employees - Repository|err when store data, err:%v", err)
		return err
	}

	return nil
}

func (e *employeeRepository) Update(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (e *employeeRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
