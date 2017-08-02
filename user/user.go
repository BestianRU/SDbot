package user

import (
	"SDbot/cfg"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"os"
	"regexp"
	//import mysqldriver
	_ "github.com/go-sql-driver/mysql"
)

//User is structure for authorized user
type User struct {
	TId      uint64 `json:"tid"`  //telegram user id
	SDId     uint64 `json:"sdid"` //SD user id
	FullName string `json:"fullanme"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

//MapUser map of authorized users with email index
type MapUser map[string]User

//AuthUser is map for authorizesd users
type AuthUser struct {
	MapUser `json:"users"`
}

//NewAuthUser AuthUser
func NewAuthUser(c *cfg.Cfg) (*AuthUser, error) {
	a := new(AuthUser)
	a.MapUser = make(map[string]User, 10)
	f, err := os.OpenFile(c.AuthUser, os.O_RDONLY, os.FileMode(0660))
	if err != nil {
		return nil, err
	}
	err = a.read(f)
	if err != nil {
		return nil, err
	}
	return a, nil
}

//Add new authorized user
func (a *AuthUser) Add(u User, c *cfg.Cfg) error {
	if u.Email == "" {
		return errors.New("User not found email" + u.FullName)
	}
	a.MapUser[u.Email] = u
	f, err := os.OpenFile(c.AuthUser, os.O_RDWR, os.FileMode(0660))
	if err != nil {
		return err
	}
	return a.save(f)
}

//save AuthUser to file
func (a *AuthUser) save(w io.Writer) error {
	jsonAuthUser, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonAuthUser)
	if err != nil {
		return err
	}
	return nil
}

//read AuthUser from file
func (a *AuthUser) read(r io.Reader) error {
	jsonAuthUser := make([]byte, 10000)
	i, err := r.Read(jsonAuthUser)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonAuthUser[:i], a)
	if err != nil {
		return err
	}
	return nil
}

//Delete user by phone from map of authorized users
func (a *AuthUser) Delete(phone string) error {
	u, err := a.GetByPhone(phone)
	if err != nil {
		return err
	}
	delete(a.MapUser, u.Email)
	return nil
}

//GetByPhone find user by phone
func (a *AuthUser) GetByPhone(p string) (User, error) {
	for _, v := range a.MapUser {
		if v.Phone == p {
			return v, nil
		}
	}
	return User{}, errors.New("User isn't authorized")

}

//GetByTId find user by telegram Id
func (a *AuthUser) GetByTId(t uint64) (User, error) {
	for _, v := range a.MapUser {
		if v.TId == t {
			return v, nil
		}
	}
	return User{}, errors.New("User isn't authorized")

}

//GetByEmail find user by email
func (a *AuthUser) GetByEmail(e string) (User, error) {
	if u, ok := a.MapUser[e]; ok {
		return u, nil
	}
	return User{}, errors.New("User isn't authorized")
}

//DBer interface for MySQL DB
type DBer interface {
	Close() error
	Query(query string, args ...interface{}) (rowser, error)
}

type rowser interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type mySQLBackend struct {
	db *sql.DB
	DBer
}

func (db *mySQLBackend) Close() error {
	return db.db.Close()
}

func (db *mySQLBackend) Query(query string, args ...interface{}) (rowser, error) {
	return db.db.Query(query, args...)
}

//newMySQL open mysql connection
func newMySQL(connectionString string) (DBer, error) {
	dbMySQL, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return &mySQLBackend{db: dbMySQL}, err
}

//getUserMail
func getUserMail(u *User, db DBer) error {
	rows, err := db.Query("SELECT email FROM glpi_useremails WHERE users_id=?", u.SDId)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&u.Email)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Email not found")

}

//GetUserFromSQLByPhone getting user data by his phone number
func GetUserFromSQLByPhone(phone string, c *cfg.Cfg) (User, error) {
	db, err := newMySQL(c.M.User + ":" + c.M.Pass + "@tcp(" + c.M.Host + ":" + c.M.Port + ")/" + c.M.Database)
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	var u User
	err = getUserFullName(phone, &u, db)
	if err != nil {
		return User{}, err
	}
	err = getUserMail(&u, db)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

//getUserFullName getting user FullName by his phone number
func getUserFullName(phone string, u *User, db DBer) error {

	rows, err := db.Query("SELECT id,mobile,comment FROM glpi_users WHERE mobile IS NOT NULL AND comment IS NOT NULL")
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&u.SDId, &u.Phone, &u.FullName)
		if err != nil {
			return err
		}
		regExp := regexp.MustCompile("\\D")
		u.Phone = regExp.ReplaceAllString(u.Phone, "")
		if u.Phone == phone {

			return nil
		}
	}
	return errors.New("user not found in SD")
}
