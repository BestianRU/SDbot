package user

import (
	"SDbot/cfg"
	"database/sql"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

//User is structure for authorized user
type User struct {
	TId       int64 //telegram id
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

//UserMap is map for authorizesd users
type UserMap map[string]User

//DBer interface for MySQL DB
type DBer interface {
	//	Ping() error
	Close() error
	//	Execute(query string, args ...interface{}) error
	Query(query string, args ...interface{}) (rowsScanner, error)
	QueryRow(query string, args ...interface{}) scanner
}

type rowsScanner interface {
	//	Columns() ([]string, error)
	Next() bool
	//	Close() error
	//	Err() error
	scanner
}

type scanner interface {
	Scan(dest ...interface{}) error
}

type mySQLBackend struct {
	db *sql.DB
	DBer
}

func (db *mySQLBackend) Close() error {
	return db.db.Close()
}

func (db *mySQLBackend) Query(query string, args ...interface{}) (rowsScanner, error) {
	return db.db.Query(query, args...)
}

func (db *mySQLBackend) QueryRow(query string, args ...interface{}) scanner {
	return db.db.QueryRow(query, args...)
}

//newMySQL open mysql connection
func newMySQL(connectionString string) (DBer, error) {
	dbMySQL, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return &mySQLBackend{db: dbMySQL}, err
}

//GetUserFromSQLByPhone Receiving user data by its phone number
func GetUserFromSQLByPhone(phone string, c *cfg.Cfg) (User, error) {
	db, err := newMySQL(c.M.User + ":" + c.M.Pass + "@tcp(" + c.M.Host + ":" + c.M.Port + ")/" + c.M.Database)
	if err != nil {
		return User{}, err
	}
	u, err := getUserFromSQLByPhone(phone, db)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

//getUserFromSQLByPhone Receiving user data by its phone number
func getUserFromSQLByPhone(phone string, db DBer) (User, error) {
	defer db.Close()
	rows, err := db.Query("SELECT id,mobile,comment FROM glpi_users WHERE mobile IS NOT NULL AND comment IS NOT NULL")
	if err != nil {
		return User{}, err
	}
	for rows.Next() {
		var id int
		var u User
		err = rows.Scan(&id, &u.Phone, &u.FirstName)
		if err != nil {
			return User{}, err
		}
		regExp := regexp.MustCompile("\\D")
		u.Phone = regExp.ReplaceAllString(u.Phone, "")
		if u.Phone == phone {
			err = db.QueryRow("SELECT email FROM glpi_useremails WHERE users_id=?", id).Scan(&u.Email)
			if err != nil {
				return User{}, err
			}
			return u, err
		}
	}
	return User{}, nil
}
