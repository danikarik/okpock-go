package sequel_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/danikarik/okpock/pkg/store/sequel"
	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertPass(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	testCase := struct {
		SerialNumber string
		AuthToken    string
		PassTypeID   string
	}{
		SerialNumber: "1967bce8-fb9c-4be7-8946-c1a3a7607a88",
		AuthToken:    uuid.NewV4().String(),
		PassTypeID:   "com.example.pass",
	}

	schema := []string{tempPassesTable}
	data := []string{}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	err = db.InsertPass(ctx, testCase.SerialNumber, testCase.AuthToken, testCase.PassTypeID)
	assert.NoError(err)
}

func TestUpdatePass(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	defaultTime, _ := time.Parse(sequel.TimeFormat, "2019-03-26 00:00:00")
	testCase := struct {
		SerialNumber string
		AuthToken    string
		PassTypeID   string
		UpdatedAt    time.Time
	}{
		SerialNumber: "2c91cb65-29ad-465a-bbbc-968f0ca224e9",
		AuthToken:    "secret",
		PassTypeID:   "com.example.pass",
		UpdatedAt:    defaultTime,
	}

	schema := []string{tempPassesTable}
	data := []string{
		fmt.Sprintf(
			insertPassesTable,
			testCase.SerialNumber,
			testCase.AuthToken,
			testCase.PassTypeID,
			testCase.UpdatedAt.Format(sequel.TimeFormat),
		),
	}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	err = db.UpdatePass(ctx, testCase.SerialNumber)
	assert.NoError(err)
}

func TestFindPass(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	testCase := struct {
		SerialNumber string
		AuthToken    string
		PassTypeID   string
	}{
		SerialNumber: "2c91cb65-29ad-465a-bbbc-968f0ca224e9",
		AuthToken:    "secret",
		PassTypeID:   "com.example.pass",
	}

	schema := []string{tempPassesTable}
	data := []string{
		fmt.Sprintf(
			insertPassesTable,
			testCase.SerialNumber,
			testCase.AuthToken,
			testCase.PassTypeID,
			time.Now().Format(sequel.TimeFormat),
		),
	}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	ok, err := db.FindPass(ctx, testCase.SerialNumber, testCase.AuthToken, testCase.PassTypeID)
	assert.NoError(err)
	assert.True(ok)
}

func TestFindRegistration(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	testCases := []struct {
		SerialNumber string
		PassTypeID   string
		DeviceID     string
		Expected     bool
	}{
		{
			SerialNumber: "2c91cb65-29ad-465a-bbbc-968f0ca224e9",
			PassTypeID:   "com.example.pass",
			DeviceID:     "52a26307-aef2-45de-af72-4e5acfa55b8d",
			Expected:     true,
		},
		{
			SerialNumber: "1967bce8-fb9c-4be7-8946-c1a3a7607a88",
			PassTypeID:   "com.example.pass",
			DeviceID:     "52a26307-aef2-45de-af72-4e5acfa55b8d",
			Expected:     false,
		},
	}

	schema := []string{tempRegistrationsTable}
	data := []string{
		fmt.Sprintf(
			insertRegistrationsTable,
			uuid.NewV4().String(),
			testCases[0].DeviceID,
			uuid.NewV4().String(),
			testCases[0].SerialNumber,
			testCases[0].PassTypeID,
		),
	}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	for _, c := range testCases {
		ok, err := db.FindRegistration(ctx, c.DeviceID, c.SerialNumber)
		assert.NoError(err)
		assert.Equal(c.Expected, ok)
	}
}

