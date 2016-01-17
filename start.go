package wechat

import (
	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/message/request"
	"github.com/chanxuehong/wechat/util"
	"net/http"
	"strconv"
)

//var AccessTokenServer *mp.DefaultAccessTokenServer

func (config *Configure) Listen() error {
	err := http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
	if err != nil {
		return err
	}
	return nil
}

func (config *Configure) InitWechatAndListen() error {
	err := config.InitWechat()
	if err != nil {
		return err
	}

	return config.Listen()
}

func (config *Configure) InitWechat() error {
	aesKey, err := util.AESKeyDecode(config.EncodingAESKey)
	if err != nil {
		return err
	}

	//AccessTokenServer = mp.NewDefaultAccessTokenServer(config.AppID, config.AppSecret, nil)

	messageServeMux := mp.NewMessageServeMux()
	messageServeMux.MessageHandleFunc(request.MsgTypeText, TextHandler)
	//messageServeMux.MessageHandleFunc(request.MsgTypeImage, ImageHandler)
	//messageServeMux.MessageHandleFunc(request.MsgTypeLocation, LocationHandler)
	messageServeMux.EventHandleFunc(request.EventTypeSubscribe, SubscribeHandler)
	messageServeMux.EventHandleFunc(request.EventTypeUnsubscribe, UnsubscribeHandler)
	//messageServeMux.EventHandleFunc(menu.EventTypeClick, MenuHandler)
	//messageServeMux.EventHandleFunc(menu.EventTypeView, MenuHandler)

	// 下面函数的几个参数设置成你自己的参数: oriId, token, appId
	mpServer := mp.NewDefaultServer(
		config.OriID,
		config.Token,
		config.AppID,
		aesKey,
		messageServeMux)

	mpServerFrontend := mp.NewServerFrontend(mpServer, mp.ErrorHandlerFunc(ErrorHandler), nil)

	// 如果你在微信后台设置的回调地址是
	//log.Printf("listen port %v ...", config.Port)
	if config.UrlBase == "" {
		http.Handle("/", mpServerFrontend)
	} else {
		http.Handle(config.UrlBase, mpServerFrontend)
	}

	return config.Listen()
}
