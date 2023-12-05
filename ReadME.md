# 使用GO撰寫網路爬蟲爬取各站點比分
 運行資料庫容器:
 docker-compose up -d 
 **藉由mod取得各項依賴選項**

 `cd test
 test\:> go mod init test & go mod tidy
 // 運行程式碼 
 test\:> go run main.go `


## 功能檢索與比對

站源發出來的時間有延遲 地域性等關係 收到的req最理想狀況是同步站源時間

客戶端要時常同步 一般網頁都會做出阻擋(單個地址10秒內可以傳出的有效REQ不超過三次)
功能在連續請求時候被
偵測回傳429:  
目前: 當前索引返回429之後保留重新進行訪問
未來: 代理的方式請求確保數據在最小的時間取得
處理方案: 使用proxy 取可用IP資源列表 在發送請求時輪循此列表取其一 並仿照頭部 若能請求成功將存入IP池
取得的資訊問題: 目前由網站取出值,但不能確定網站是根據總站有多少的延遲,所以根據彩卷名稱找出對應的其他資源


## note
彩果資訊DB完成 
今日:新增一個欄位 checkstate 將比分場次與不同站源比對，
位於 lottery_types  正常時為空值
將其他彩果對資料比對成功後 確認後欄位值為1
目前進度為16/160 
問題:因為站源的排版方向不一致 需要個別解析 較耗費時間
請求目前 會是T=160*8s 訪問每個頁面不返回429的最低時間

預計:
- 新增PROXY功能 T/代理數量 約10分鐘取得同步
- 將爬蟲設為背景 規劃同步流程


-- Adminer 4.8.1 MySQL 5.7.44 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `lottery_types`;
CREATE TABLE `lottery_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `namech` longtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


DROP TABLE IF EXISTS `ticket_numbers`;
CREATE TABLE `ticket_numbers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `lottery_type_id` int(11) NOT NULL,
  `checkstate` bit(1) DEFAULT NULL,
  `winning_number` varchar(255) DEFAULT NULL,
  `additional_number` varchar(255) DEFAULT NULL,
  `lottery day` varchar(55) DEFAULT NULL,
  `start_time` varchar(55) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 2023-12-05 11:00:18