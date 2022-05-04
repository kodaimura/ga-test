package entity


type User struct {
	UId int `db:"UID" json:"uid"`
	Username string `db:"USERNAME" json:"username"`
	Password string `db:"PASSWORD" json:"password"`
	CreateAt string `db:"CREATE_AT" json:"createat"`
	UpdateAt string `db:"UPDATE_AT" json:"updateat"`
}
