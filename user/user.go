package user

import (
	"SDbot/cfg"
	"database/sql"
	"regexp"
	_ "github.com/go-sql-driver/mysql"
)

//User is structure for authorized user
type User struct {
	TId	int64 //telegram id
	FirstName	string
	LastName	string
	Email		string
	Phone		string
}

//UserMap is map for authorizesd users
type UserMap map[string]User 	

//GetUserFromSQLByPhone Receiving user data by its phone number
func GetUserFromSQLByPhone(phone string, c *cfg.Cfg) (User,error) {
	db, err := sql.Open("mysql", c.M.User+":"+c.M.Pass+"@tcp("+c.M.Host+":"+c.M.Port+")/"+c.M.Database)
	if err!=nil {
		return User{},err
	}
	defer db.Close()
	rows, err := db.Query("SELECT id,mobile,comment FROM glpi_users WHERE mobile IS NOT NULL")
	if err!=nil {
		return User{},err
	}
	for rows.Next() {
		var id int
		var mobile string
		var name string
		rows.Scan(&id,&mobile,&name)
		regExp:=regexp.MustCompile("\\D")
		mobile=regExp.ReplaceAllString(mobile,"")
		if mobile==phone {
			
		}
		println(id,"\t",mobile,"\t",name)
	}
	return User{},nil
}


