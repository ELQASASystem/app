package http

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/ELQASASystem/server/internal/app"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"

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
func (a *auth) generateLoginToken(u string, c *context.Context) {

	var (
		original = bytes.NewBuffer(nil)
		ti, _    = time.Now().MarshalText()
		salt     = make([]byte, 16)
	)

	_, _ = rand.Read(salt)
	original.WriteString(u)
	original.Write(ti)
	original.Write(salt)

	token := sha256.Sum256(original.Bytes())
	t := loginToken(fmt.Sprintf("%x", token))

	err := app.AC.DB.Account().UpdateLoginToken(string(t), u)
	if err != nil {
		log.Error().Err(err).Msg("更新数据库 LoginToken 失败")
		return
	}

	c.SetCookie(&http.Cookie{
		Name: "loginToken", Value: string(t),
		Expires: time.Now().AddDate(0, 1, 0), Secure: true,
	})
	return
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
