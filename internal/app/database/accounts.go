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
	err = row.Scan(&data.ID, &data.Password, &data.Class)
	if err != nil {
		return
	}

	return
}
