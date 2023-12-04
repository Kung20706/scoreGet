package main 

func Dbconn(){// 連線到 MySQL 資料庫
dsn := "db_user:db_password@tcp(127.0.0.1:3306)/db_name?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
	log.Fatal(err)
}
defer func() {
	// 正確的方式是使用 sql.DB 的 Close 方法
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Close()
}()
}