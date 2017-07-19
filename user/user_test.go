package user

import (
	"errors"
	"testing"
)

type testData struct {
	id    uint64
	phone string
	mail  string
	name  string
}

var d = []testData{
	{phone: "+1-2(3)456 7890 ",
		mail: "abc_123@cde.com",
		id:   12345,
		name: "Ivan",
	},
	{phone: " +1-2(3)456 789990 ",
		mail: "abc@cde_567.com",
		id:   12346,
		name: "test name"},
}

type testDB struct {
	//DBer
	rows *testRows
}
type testRows struct {
	thisRow int
}

func (r *testRows) Next() bool {
	if len(d) > r.thisRow-1 {
		return true
	}
	return false
}

func (r *testRows) Scan(dest ...interface{}) error {
	if len(dest) != 3 {
		return errors.New("Scan must have 3 arguments")
	}
	*dest[0].(*uint64) = d[r.thisRow].id
	*dest[1].(*string) = d[r.thisRow].phone
	*dest[2].(*string) = d[r.thisRow].name
	r.thisRow++
	return nil
}

type testRow struct {
	id       int
	firsName string
	phone    string
	mail     string
}

func (r testRow) Scan(dest ...interface{}) error {

	dest[0] = "email@gmail.com"
	return nil
}

func (d testDB) Close() error {
	return nil
}

func (d testDB) Query(query string, args ...interface{}) (rowser, error) {
	return d.rows, nil
}

func (d testDB) QueryRow(query string, args ...interface{}) scanner {
	return d.QueryRow(query, args...)
}

// func TestGetUserFromSQLByPhone(t *testing.T) {
// 	var db testDB

// 	u, err := getUserFromSQLByPhone("1234567890", db)
// 	if err != nil {
// 		t.Error("Error in getUserFromSQLByPhone", err)
// 	}
// 	if u.Phone != "1234567890" {
// 		t.Error("Expected return phone 1234567890, got:", u.Phone)
// 	}
// 	if u.Email != "email@gmail.com" {
// 		t.Error("Expected return Email email@gmail.com, got:", u.Email)
// 	}
// }

func TestGetUserFullName(t *testing.T) {
	var u User
	var db testDB
	db.rows = new(testRows)

	testPhone := "123456789990"
	err := getUserFullName(testPhone, &u, db)
	if err != nil {
		t.Error("Error in getUserFullName", err)
	}
	if u.Phone != testPhone {
		t.Error("Error returning phone from getUserFullName", u.Phone)
	}
	if u.FullName != "test name" {
		t.Error("Error returning FullName from getUserFullName", u.FullName)
	}
	if u.SDId != 12346 {
		t.Error("Error returning SDId from getUserFullName", u.SDId)
	}
}
