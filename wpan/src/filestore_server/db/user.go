package db

import (
	"fmt"
	mydb "filestore_server/db/mysql"
)

// 通过用户名和密码完成user表的注册操作
func UserSignup(username string, password string, phone string) bool {
	stmt,err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`,`user_pwd`,`phone`)values(?,?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err" + err.Error())
		return false
	}
	defer stmt.Close()

	ret,err := stmt.Exec(username, password, phone)
	if err != nil {
		fmt.Println("Failed to insert, err" + err.Error())
		return false
	}
	rowsAffected, err := ret.RowsAffected();
	if nil == err && rowsAffected>0 {
		return true
	}
	return false
}

// 用户登录校验
func UserSignin(username string, encpwd string) bool {
	stmt,err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows,err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// 刷新用户登录token
func UpdateToken(username string, token string) bool {
	stmt,err := mydb.DBConn().Prepare("replace into tbl_user_token(`user_name`, `user_token`)values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_,err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username string
	Email string
	Phone string
	SignupAt string
	LastActive string
	Status int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt,err := mydb.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}