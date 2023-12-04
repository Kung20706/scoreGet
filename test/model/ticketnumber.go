package models

import (
)

type LotteryType struct {
	ID   int    `gorm:"column:id;primary_key;auto_increment"`
	Name string `gorm:"column:name;not null"`
}

// TableName specifies the database table name
func (LotteryType) TableName() string {
	return "lottery_types"
}