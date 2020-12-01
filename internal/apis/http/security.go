package http

// auth 鉴权
type auth struct {
	onlineToken map[string]string // onlineToken 用户在线 Token
}

type (
	loginToken  string // loginToken 登录永久 Token
	onlineToken string // onlineToken 用户在线 Token
)

// NewAuth 新建一个鉴权
func NewAuth() *auth { return &auth{} }

// generateLoginToken 使用 u：用户名(User) 生成 loginToken
func generateLoginToken(u string) loginToken {
	return ""
}

// generateOnlineToken 使用 u：用户名(User) t：登录永久 Token(loginToken) 生成 onlineToken
func generateOnlineToken(u string, t loginToken) onlineToken {
	return ""
}

// verifyOnlineToken 使用 u：用户名(User) t：用户在线 Token(onlineToken) 鉴权
// bool 值返回鉴权结果 T：允许 F：禁止
func (a *auth) verifyOnlineToken(u string, t onlineToken) bool {
	return true
}
