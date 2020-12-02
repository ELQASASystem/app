package http

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ELQASASystem/server/internal/app"

	"github.com/kataras/iris/v12/context"
	"github.com/rs/zerolog/log"
)

type (
	// auth 鉴权
	auth struct {
		app             *app.App               // app App
		onlineTokenList map[string]onlineToken // onlineTokenList 用户在线 Token 列表
	}

	loginToken  string // loginToken 登录永久 Token
	onlineToken string // onlineToken 用户在线 Token
)

var Banned = errors.New("禁止登录")

// NewAuth 新建一个鉴权
func NewAuth() *auth { return &auth{app.AC, map[string]onlineToken{}} }

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

	err := a.app.DB.Account().UpdateLoginToken(string(t), u)
	if err != nil {
		log.Error().Err(err).Msg("更新数据库 LoginToken 失败")
		return
	}

	c.SetCookie(&http.Cookie{
		Name: "loginToken", Value: string(t),
		Path: "/", Expires: time.Now().AddDate(0, 1, 0), Secure: true,
	})

	_ = a.generateOnlineToken(u, t, c, false)
}

// generateOnlineToken 使用 u：用户名 t：loginToken 生成 onlineToken
func (a *auth) generateOnlineToken(u string, lt loginToken, c *context.Context, check bool) (err error) {

	if check {
		if !a.checkLoginToken(u, lt) {
			c.StatusCode(401)
			return Banned
		}
	}

	var (
		original = bytes.NewBuffer(nil)
		salt     = make([]byte, 8)
	)

	original.WriteString(string(lt))
	original.WriteString(u)
	original.Write(salt)

	token := sha1.Sum(original.Bytes())
	t := onlineToken(fmt.Sprintf("%x", token))

	a.onlineTokenList[u] = t
	c.SetCookie(&http.Cookie{
		Name: "onlineToken", Value: string(t),
		Path: "/", Expires: time.Now().Add(time.Hour), Secure: true,
	})
	return
}

// checkLoginToken 检查 loginToken
func (a *auth) checkLoginToken(u string, t loginToken) (right bool) {

	res, err := a.app.DB.Account().ReadAccountsList(u)
	if err != nil {
		log.Error().Err(err).Msg("读取用户信息出错")
		return
	}

	if string(t) == res.LoginToken {
		return true
	}

	return
}

// verifyOnlineToken onlineToken 鉴权
// bool 值返回鉴权结果 T：允许 F：禁止
func (a *auth) verifyOnlineToken(c *context.Context) (r bool) {

	u := c.GetCookie("user")

	if c.GetCookie("onlineToken") == string(a.onlineTokenList[u]) {
		return true
	}

	if a.generateOnlineToken(u, loginToken(c.GetCookie("loginToken")), c, true) == nil {
		return true
	}

	c.StatusCode(401)
	return
}
