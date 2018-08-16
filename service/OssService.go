package service

import (
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"oss/config"
	gosw "oss/goseaweed"
	"path"
	"regexp"
	"strings"
)

var sw *gosw.Seaweed
var filenameRegexp = regexp.MustCompile(`.*filename="(.*)".*`)

func init() {
	sw = gosw.NewSeaweed(config.SeaWeedMasterHost)
}
func OssUploadFile(reader io.Reader, header *multipart.FileHeader) (string, error) {
	fid, err := sw.UploadMultiPartFile(reader, header)
	if err != nil {
		log.Printf("upload file err %v \n", err)
	}
	return fid, err
}

func GetMimeType(fid string) (filename string, mimeType string) {
	filename = OssGetFileName(fid)
	ext := strings.ToLower(path.Ext(filename))
	if ext != "" {
		mimeType = mime.TypeByExtension(ext)
	}
	return filename, mimeType
}

func OssGetFileUrl(fid string, params map[string]string) string {
	return sw.GetFileDownloadUrl(fid, params)
}

func OssDownLoadFile(fid string, params map[string]string) (string, io.ReadCloser, error) {
	return sw.HC.DownloadUrl(OssGetFileUrl(fid, params))
}

func OssGetFileName(fid string) string {
	url := OssGetFileUrl(fid, nil)
	response, err := sw.HC.Client.Head(url)
	if err != nil {
		log.Printf("fid=%s url=%s err=%v\n", fid, url, err)
		return ""
	}
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		log.Printf("http error code = %d\n", response.StatusCode)
		return ""
	}
	contentDisposition := response.Header["Content-Disposition"]
	if len(contentDisposition) > 0 {
		strs := filenameRegexp.FindStringSubmatch(contentDisposition[0])
		if len(strs) > 1 {
			return strs[1]
		}
	}
	return ""
}

func OssDeleteFiles(fids []string) error {
	if len(fids) == 0 {
		return nil
	}
	for _, fid := range fids {
		if err := sw.DeleteFile(fid, ""); err != nil {
			log.Printf("delete file fid=%s err=%v \n", fid, err)
			return err
		}
	}
	return nil
}
