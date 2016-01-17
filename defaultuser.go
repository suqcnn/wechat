package wechat

import (
	//"fmt"
	"errors"
	"time"
)

type User struct {
	UserID          string    `json:"userid" xorm:"PK 'userid'"`
	Name            string    `json:"_name" xorm:"'name'"`
	Location        string    `json:"_location" xorm:"'location'"`
	Address         string    `json:"_address" xorm:"'address'`
	DetailAddr      string    `json:"detailaddress" xorm:"'detailaddress'"`
	TelePhone       string    `json:"telephone" xorm:"'telephone'"`
	HeadUrl         string    `json:"headurl" xorm:"'headurl'"`
	Aoi             string    `json:"aoi" xorm:"'aoi'"`
	AmapID          string    `json:"_id" xorm:"'amapid'"`
	ChangeTime      time.Time `json:"changetime" xorm:"updated 'changetime'"`
	Groups          string    `json:"groups" xorm:"'groups'"`
	LastMsgType     string    `json:"lastmsgtype" xorm:"'lastmsgtype'"`
	Tally           int       `json:"tally" xorm:"'tally'"`
	SubscribeTime   time.Time `json:"subscribetime" xorm:"notnull created 'subscribetime'"`
	UnSubscribeTime time.Time `json:"unsubscribetime" xorm:"deleted 'unsubscribetime'"`
}

func (u *User) TableName() string {
	return "user"
}

func (user *User) IncTally() error {
	use_db := &User{}
	find_bool, query_err := Orm.Id(user.UserID).Get(use_db)
	if !find_bool || query_err != nil {
		return query_err
	}

	user.Tally = use_db.Tally + 1
	_, err := Orm.Id(user.UserID).Cols("tally").Update(user)
	if err != nil {
		return err
	}

	return nil
}

const (
	FMT_REG_NAME_FAIL        = `尊敬的%s，注册姓名失败，失败编码：%s，请联系公众号管理员！`
	FMT_REG_NAME_SUC         = `尊敬的%s，注册姓名成功！`
	FMT_REG_HEADPIC_TEXT_SUC = `请上传头像照片！`
	FMT_REG_HEADPIC_FAIL     = `上传图片失败，失败编码：%s，请联系公众号管理员！`

	MSGTYPE_NULL            = "无"
	MSGTYPE_REG_NAME        = "@姓名"
	MSGTYPE_REG_PHONE       = "@手机"
	MSGTYPE_REG_HEADPIC_IMG = "@头像完成"
)

func (user *User) SaveStatus() error {
	user.Tally = 0
	rownum, err := Orm.Id(user.UserID).Cols("lastmsgtype", "tally").Update(user)
	if err != nil {
		return errors.New("UPDATE_LASTMSGTYPE_ERR:" + err.Error())
	}

	if rownum != 1 {
		rownum, err = Orm.InsertOne(user)
		if err != nil || rownum != 1 {
			return errors.New("INSERT_LASTMSGTYPE_ERR:" + err.Error())
		}
	}

	return nil
}

func (user *User) UploadHeadPic() error {
	rownum, err := Orm.Id(user.UserID).Cols("headurl", "lastmsgtype").Update(user)
	if err != nil {
		return errors.New("UPDATE_HEADURL_ERR:" + err.Error())
	}

	if rownum != 1 {
		rownum, err = Orm.InsertOne(user)
		if err != nil || rownum != 1 {
			return errors.New("INSERT_USER_HEADURL_ERR:" + err.Error())
		}
	}

	return nil
}

func (user *User) RegName() error {
	user.LastMsgType = MSGTYPE_REG_NAME
	rownum, err := Orm.Id(user.UserID).Cols("name", "lastmsgtype").Update(user)
	if err != nil {
		return errors.New("UPDATE_NAME_ERR:" + err.Error())
	}

	if rownum != 1 {
		rownum, err = Orm.InsertOne(user)
		if err != nil || rownum != 1 {
			return errors.New("INSERT_USER_ERR:" + err.Error())
		}
	}

	//fmt.Sprintf(FMT_REG_NAME_SUC, user.Name),
	return nil
}