func TestFindSerialNumbers(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	testCase := struct {
		SerialNumber string
		AuthToken    string
		PassTypeID   string
		DeviceID     string
		Expected     []string
	}{
		SerialNumber: "1967bce8-fb9c-4be7-8946-c1a3a7607a88",
		AuthToken:    "secret",
		PassTypeID:   "com.example.pass",
		DeviceID:     "52a26307-aef2-45de-af72-4e5acfa55b8d",
		Expected:     []string{"1967bce8-fb9c-4be7-8946-c1a3a7607a88"},
	}

	schema := []string{
		tempPassesTable,
		tempRegistrationsTable,
	}
	data := []string{
		fmt.Sprintf(
			insertPassesTable,
			testCase.SerialNumber,
			testCase.AuthToken,
			testCase.PassTypeID,
			time.Now().Format(sequel.TimeFormat),
		),
		fmt.Sprintf(
			insertRegistrationsTable,
			uuid.NewV4().String(),
			testCase.DeviceID,
			uuid.NewV4().String(),
			testCase.SerialNumber,
			testCase.PassTypeID,
		),
	}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	sns, err := db.FindSerialNumbers(ctx, testCase.DeviceID, testCase.PassTypeID, "")
	assert.NoError(err)
	assert.Equal(testCase.Expected, sns)
}

func TestLatestPass(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	defaultTime, _ := time.Parse(sequel.TimeFormat, "2019-03-22 00:00:00")
	testCase := struct {
		SerialNumber string
		AuthToken    string
		PassTypeID   string
		UpdatedAt    time.Time
	}{
		SerialNumber: "1967bce8-fb9c-4be7-8946-c1a3a7607a88",
		AuthToken:    "secret",
		PassTypeID:   "com.example.pass",
		UpdatedAt:    defaultTime,
	}

	schema := []string{tempPassesTable}
	data := []string{
		fmt.Sprintf(
			insertPassesTable,
			testCase.SerialNumber,
			testCase.AuthToken,
			testCase.PassTypeID,
			testCase.UpdatedAt.Format(sequel.TimeFormat),
		),
	}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	lastUpdate, err := db.LatestPass(ctx, testCase.SerialNumber, testCase.AuthToken, testCase.PassTypeID)
	assert.NoError(err)
	assert.False(lastUpdate.Sub(testCase.UpdatedAt) > 0)
}

func TestInsertRegistration(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	testCase := struct {
		DeviceID     string
		PushToken    string
		SerialNumber string
		PassTypeID   string
	}{
		DeviceID:     uuid.NewV4().String(),
		PushToken:    uuid.NewV4().String(),
		SerialNumber: "1967bce8-fb9c-4be7-8946-c1a3a7607a88",
		PassTypeID:   "com.example.pass",
	}

	schema := []string{tempRegistrationsTable}
	data := []string{}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	err = db.InsertRegistration(ctx, testCase.DeviceID, testCase.PushToken, testCase.SerialNumber, testCase.PassTypeID)
	assert.NoError(err)
}

func TestDeleteRegistration(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	testCase := struct {
		DeviceID     string
		PushToken    string
		SerialNumber string
		PassTypeID   string
	}{
		DeviceID:     uuid.NewV4().String(),
		PushToken:    uuid.NewV4().String(),
		SerialNumber: "1967bce8-fb9c-4be7-8946-c1a3a7607a88",
		PassTypeID:   "com.example.pass",
	}

	schema := []string{tempRegistrationsTable}
	data := []string{
		fmt.Sprintf(
			insertRegistrationsTable,
			uuid.NewV4().String(),
			testCase.DeviceID,
			testCase.PushToken,
			testCase.SerialNumber,
			testCase.PassTypeID,
		),
	}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	ok, err := db.DeleteRegistration(ctx, testCase.DeviceID, testCase.SerialNumber, testCase.PassTypeID)
	assert.NoError(err)
	assert.True(ok)
}

func TestInsertLog(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	testCase := struct {
		RemoteAddr string
		RequestID  string
		Message    string
	}{
		RemoteAddr: uuid.NewV4().String(),
		RequestID:  uuid.NewV4().String(),
		Message:    "test",
	}

	schema := []string{tempLogsTable}
	data := []string{}

	conn, err := executeTempScripts(ctx, t, schema, data)
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	db := sequel.New(conn)

	err = db.InsertLog(ctx, testCase.RemoteAddr, testCase.RequestID, testCase.Message)
	assert.NoError(err)
}
