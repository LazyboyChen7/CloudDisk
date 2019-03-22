package meta

import ( 
	mydb "filestore_server/db"
)

// 文件元信息结构
type FileMeta struct {
	Filesha256 string
	FileName string
	FileSize int64
	FileLocation string
	FileUploadTime string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 新增/更新文件元数据 到mysql中
func UpdataFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.Filesha256, fmeta.FileName, fmeta.FileSize, fmeta.FileLocation)
}

// 从mysql获取文件元信息
func GetFileMetaDB(filesha256 string) (FileMeta, error) {
	tfile,err := mydb.GetFileMeta(filesha256)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		Filesha256 : tfile.FileHash,
		FileName : tfile.FileName.String,
		FileSize : tfile.FileSize.Int64,
		FileLocation : tfile.FileAddr.String,
	}
	return fmeta, nil
}

// 新增/更新文件元数据
func UpdataFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.Filesha256] = fmeta
}

// 通过哈希获取文件元对象
func GetFileMeta(FileHash string) FileMeta {
	return fileMetas[FileHash]
}

// 删除文件
func RemoveFileMeta(filesha256 string) {
	delete(fileMetas, filesha256)
}