package model

type TicketNumber struct {
	ID               int    `gorm:"column:id;primary_key;type:int(20);NOT NULL;DEFAULT:0"`
	LotteryTypeID    int    `gorm:"column:lottery_type_id;type:int(20);"`
	CheckState       int    `gorm:"column:check_state;type:int(1);"`
	WinningNumber    string `gorm:"column:winning_number;type:varchar(255);"`
	AdditionalNumber string `gorm:"column:additional_number;type:varchar(255);"`
	LotteryDay       string `gorm:"column:lottery_day;type:varchar(55);"` // 数据库中是 date 类型，这里使用 string 类型
	StartTime        string `gorm:"column:start_time;type:varchar(55);"`
}

// TableName 指定数据库表名
func (TicketNumber) TableName() string {
	return "Ticket_Numbers"
}
