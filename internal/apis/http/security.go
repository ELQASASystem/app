package http

import (
	"github.com/kataras/iris/v12/context"
)

// auth 鉴权
type auth struct {
	onlineTokenList map[string]onlineToken // onlineTokenList 用户在线 Token 列表
}

type (
	loginToken  string // loginToken 登录永久 Token
	onlineToken string // onlineToken 用户在线 Token
)

// NewAuth 新建一个鉴权
func NewAuth() *auth { return &auth{} }

// generateLoginToken 使用 u：用户名 生成 loginToken
func (a *auth) generateLoginToken(u string) loginToken {
	return ""
}

// generateOnlineToken 使用 u：用户名 t：loginToken 生成 onlineToken
func (a *auth) generateOnlineToken(u string, t loginToken) onlineToken {
	return ""
}

// verifyOnlineToken onlineToken 鉴权
// bool 值返回鉴权结果 T：允许 F：禁止
func (a *auth) verifyOnlineToken(c *context.Context) (r bool) {

	if c.GetCookie("onlineToken") == string(a.onlineTokenList[c.GetCookie("user")]) {
		return true
	}

	c.StatusCode(401)
	return
}
