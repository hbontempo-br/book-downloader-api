package models

import "fmt"

type BookStatusModel struct {
	ID         uint   `gorm:"primary_key"`
	Enumerator string `gorm:"column:enumerator"`
}

func (BookStatusModel) TableName() string {
	return "BookStatus"
}

func (bsm BookStatusModel) String() string {
	return fmt.Sprintf("{ID:%v, Enumerator:%v}", bsm.ID, bsm.Enumerator)
}
