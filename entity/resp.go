package entity

import (
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/utils"
)

type (

	JsonResponse struct {
		// httpCode
		HttpCode int `json:"ret"`
		// Data
		Data *JsonData `json:"data,omitempty"`
	}

	JsonData struct {
		Code int     `json:"code"`
		Msg  string  `json:"msg,omitempty"`
		Info []KvMap `json:"info,omitempty"`
	}
)

const (
	CodeOk     = 0
	CodeFailed = -1
)

// CreateResponse 创建json 响应
func CreateResponse(content ...[]byte) *JsonResponse {
	var resp = &JsonResponse{}
	if len(content) == 0 {
		resp.Data = &JsonData{
			Code: 0,
			Msg:  "OK",
		}
		resp.HttpCode = 200
		return resp
	}
	err := utils.JsonDecode(content[0], resp)
	if err != nil {
		log.Infoln("CreateResponse error:", err)
	}
	return resp
}

func (resp *JsonResponse) GetData() JsonData {
	if resp.Data != nil {
		return *resp.Data
	}
	return JsonData{}
}

func (resp *JsonResponse) GetMsg() string {
	if resp.Data != nil {
		return resp.Data.Msg
	}
	return ""
}

func (resp *JsonResponse) GetCode() int {
	if resp.Data != nil {
		return resp.Data.Code
	}
	return CodeFailed
}

func (resp *JsonResponse) Empty() bool {
	if resp.Data == nil {
		return true
	}
	if resp.Data.Info == nil || len(resp.Data.Info) == 0 {
		return true
	}
	return false
}

func (resp *JsonResponse) IsSuccess() bool {
	if resp.HttpCode != 200 {
		return false
	}
	if resp.Data != nil && resp.Data.Code == CodeOk {
		return true
	}
	return false
}

func (resp *JsonResponse) Decode(data interface{}) error {
	var json = utils.JsonEncode(resp.Data)
	if json.HasErr() {
		return json.Error()
	}
	if json.Empty() {
		return json.EmptyErr()
	}
	return json.Decode(data)
}

func (resp *JsonResponse) First() KvMap {
	if resp.Data != nil && resp.Data.Info != nil {
		return resp.Data.Info[0]
	}
	return nil
}

func (resp *JsonResponse) Last() KvMap {
	if resp.Data != nil && resp.Data.Info != nil {
		return resp.Data.Info[len(resp.Data.Info)-1]
	}
	return nil
}

func (resp *JsonResponse) IndexOf(i int) KvMap {
	if resp.Data != nil && resp.Data.Info != nil {
		var size = len(resp.Data.Info)
		if size > i {
			return resp.Data.Info[i]
		}
	}
	return nil
}

func (resp *JsonResponse) Count() int {
	if resp.Data != nil && resp.Data.Info != nil {
		var size = len(resp.Data.Info)
		return size
	}
	return 0
}

func (resp *JsonResponse) DecodeData(v interface{}) error {
	if resp.Data != nil && resp.Data.Info != nil {
		var size = len(resp.Data.Info)
		if size == 1 {
			json := utils.JsonEncode(resp.Data.Info[0])
			return utils.JsonDecode(json.Bytes(), v)
		} else {
			json := utils.JsonEncode(resp.Data.Info)
			return utils.JsonDecode(json.Bytes(), v)
		}
	}
	return nil
}

func (resp *JsonResponse) SizeOf() int {
	if resp.Data != nil && resp.Data.Info != nil {
		return len(resp.Data.Info)
	}
	return 0
}

func (resp *JsonResponse) String() string {
	return utils.JsonEncode(resp).String()
}
