package user

import "testing"
import "errors"
import "reflect"

type testFile struct {
	data []byte
}

func (f *testFile) Read(ret []byte) (int, error) {
	var i int

	if i = copy(ret, f.data); i == 0 {
		return 0, errors.New("Error reading file")
	}
	return i, nil
}

func (f *testFile) Write(ret []byte) (int, error) {
	var i int
	f.data = make([]byte, len(ret))
	if i = copy(f.data, ret); i == 0 {
		return 0, errors.New("Error writing file")
	}
	return i, nil
}

type testData struct {
	id    uint64
	phone string
	mail  string
	name  string
}
type testDB struct {
	rows *testRows
}
type testRows struct {
	thisRow int
	date    []testData
}

func (r *testRows) Next() bool {
	if len(r.date) > r.thisRow {
		return true
	}
	return false
}

func (r *testRows) Scan(dest ...interface{}) error {
	//test for Query("SELECT id,mobile,comment FROM glpi_users WHERE mobile IS NOT NULL AND comment IS NOT NULL")
	if len(dest) == 3 {
		*dest[0].(*uint64) = r.date[r.thisRow].id
		*dest[1].(*string) = r.date[r.thisRow].phone
		*dest[2].(*string) = r.date[r.thisRow].name
	}
	//test for Query("SELECT email FROM glpi_useremails WHERE users_id=?", u.SDId)
	if len(dest) == 1 {
		*dest[0].(*string) = r.date[r.thisRow].mail
	}
	r.thisRow++
	return nil
}

type testRow struct {
	id       int
	firsName string
	phone    string
	mail     string
}

func (d testDB) Close() error {
	return nil
}

func (d testDB) Query(query string, args ...interface{}) (rowser, error) {
	d.rows = new(testRows)
	switch query {
	case "SELECT email FROM glpi_useremails WHERE users_id=?":
		if args[0].(uint64) == 12346 {
			d.rows.date = []testData{
				{
					phone: " +1-2(3)456 789990 ",
					mail:  "abc@cde_567.com",
					id:    12346,
					name:  "test name"},
			}
		}

		return d.rows, nil
	case "SELECT id,mobile,comment FROM glpi_users WHERE mobile IS NOT NULL AND comment IS NOT NULL":
		d.rows.date = []testData{
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
		return d.rows, nil

	}
	return d.rows, nil
}

func TestAuthUser(t *testing.T) {
	tf := new(testFile)
	var au AuthUser
	au.MapUser = make(map[string]User)
	au.MapUser["abc_123@cde.com"] = User{
		TId:      123,
		SDId:     345,
		Phone:    "+1-2(3)456 7890 ",
		Email:    "abc_123@cde.com",
		FullName: "Ivan",
	}
	if au.save(tf) != nil {
		t.Error("Error test writing for AuthUser")
	}
	var newAU AuthUser
	newAU.MapUser = make(map[string]User)
	if newAU.read(tf) != nil {
		t.Error("Error test reading for AuthUser")
	}
	if !reflect.DeepEqual(au.MapUser, newAU.MapUser) {
		t.Error("Error test to compare structure after writing and readin for AuthUser")
	}

}
func TestGetUserFullName(t *testing.T) {
	var u User
	var db testDB

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

	testPhone = "1"
	err = getUserFullName(testPhone, &u, db)
	if err == nil {
		t.Error("Error in getUserFullName expected error if user not found but return nil", err)
	}

}

func TestGetUserMail(t *testing.T) {
	var u User
	var db testDB
	//	db.rows = new(testRows)
	u.SDId = 12346

	err := getUserMail(&u, db)
	if err != nil {
		t.Error("Error in getUserMail", err)
	}
	if u.Email != "abc@cde_567.com" {
		t.Error("Error returning Email from getUserMail expected abc@cde_567.com but return", u.Email)
	}

	u.SDId = 123
	err = getUserMail(&u, db)
	if err == nil {
		t.Error("Error in getUserMail expected error if user not found but return nil")
	}

}
