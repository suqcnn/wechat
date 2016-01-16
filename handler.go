package wechat

import (
	"errors"
	"github.com/chanxuehong/wechat/mp/media"
	"github.com/chanxuehong/wechat/mp/message/response"
	"log"
	"reflect"
	"strings"
)

const (
	MATCHTYPE_PREFIX = "PREFIX"
	MATCHTYPE_EQUAL  = "EQUAL"
	MATCHTYPE_REGEXP = "REGEXP"
	MATCHTYPE_CONTEX = "CONTEX"
)

type Helper struct {
	MatchType  string
	MatchValue string
	Operate    interface{}
	OperArgs   []interface{}
}

var GHelpers []Helper = []Helper{}

type HandleBean struct {
	UserID     string
	ServID     string
	Createtime int64
	Content    string
	RawContent string
	ReqObj     interface{}
}

func (h *HandleBean) Unmarshal(ret_i interface{}) (interface{}, error) {
	ret_type := reflect.TypeOf(ret_i).String()
	switch ret_type {
	case "string":
		resp_str, _ := ret_i.(string)
		if resp_str == "" {
			return nil, errors.New("no response")
		} else {
			return response.NewText(h.UserID, h.ServID, h.Createtime, resp_str), nil
		}
	case "*media.MediaInfo":
		media_info, _ := ret_i.(*media.MediaInfo)
		return response.NewImage(h.UserID, h.ServID, h.Createtime, media_info.MediaId), nil
	default:
		return nil, errors.New("unknown type:" + ret_type)
	}
}

func (h *HandleBean) Response() (interface{}, error) {
	h.Content = h.RawContent

	//match
	for _, helper := range GHelpers {
		//MATCHTYPE_EQUAL
		if helper.MatchType == MATCHTYPE_EQUAL &&
			h.RawContent == helper.MatchValue {
			return h.do(helper.Operate, helper.OperArgs...)
		} else if helper.MatchType == MATCHTYPE_PREFIX &&
			strings.HasPrefix(h.RawContent, helper.MatchValue) {
			//MATCHTYPE_PREFIX
			h.Content = strings.TrimPrefix(h.RawContent, helper.MatchValue)
			return h.do(helper.Operate, helper.OperArgs...)
		} else if helper.MatchType == MATCHTYPE_CONTEX && h.MatchContext(helper.MatchValue) {
			//MATCHTYPE_CONTEX
			return h.do(helper.Operate, helper.OperArgs...)
		}
	}
	//return "", errors.New("NO_MATCH_OPER")
	return "", nil
}

func (h *HandleBean) MatchContext(matched_status string) bool {
	user := &User{UserID: h.UserID}
	get_bool, query_err := Orm.Id(h.UserID).Get(user)
	if !get_bool || query_err != nil {
		return false
	}

	if user.LastMsgType == matched_status {
		return true
	} else {
		return false
	}
}

func (h *HandleBean) do(operate interface{}, operargs ...interface{}) (interface{}, error) {
	inter_type := reflect.TypeOf(operate).String()

	if inter_type == "string" {
		return operate.(string), nil
	} else if strings.HasPrefix(inter_type, "func(") {
		func_v := reflect.ValueOf(operate)

		params := make([]reflect.Value, len(operargs)+1)
		params[0] = reflect.ValueOf(h)

		for i, operarg := range operargs {
			params[i+1] = reflect.ValueOf(operarg)
		}

		call_rets := func_v.Call(params)
		if call_rets[1].IsNil() {
			return call_rets[0].Interface(), nil
		} else {
			return call_rets[0].Interface(), call_rets[1].Interface().(error)
		}
	} else {
		return "", errors.New("UN_FUNC:" + inter_type)
	}

	return "", nil
}

func Register(m_type, m_value string, operate interface{}, operargs ...interface{}) {
	log.Println(operargs)
	GHelpers = append(GHelpers, Helper{m_type, m_value, operate, operargs})
}
