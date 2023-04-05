// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/docker-mariadb-temporal-tables-tutorial/db/models"
)

var (
	originalCurrentDateProvider = currentDate
	originalRandomIntProvider   = randomInt
)

func TestConnectToPostgres(t *testing.T) {
	testCases := []struct {
		name          string
		mockSqlOpen   func(driverName string, dataSourceName string) (*sql.DB, error)
		expectedError error
	}{
		{
			name: "happy path",
			mockSqlOpen: func(driverName, dataSourceName string) (*sql.DB, error) {
				db, _, _ := sqlmock.New()
				return db, nil
			},
		},
		{
			name: "error",
			mockSqlOpen: func(driverName, dataSourceName string) (*sql.DB, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("random error"),
		},
		{
			name: "error pinging",
			mockSqlOpen: func(driverName, dataSourceName string) (*sql.DB, error) {
				db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
				mock.ExpectPing().WillReturnError(errors.New("random error"))
				return db, nil
			},
			expectedError: errors.New("random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlOpen = tc.mockSqlOpen
			db, err := ConnectToMariaDb("user", "pass", "host", "port", "schema")
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
				require.NotNil(t, db)
			}
		})
	}
}

func TestAdvanceMariaDbTimestamp(t *testing.T) {
	testCases := []struct {
		name                string
		currentDateProvider func() time.Time
		randomIntProvider   func() int
		mockClosure         func() *sql.DB
		expectedOutput      *models.DbTimestamp
		expectedError       error
	}{
		{
			name: "happy path",
			currentDateProvider: func() time.Time {
				t, _ := time.Parse(time.DateOnly, "2023-01-01")
				return t
			},
			randomIntProvider: func() int {
				return 1
			},
			mockClosure: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				q := `SET @@timestamp = UNIX_TIMESTAMP('2024-01-01')`
				mock.ExpectExec(regexp.QuoteMeta(q)).WillReturnResult(driver.ResultNoRows)
				return db
			},
			expectedOutput: &models.DbTimestamp{
				Timestamp: "2024-01-01",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handleCurrentDateProvider(tc.currentDateProvider)
			handleRandomIntProvider(tc.randomIntProvider)
			db := tc.mockClosure()
			output, err := AdvanceMariaDbTimestamp(context.TODO(), db)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestSetDefaultMariaDbTimestamp(t *testing.T) {
	testCases := []struct {
		name                string
		currentDateProvider func() time.Time
		mockClosure         func() *sql.DB
		expectedOutput      *models.DbTimestamp
		expectedError       error
	}{
		{
			name: "happy path",
			currentDateProvider: func() time.Time {
				t, _ := time.Parse(time.DateOnly, "2023-01-01")
				return t
			},
			mockClosure: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				q := `SET @@timestamp = default`
				mock.ExpectExec(regexp.QuoteMeta(q)).WillReturnResult(driver.ResultNoRows)
				return db
			},
			expectedOutput: &models.DbTimestamp{
				Timestamp: "2023-01-01",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handleCurrentDateProvider(tc.currentDateProvider)
			db := tc.mockClosure()
			output, err := SetDefaultMariaDbTimestamp(context.TODO(), db)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func handleCurrentDateProvider(mocked func() time.Time) {
	if mocked != nil {
		currentDate = mocked
	} else {
		currentDate = originalCurrentDateProvider
	}
}

func handleRandomIntProvider(mocked func() int) {
	if mocked != nil {
		randomInt = mocked
	} else {
		randomInt = originalRandomIntProvider
	}
}
