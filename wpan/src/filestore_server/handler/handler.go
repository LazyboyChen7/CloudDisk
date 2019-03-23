package handler

import (
	"net/http"
	"io/ioutil"
	"io"
	"os"
	"fmt"
	"time"
	//"strconv"
	//dblayer "filestore_server/db"
	"encoding/json"
	"filestore_server/meta"
	"filestore_server/util"
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data,err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, err.Error())
			return 
		}
		io.WriteString(w, string(data)) 
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data, err:%s", err.Error())
			return 
		}
		defer file.Close()
		
		fileMeta := meta.FileMeta{
			FileName : head.Filename,
			FileLocation : "/tmp/" + head.Filename,
			FileUploadTime : time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.FileLocation)
		if err != nil {
			fmt.Printf("Failed to create file, err:%s", err.Error())
			return 
		}
		defer newFile.Close()

		fileMeta.FileSize,err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data to file, err:%s", err.Error())
			return 
		}

		newFile.Seek(0, 0)
		fileMeta.Filesha256 = util.FileSha256(newFile)
		//meta.UpdataFileMeta(fileMeta)
		_ = meta.UpdataFileMetaDB(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Finished!")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	fMeta,err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	data,err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.Write(data)
}

// 下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha256 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha256)

	f,err := os.Open(fm.FileLocation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-disposition", "attachment;filename=\"" + fm.FileName + "\"")
	w.Write(data)
}

// 更新文件元信息
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	filesha256 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return 
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return 
	}

	curFileMeta := meta.GetFileMeta(filesha256)
	curFileMeta.FileName = newFileName
	meta.UpdataFileMeta(curFileMeta)

	data,err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 删除文件及元信息
func FileDelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filesha256 := r.Form.Get("filehash")
	
	fMeta := meta.GetFileMeta(filesha256)
	os.Remove(fMeta.FileLocation)
	
	meta.RemoveFileMeta(filesha256)

	w.WriteHeader(http.StatusOK)
}

// FileQueryHandler : 查询批量的文件元信息
// func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
// 	username := r.Form.Get("username")
// 	//fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
// 	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	data, err := json.Marshal(userFiles)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(data)
// }