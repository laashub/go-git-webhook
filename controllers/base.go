package controllers

import (
	"github.com/astaxie/beego"
	"go-git-webhook/models"
	"go-git-webhook/conf"
)

type BaseController struct {
	beego.Controller
	Member *models.Member
	Scheme string
}

func (c *BaseController) Prepare (){
	c.Data["SiteName"] = "Git WebHook"
	c.Data["Member"] = models.Member{}

	if member,ok := c.GetSession(conf.LoginSessionName).(models.Member); ok && member.MemberId > 0{
		c.Member = &member
		c.Data["Member"] = c.Member
	}
	scheme := "http"

	if c.Ctx.Request.TLS != nil {
		scheme = "https"
	}
	c.Scheme = scheme
}

//获取或设置当前登录用户信息,如果 MemberId 小于 0 则标识删除 Session
func (c *BaseController) SetMember(member models.Member) {

	if member.MemberId <= 0 {
		c.DelSession(conf.LoginSessionName)
		c.DelSession("uid")
	} else {
		c.SetSession(conf.LoginSessionName, member)
		c.SetSession("uid", member.MemberId)
	}
}

//响应 json 结果
func (c *BaseController) JsonResult(errCode int,errMsg string,data ...interface{}){
	json := make(map[string]interface{},3)

	json["errcode"] = errCode
	json["message"] = errMsg

	if len(data) > 0 && data[0] != nil{
		json["data"] = data[0]
	}

	c.Data["json"] = json
	c.ServeJSON(true)
	c.StopRun()
}

func (c *BaseController) UrlFor (endpoint string, values ...interface{}) string {

	return c.BaseUrl() + beego.URLFor(endpoint,values...)
}

func (c *BaseController) BaseUrl() string {
	scheme := "http://"

	if c.Ctx.Request.TLS != nil {
		scheme = "https://"
	}
	return scheme + c.Ctx.Request.Host
}