package wechat

import (
	//"fmt"
	//"git.oschina.net/xuebing1110/queryapi"
	"errors"
	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/chanxuehong/wechat/mp/message/request"
	"github.com/chanxuehong/wechat/mp/message/response"
	"log"
	"net/http"
	"reflect"
)

const (
	DEFAULT_TEXT = "●注册姓名：回复格式【姓名@**】，例如“姓名@张三”；\n●上传头像：发送【头像】，然后上传图片；\n●定位位置：发送位置；\n●琴岛通绑定：发送【琴岛通@xxxxxxxx】绑定琴岛通卡；\n●公交卡余额查询：发送【公交卡余额】查询余额。"
	MAP_URL      = "http://123.56.66.219/amap/map.html"
	HELP_URL     = "http://mp.weixin.qq.com/s?__biz=MzAxMDYwMDM1Mw==&mid=400839956&idx=1&sn=bb380e0ebc832851a570ed8673aee169&scene=18#wechat_redirect"
)

func GetHandleBean(req_msg interface{}) (*HandleBean, error) {
	handle_bean := &HandleBean{ReqObj: req_msg}
	req_type := reflect.TypeOf(req_msg).String()

	switch req_type {
	case "*request.Text":
		text, ok := req_msg.(*request.Text)
		if !ok {
			log.Println("parse to request.Text error!")
		}

		handle_bean.UserID = text.FromUserName
		handle_bean.ServID = text.ToUserName
		handle_bean.Createtime = text.CreateTime
		handle_bean.RawContent = text.Content
		return handle_bean, nil
	default:
		return handle_bean, errors.New("unknown type:" + req_type)
	}
}

// 文本消息的 Handler
func TextHandler(w http.ResponseWriter, r *mp.Request) {
	log.Println(string(r.RawMsgXML))
	handle_bean, hb_err := GetHandleBean(request.GetText(r.MixedMsg))
	if hb_err != nil {
		return
	}

	resp_ret, err := handle_bean.Response()
	if err != nil {
		log.Println(err.Error())
		resp_ret = "你这话我没法接"
	}

	resp, resp_err := handle_bean.Unmarshal(resp_ret)
	if resp_err != nil {
		resp, resp_err = APIQuery.Query(handle_bean.UserID, "TURING://"+handle_bean.Content)
		if resp_err == nil {
			resp, resp_err = handle_bean.Unmarshal(resp)
		}
	}

	if resp_err == nil {
		mp.WriteAESResponse(w, r, resp)
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}

/*
func ImageHandler(w http.ResponseWriter, r *mp.Request) {
	log.Println(string(r.RawMsgXML))
	image := request.GetImage(r.MixedMsg)

	//初始化变量
	var err error
	var tip_content string

	//查询用户信息
	user := &User{UserID: image.FromUserName}
	get_bool, get_err := Orm.Id(image.FromUserName).Get(user)

	//上文为上传头像操作"@头像"
	if get_err == nil && get_bool && user.LastMsgType == MSGTYPE_REG_HEADPIC_TEXT {
		user.HeadUrl = image.PicURL
		user.LastMsgType = MSGTYPE_REG_HEADPIC_IMG
		tip_content, err = user.UploadHeadPic()
		if err != nil {
			log.Println(err)
		}
	} else {

	}

	//上下文操作
	if tip_content != "" {
		resp := response.NewText(image.FromUserName, image.ToUserName, image.CreateTime, tip_content)
		mp.WriteAESResponse(w, r, resp)
	} else {
		//默认图片操作(上传设备图片)
		user = &User{UserID: image.FromUserName}
		mp.WriteAESResponse(w, r, user.UploadDevicePic(image))
	}
}

func LocationHandler(w http.ResponseWriter, r *mp.Request) {
	log.Println(string(r.RawMsgXML))
	location := request.GetLocation(r.MixedMsg)

	//初始化变量
	var err error
	var tip_content, tip_descr, tip_url string

	//查询用户信息
	user := &User{UserID: location.FromUserName}
	get_bool, get_err := Orm.Id(location.FromUserName).Get(user)

	if get_err != nil {
		tip_content = "定位失败！"
		tip_descr = "查询用户信息失败：" + get_err.Error()
	} else if !get_bool || user.Name == user.UserID {
		tip_content = "定位失败！"
		tip_descr = "请先注册用户信息!"
		tip_url = HELP_URL
	} else {
		user.Location = fmt.Sprintf("%v,%v", location.LocationY, location.LocationX)
		user.Address = location.Label
		tip_content, err = user.UploadLocation()
		if err != nil {
			log.Println(err)
		} else {
			tip_descr = fmt.Sprintf("%s\n街道:%s\n点击看地图！", location.Label, user.DetailAddr)
			tip_url = MAP_URL
		}
	}

	resp := response.NewNews(
		location.FromUserName,
		location.ToUserName,
		location.CreateTime,
		[]response.Article{
			response.Article{
				Title:       tip_content,
				Description: tip_descr,
				URL:         tip_url}})
	mp.WriteAESResponse(w, r, resp)
}
*/
func UnsubscribeHandler(w http.ResponseWriter, r *mp.Request) {
	log.Println(string(r.RawMsgXML))

	msg_header := request.GetSubscribeEvent(r.MixedMsg).MessageHeader
	log.Printf(`the user "%s" Unsubscribe!!!`, msg_header.FromUserName)
}

func MenuHandler(w http.ResponseWriter, r *mp.Request) {
	log.Println(string(r.RawMsgXML))
	click_e := menu.GetClickEvent(r.MixedMsg)

	content := "收到点击菜单事件！"
	resp := response.NewText(click_e.FromUserName, click_e.ToUserName, click_e.CreateTime, content)
	mp.WriteAESResponse(w, r, resp)
}

func SubscribeHandler(w http.ResponseWriter, r *mp.Request) {
	log.Println(string(r.RawMsgXML))
	msg_header := request.GetSubscribeEvent(r.MixedMsg).MessageHeader
	log.Printf(`the user "%s" subscribe!!!`, msg_header.FromUserName)

	msg := request.GetText(r.MixedMsg)
	/*resp := response.NewNews(
	msg.FromUserName,
	msg.ToUserName,
	msg.CreateTime,
	[]response.Article{
		response.Article{
			Title:       "使用帮助说明",
			Description: "欢迎关注TalkAbout公众号，点击请看使用帮助...",
			PicURL:      "http://mmbiz.qpic.cn/mmbiz/oITYCaq3I7fvWXu5LYreDhxEgvPM2NZHzGedlk79GzFRsuV8XYxSU1usEibLqIROpQib0YccnlUuL6icYSzicJiaBRA/640?wx_fmt=jpeg&tp=webp&wxfrom=5",
			URL:         "http://mp.weixin.qq.com/s?__biz=MzAxMDYwMDM1Mw==&mid=400785377&idx=1&sn=adbc968611ba1889f830fa0154b45025&scene=18#wechat_redirect"}})
	*/
	content := "请先注册个人信息：\n※注册姓名：回复格式“姓名**”，例如“姓名张三”；\n※上传头像：发送“头像”，然后上传图片；\n※定位位置：发送位置；"
	resp := response.NewText(msg.FromUserName, msg.ToUserName, msg.CreateTime, content)
	mp.WriteAESResponse(w, r, resp)
}
