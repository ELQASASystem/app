package database

import "fmt"

// ReadAccountsList 使用 u：用户名(ID) 查询 AccountsListTab 表。
// 校验密码
func (a *account) ReadAccountsList(u string) (data *AccountsListTab, err error) {

	sq := fmt.Sprintf(`SELECT accounts_list.* FROM accounts_list WHERE accounts_list.id = "%v"`, u)

	row, err := a.conn.Query(sq)
	if err != nil {
		return
	}
	defer row.Close()

	if !row.Next() {
		return
	}
	data = new(AccountsListTab)
	err = row.Scan(&data.ID, &data.Password, &data.Class, &data.LoginToken)
	if err != nil {
		return
	}

	return
}

// UpdateLoginToken 更新 AccountsListTab 表中 u：用户名 的 loginToken 字段为 t：loginToken。
// 更新 loginToken
func (a *account) UpdateLoginToken(t, u string) (err error) {

	l, err := a.conn.Prepare(`UPDATE accounts_list SET accounts_list.login_token = ? WHERE accounts_list.id = ?`)
	if err != nil {
		return
	}
	defer l.Close()

	if _, err = l.Exec(t, u); err != nil {
		return
	}

	return
}
