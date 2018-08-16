package handler

import (
	"encoding/json"
	"github.com/naoina/denco"
	"io"
	"log"
	"net/http"
	"oss/config"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

}

func sendJson(w http.ResponseWriter, data interface{}) {
	setJsonHeader(w)
	json.NewEncoder(w).Encode(data)
}
func SendErrJson(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	sendJson(w, data)
}
func SendOkJson(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	sendJson(w, data)
}

func SendOkText(w http.ResponseWriter, data string) {
	setJsonHeader(w)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, data)
}

func SendErrText(w http.ResponseWriter, data string) {
	setJsonHeader(w)
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, data)
}

func getOssPath(basePath string) string {
	prefix := "/api/v1"
	strPath := config.GetHandlerPath(prefix + basePath)
	log.Printf("ROUTE PATH %s \n", strPath)
	return strPath
}

func StartOssServer() {
	mux := denco.NewMux()
	muxhandler, err := mux.Build([]denco.Handler{
		mux.GET("/", IndexHandler),
		mux.POST(getOssPath("/oss"), UploadFileHandler),
		mux.POST(getOssPath("/oss/"), UploadFileHandler),

		mux.GET(getOssPath("/oss/:fileid"), DownFileHandler),
		mux.POST(getOssPath("/oss/:fileid"), DownFileHandler),

		mux.GET(getOssPath("/oss/delete/:fileid"), DeleteFileHandler),
		mux.POST(getOssPath("/oss/delete/:fileid"), DeleteFileHandler),

		mux.GET(getOssPath("/oss/deletemulti/:fileid"), DeleteBatchFileHandler),
		mux.POST(getOssPath("/oss/deletemulti/:fileid"), DeleteBatchFileHandler),

		mux.GET(getOssPath("/oss/thumb/:fileid"), DownThumbFile),
		mux.POST(getOssPath("/oss/thumb/:fileid"), DownThumbFile),

		mux.GET(getOssPath("/oss/filename"), GetFileName),
		mux.POST(getOssPath("/oss/filename"), GetFileName),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("APP Listening PORT %s \n", config.AppPort)

	log.Fatal(http.ListenAndServe(config.AppPort, muxhandler))
}
