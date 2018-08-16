package handler

import (
	"github.com/naoina/denco"
	"io"
	"log"
	"net/http"
	"net/url"
	"oss/model"
	"oss/service"
	"strings"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

//index 首页 测试
func IndexHandler(w http.ResponseWriter, r *http.Request, params denco.Params) {
	io.WriteString(w, " Hello World ")
}

//上传文件
func UploadFileHandler(w http.ResponseWriter, r *http.Request, params denco.Params) {

	r.ParseMultipartForm(32 << 20)
	uploadfile, header, _ := r.FormFile("file")

	fileSize := header.Size
	if fileSize == 0 {
		SendErrText(w, "请选择上传的文件")
		return
	}

	if fileSize > 128*1024*1024 {
		SendErrText(w, "文件大小超出128M")
		return
	}

	fid, err := service.OssUploadFile(uploadfile, header)
	if err != nil {
		log.Printf("upload err %v\n", err)
		SendErrText(w, "服务器开小差！")
		return
	}

	log.Printf("[ fileid=%s filename=%s filesize=%d ] \n", fid, header.Filename, header.Size)

	SendOkText(w, fid)
}

func sendFile(w http.ResponseWriter, r *http.Request, params denco.Params, formParams map[string]string) {
	fid := params.Get("fileid")
	filename, fileReader, err := service.OssDownLoadFile(fid, formParams)
	log.Printf("fileid=%s filename=%s \n", fid, filename)
	if err != nil || fid == "" {
		log.Printf("download fileid=%s err=%v \n", fid, err)
		SendErrText(w, "文件不存在")
		return
	}
	defer fileReader.Close()

	filename, mimeType := service.GetMimeType(fid)
	filename = url.QueryEscape(filename)
	w.Header().Set("Content-Disposition", "inline;filename=\""+filename+"\"")
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	io.Copy(w, fileReader)
}

//下载文件
func DownFileHandler(w http.ResponseWriter, r *http.Request, params denco.Params) {
	r.ParseForm()
	formParams := make(map[string]string)
	formParams["height"] = r.FormValue("height")
	formParams["width"] = r.FormValue("width")
	sendFile(w, r, params, formParams)
}

//下载缩略图
func DownThumbFile(w http.ResponseWriter, r *http.Request, params denco.Params) {
	r.ParseForm()
	width := r.FormValue("width")
	if width == "" {
		width = "400"
	}

	formParams := make(map[string]string)
	formParams["width"] = width
	formParams["height"] = r.FormValue("height")
	sendFile(w, r, params, formParams)
}

//删除文件
func DeleteFileHandler(w http.ResponseWriter, r *http.Request, params denco.Params) {
	fileid := params.Get("fileid")
	result := &model.AjaxResult{}
	if fileid == "" {
		SendErrJson(w, result.ErrResult("缺失文件id,操作失败"))
		return
	}

	fids := strings.Split(fileid, "@")
	err := service.OssDeleteFiles(fids)

	if err != nil {
		log.Printf("文件删除失败 fileid=%s err=%v", fileid, err)
		SendErrJson(w, result.ErrResult("服务器开小差!"))
		return
	}
	SendOkJson(w, result.OkResult("操作成功", fileid))
}

//批量删除文件
func DeleteBatchFileHandler(w http.ResponseWriter, r *http.Request, params denco.Params) {
	DeleteFileHandler(w, r, params)
}

//获取文件名称
func GetFileName(w http.ResponseWriter, r *http.Request, params denco.Params) {
	r.ParseForm()
	fileid := r.FormValue("fileid")
	if fileid == "" {
		SendErrText(w, "缺失文件id参数")
		return
	}

	filename := service.OssGetFileName(fileid)
	if filename == "" {
		SendErrText(w, "文件不存在")
		return
	}
	SendOkText(w, filename)
}
