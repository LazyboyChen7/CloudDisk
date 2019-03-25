package db

import (
	"time"
	"fmt"
	mydb "filestore_server/db/mysql"
)

// 用户文件表结构体
type UserFile struct {
	UserName string
	FileHash string
	FileName string 
	FileSize int64
	UploadAt string
	LastUpload string
}

// 插入表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt,err := mydb.DBConn().Prepare("insert ignore into tbl_user_file(`user_name`, `file_sha256`, `file_name`,`file_size`,`upload_at`)values(?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_,err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		return false
	}
	return true
}

// 获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt,err := mydb.DBConn().Prepare("select file_sha256,file_name,file_size,upload_at,last_update from tbl_user_file where user_name = ? limit ?")
	if err != nil {
		return nil,err
	}
	defer stmt.Close()

	rows,err := stmt.Query(username, limit)
	if err != nil {
		return nil,err
	}

	var userfiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpload)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userfiles = append(userfiles, ufile)
	}
	return userfiles, nil
}