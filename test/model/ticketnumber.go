package model

type LotteryType struct {
	ID     int    `gorm:"column:id;primary_key;auto_increment"`
	Name   string `gorm:"column:name;not null"`
	Namech string `gorm:"column:namech;not null"`
}

// TableName specifies the database table name
func (LotteryType) TableName() string {
	return "lottery_types"
}

type TicketNumber struct {
	ID               int    `gorm:"column:id;primary_key;auto_increment;NOT NULL;DEFAULT:0;"` // 編號
	LotteryTypeID    int    `gorm:"column:lottery_type_id;type:int(20);"`                     // 彩種編號
	CheckState       int    `gorm:"column:check_state;type:int(1);"`                          // 狀態 0為未確認,1為確認
	WinningNumber    string `gorm:"column:winning_number;type:varchar(255);"`                 // 該場次的球號
	AdditionalNumber string `gorm:"column:additional_number;type:varchar(255);"`              // 備註球號
	LotteryDay       string `gorm:"column:lottery_day;type:varchar(55);"`                     // 期號
	StartTime        string `gorm:"column:start_time;type:varchar(55);"`                      // 開始時間
	Special_Number   string `gorm:"column:special_number;type:varchar(55);"`                  // 特別號
	Original_Number  string `gorm:"column:original_number;type:varchar(55);"`                 // 特別號
}

// TableName 指定数据库表名
func (TicketNumber) TableName() string {
	return "Ticket_Numbers"
}
