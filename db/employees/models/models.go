// Copyright (c) 2022 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package models

import "time"

// Employee represents an employee record in the database.
type Employee struct {
	Id         int64   `json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
}

// EmployeeHistory represents an employee history record in the database.
type EmployeeHistory struct {
	Employee
	RowStart time.Time `json:"row_start"`
	RowEnd   time.Time `json:"row_end"`
}

// NewEmployee is used to create a new employee record.
type NewEmployee struct {
	FirstName  string  `json:"first_name" validate:"required"`
	LastName   string  `json:"last_name" validate:"required"`
	Department string  `json:"department" validate:"required"`
	Salary     float64 `json:"salary" validate:"required,gt=0"`
}

// UpdateEmployee is used to update an employee record.
type UpdateEmployee struct {
	FirstName  *string  `json:"first_name"`
	LastName   *string  `json:"last_name"`
	Department *string  `json:"department"`
	Salary     *float64 `json:"salary"`
}

// FirstNameIsFulfilled checks whether FirstName is fulfilled.
func (u *UpdateEmployee) FirstNameIsFulfilled() bool {
	return u.FirstName != nil && *u.FirstName != ""
}

// LastNameIsFulfilled checks whether LastName is fulfilled.
func (u *UpdateEmployee) LastNameIsFulfilled() bool {
	return u.LastName != nil && *u.LastName != ""
}

// DepartmentIsFulfilled checks whether Department is fulfilled.
func (u *UpdateEmployee) DepartmentIsFulfilled() bool {
	return u.Department != nil && *u.Department != ""
}

// SalaryIsFulfilled checks whether Salary is fulfilled.
func (u *UpdateEmployee) SalaryIsFulfilled() bool {
	return u.Salary != nil && *u.Salary > 0
}
