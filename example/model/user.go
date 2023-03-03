package model

type User struct {
	Id          uint   `gorm:"column:id;autoIncrement;primaryKey"`
	Name        string `gorm:"column:name;type:varchar(255)"`
	PhoneNumber string `gorm:"column:phone_number;type:varchar(255)"`
	Address     string `gorm:"column:address;type:varchar(255)"`
	IdNo        string `gorm:"column:id_no;type:varchar(255)"`
	UpdateAt    int64  `gorm:"column:update_at;type:bigint"`
}

func (u *User) CryptoColumns() []string {
	return []string{
		userColumn.IdNo,
		userColumn.Address,
		userColumn.PhoneNumber,
	}
}

func (u *User) TableName() string {
	return "user"
}

func GetUserColumn() UserColumn {
	return userColumn
}

var userColumn = UserColumn{
	Id:          "id",
	Name:        "name",
	PhoneNumber: "phone_number",
	Address:     "address",
	IdNo:        "id_no",
	UpdateAt:    "update_at",
}

type UserColumn struct {
	Id          string
	Name        string
	PhoneNumber string
	Address     string
	IdNo        string
	UpdateAt    string
}
