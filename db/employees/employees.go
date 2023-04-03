// Copyright (c) 2022 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package employees

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	mariaDb "github.com/tiagomelo/docker-mariadb-temporal-tables-tutorial/db"
	"github.com/tiagomelo/docker-mariadb-temporal-tables-tutorial/db/employees/models"
)

// For ease of unit testing.
var (
	readEmployee = func(row *sql.Row, dest ...any) error {
		return row.Scan(dest...)
	}
	readEmployees = func(rows *sql.Rows, dest ...any) error {
		return rows.Scan(dest...)
	}
)

// GetById returns a current employee with given id.
func GetById(ctx context.Context, db *sql.DB, id uint) (*models.Employee, error) {
	q := `
	SELECT id, first_name, last_name, salary, department
	FROM employees
	WHERE id = ?
	`
	var employee models.Employee
	row := db.QueryRowContext(ctx, q, id)
	if err := readEmployee(row,
		&employee.Id,
		&employee.FirstName,
		&employee.LastName,
		&employee.Salary,
		&employee.Department,
	); err != nil {
		return nil, mariaDb.ErrDBNotFound
	}
	return &employee, nil
}

// GetAll returns all current employees.
func GetAll(ctx context.Context, db *sql.DB) ([]models.Employee, error) {
	q := `
	SELECT id, first_name, last_name, salary, department
	FROM employees
	`
	employees := make([]models.Employee, 0)
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return employees, err
	}
	defer rows.Close()
	for rows.Next() {
		var employee models.Employee
		if err = readEmployees(rows,
			&employee.Id,
			&employee.FirstName,
			&employee.LastName,
			&employee.Salary,
			&employee.Department,
		); err != nil {
			return employees, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

// Create creates an employee.
func Create(ctx context.Context, db *sql.DB, newEmployee *models.NewEmployee) (*models.Employee, error) {
	q := `
	INSERT INTO
		employees(first_name, last_name, salary, department)
	VALUES
		(?, ?, ?, ?)
	RETURNING
		id, first_name, last_name, salary, department
	`
	var employee models.Employee
	row := db.QueryRowContext(ctx, q, newEmployee.FirstName, newEmployee.LastName, newEmployee.Salary, newEmployee.Department)
	if err := readEmployee(row,
		&employee.Id,
		&employee.FirstName,
		&employee.LastName,
		&employee.Salary,
		&employee.Department,
	); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == mariaDb.UniqueViolation {
			return nil, mariaDb.ErrDBDuplicatedEntry
		}
		return nil, err
	}
	return &employee, nil
}

// handleEmployeeChanges updates the changed properties.
func handleEmployeeChanges(updateEmployee *models.UpdateEmployee, dbEmployee *models.Employee) {
	if updateEmployee.FirstNameIsFulfilled() {
		dbEmployee.FirstName = *updateEmployee.FirstName
	}
	if updateEmployee.LastNameIsFulfilled() {
		dbEmployee.LastName = *updateEmployee.LastName
	}
	if updateEmployee.SalaryIsFulfilled() {
		dbEmployee.Salary = *updateEmployee.Salary
	}
	if updateEmployee.DepartmentIsFulfilled() {
		dbEmployee.Department = *updateEmployee.Department
	}
}

// Update updates an employee.
func Update(ctx context.Context, db *sql.DB, employeeId uint, updateEmployee *models.UpdateEmployee) (*models.Employee, error) {
	q := `
	UPDATE employees
	SET 
		first_name = ?,
		last_name = ?,
		salary = ?,
		department = ?
	WHERE
		id = ?
	`
	dbEmployee, err := GetById(ctx, db, employeeId)
	if err != nil {
		return nil, err
	}
	handleEmployeeChanges(updateEmployee, dbEmployee)
	_, err = db.ExecContext(ctx, q, dbEmployee.FirstName, dbEmployee.LastName, dbEmployee.Salary, dbEmployee.Department, dbEmployee.Id)
	return dbEmployee, err
}

// Delete deletes an employee.
func Delete(ctx context.Context, db *sql.DB, id uint) error {
	q := `
	DELETE FROM
		employees
	WHERE id = ?
	`
	_, err := db.ExecContext(ctx, q, id)
	return err
}
