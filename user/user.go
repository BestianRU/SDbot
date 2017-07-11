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
	rows, err := db.Query("SELECT id,mobile,comment FROM glpi_users WHERE mobile IS NOT NULL AND comment IS NOT NULL")
	if err!=nil {
		return User{},err
	}
	for rows.Next() {
		var id int
		//var mobile string
		//var name string
		//var mail string
		var u User
		err=rows.Scan(&id,&u.Phone,&u.FirstName)
		if err!=nil {
				return User{},err
		}
		regExp:=regexp.MustCompile("\\D")
		u.Phone=regExp.ReplaceAllString(u.Phone,"")
		if u.Phone==phone {
			err = db.QueryRow("SELECT email FROM glpi_useremails WHERE users_id=?",id).Scan(&u.Email)
			if err!=nil {
				return User{},err
			}
			
			return User{},err
			println(u.Email,"\t",u.Phone)
			
		}
	//	println(id,"\t",mobile,"\t",name)
	
	}
	return User{},nil
}


