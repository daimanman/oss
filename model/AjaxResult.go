package model

const (
	AJAX_STATUS_CODE_SUCCESS int = 0
	AJAX_STATUS_CODE_WARN    int = 1
	AJAX_STATUS_CODE_ERROR   int = 2
)

//返回给客户端的对象信息
type AjaxResult struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
}

func (result *AjaxResult) ErrResult(msg string) *AjaxResult {
	result.Message = msg
	result.StatusCode = AJAX_STATUS_CODE_ERROR
	return result
}

func (result *AjaxResult) OkResult(msg string, data interface{}) *AjaxResult {
	result.Message = msg
	result.Data = data
	result.StatusCode = AJAX_STATUS_CODE_SUCCESS
	return result
}
