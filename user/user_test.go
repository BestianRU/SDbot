package user

import "testing"

type testDB struct {
	//DBer
	rows testRows
}
type testRows struct {
	//	rows []testRow
}

func (r testRows) Next() bool {
	return true
}

func (r testRows) Scan(dest ...interface{}) error {
	m := "+1-2(3)456 7890 "
	*dest[1].(*string) = m
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

func (d testDB) Query(query string, args ...interface{}) (rowsScanner, error) {
	return d.rows, nil
}

func (d testDB) QueryRow(query string, args ...interface{}) scanner {
	return d.QueryRow(query, args...)
}

func TestGetUserFromSQLByPhone(t *testing.T) {
	var db testDB

	u, err := getUserFromSQLByPhone("1234567890", db)
	if err != nil {
		t.Error("Error in getUserFromSQLByPhone", err)
	}
	if u.Phone != "1234567890" {
		t.Error("Expected return phone 1234567890, got:", u.Phone)
	}
	if u.Email != "email@gmail.com" {
		t.Error("Expected return Email email@gmail.com, got:", u.Email)
	}
}
