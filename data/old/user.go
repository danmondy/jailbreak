package data

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	userSchema = `CREATE table user (id varchar(50) not null primary key, email varchar(50), hashword varchar(100), role integer, since datetime)`
)

type User struct {
	ID       string    `db:"id",`
	Email    string    `db:"email"`
	Hashword string    `db:"hashword"`
	Role     int       `db:"role"`
	Since    time.Time `db:"since"`
}

func NewUser(email string, password string, role int) User {
	u := User{Email: email, Role: role}
	u.SetHashword([]byte(password))
	u.Since = time.Now()
	return u
}

func (u *User) SetHashword(password []byte) error {
	hash, err := bcrypt.GenerateFromPassword(password, 0) //0 -> uses default cost (10 at time of writing, 4 min / 31 max)
	if err != nil {
		return err
	} else {

		u.Hashword = base64.StdEncoding.EncodeToString(hash)
		return nil
	}
}

func (u *User) CompareHash(p []byte) error {
	hash, err := base64.StdEncoding.DecodeString(u.Hashword)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(hash, p)
}

func GetAllUsers() ([]User, error) {
	users := []User{}
	fmt.Println("getting all users")
	err := db.Select(&users, "SELECT * FROM user")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(users)
	return users, nil
}

func GetUserByID(id string) (User, error) {
	user := User{}
	err := db.Get(&user, "SELECT * FROM user where id = ?", id)
	return user, err
}

func GetUserByEmail(email string) (User, error) {
	fmt.Println(email)
	user := User{}
	err := db.Get(&user, "SELECT * FROM user WHERE email=?", email)
	return user, err
}

func UpdateUser(user User) error {
	log.Println("UPDATING USER", user)
	_, err := db.NamedExec("UPDATE user set email = :email, hashword = :hashword, role = :role where id = :id;", user)
	log.Println("UPDATED USER", user)
	return err
}

func DeleteUser(id string) error {
	_, err := db.Exec("delete from user where id = ?", id)
	return err
}

/*


package data

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int       `db:"id,primarykey,autoinsert",dbname:"user"`
	Email    string    `db:"email"`
	Hashword string    `db:"hashword"`
	Role     int       `db:"role"`
	Since    time.Time `db:"since"`
}

func NewUser(email string, password string, role int) User {
	u := User{Email: email, Role: role}
	u.SetHashword([]byte(password))
	u.Since = time.Now()
	return u
}

func (u *User) SetHashword(password []byte) error {
	hash, err := bcrypt.GenerateFromPassword(password, 0) //0 -> uses default cost (10 at time of writing, 4 min / 31 max)
	if err != nil {
		return err
	} else {
		u.Hashword = string(hash)
		return nil
	}
}

func (u *User) CompareHash(p []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Hashword), p)
}

func GetAllUsers() ([]User, error) {
	users := []User{}
	err := db.Select(&users, "Select * FROM user")
	if err != nil {
		return nil, err
	}
	fmt.Println(users)
	return users, nil
}

func GetUserById(id int) (*User, error) {
	query := "SELECT * FROM user WHERE Id = ?"
	row := db.QueryRow(query, id)
	return MapUser(row)
}

func InsertUser(u *User) error {
	fmt.Println(u)
	_, err := db.Exec(fmt.Sprintf("INSERT into user (email, hashword, role, since) VALUES ('%s', '%s', %v, '%s')", u.Email, u.Hashword, u.Role, u.Since))
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*User, error) {
	query := "SELECT * FROM user WHERE EMAIL = ?"
	row := db.QueryRow(query, email)
	return MapUser(row)
}

func MapUser(r *sql.Row) (*User, error) {
	var u User
	var t string
	err := r.Scan(&u.Id, &u.Email, &u.Hashword, &u.Role, &t)
	if err != nil {
		return nil, err
	}
	u.Since, err = StringToTime(t)
	return &u, err
}

func MapUsers(r *sql.Rows) ([]User, error) {
	var users []User
	for r.Next() {
		var u User
		err := r.Scan(&u.Id, &u.Email, &u.Hashword, &u.Role, &u.Since)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

/*
func NewUserProvider() {
	return Provider{
		Q_GetById:    "SELECT * FROM users WHERE EMAIL = ?",
		Q_GetByField: "SELECT * FROM users where %s = ?",
		//Q_Insert:            "INSERT into users (email, hashword, role, since) VALUES ('%s', '%s', %v, '%s')",
		MapMany:   MapUsers,
		MapSingle: MapUser,
	}
}
*/
//Experimental
/*type MapOne func(r *sql.Row) (interface{}, error)
type MapMany func(r *sql.Rows) (interface{}, error)

type Provider struct {
	//	Q_Insert     string
	Q_GetById    string
	Q_GetByField string
	MapSingle    MapOne
	MapMany      MapMany
}

func (p Provider) GetById(id string) (interface{}, error) {
	row := db.QueryRow(p.Q_GetById, id)
	return p.MapOne(row)
}
func (p Provider) GetByField(fieldName string, value string) (interface{}, error) {
	row := db.QueryRow(fmt.Sprintf(query, fieldName), value)
	return p.MapOne(row)
}
*/
/*
/*func (p Provider) Insert() error {
	fmt.Println(u)
	_, err := db.Exec(fmt.Sprintf(Q_Insert, u.Email, u.Hashword, u.Role, u.Since))
	if err != nil {
		return err
	}
	return nil
}*/
