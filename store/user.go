package store

import "database/sql"

type User struct {
	UID int
	VKID int
	Node int
	UserName string

	SelectDate int
	SelectCollection int
	SelectType int
	SelectStoreys int
	SelectProject int
	Data string
}

func createUser(vkid int) {
	_, err := DB.Exec("insert into users (vkid) values ($1)", vkid)
	CheckErr(err)
}

func GetUser(vkid int) *User {
	u := &User{}
	row := DB.QueryRow("select uid, vkid, node, UserName, SelectDate, SelectCollection, SelectType, SelectStoreys, SelectRroject, Data from users where vkid=$1", vkid)
	err := row.Scan(&u.UID, &u.VKID, &u.Node, &u.UserName, &u.SelectDate, &u.SelectCollection, &u.SelectStoreys, &u.SelectProject, &u.Data)
	if err == sql.ErrNoRows {
		createUser(vkid)
		return GetUser(vkid)
	}
	CheckErr(err)
	return u
}

func (u *User) set(param string, argument interface{}) {
	_, err := DB.Exec("update users set "+param+"=$1 where uid=$2", argument, u.UID)
	CheckErr(err)
}

func (u *User) SetNode(node int) {
	u.set("node", node)
}

func (u *User) SetSelectDate(date int) {
	u.SelectDate = date
	u.set("SelectDate", date)
}


func (u *User) SetSelectCollection(c int) {
	u.SelectCollection = c
	u.set("SelectCollection", c)
}


func (u *User) SetSelectType(t int) {
	u.SelectType = t
	u.set("SelectType", t)
}


func (u *User) SetSelectStoreys(s int) {
	u.SelectStoreys = s
	u.set("SelectStoreys", s)
}


func (u *User) SetSelectProject(p int) {
	u.SelectProject = p
	u.set("SelectProject", p)
}


func (u *User) SetUserName(n string) {
	u.UserName = n
	u.set("Username", n)
}

func (u *User) SetData(d string) {
	u.Data = d
	u.set("data", d)
}

func (u *User) Delete() {
	_, err := DB.Exec("delete from users where uid=$1", u.UID)
	CheckErr(err)
}
