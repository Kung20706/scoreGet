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


DROP TABLE IF EXISTS `Ticket_Numbers`;
CREATE TABLE `Ticket_Numbers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `lottery_type_id` int(11) NOT NULL,
  `check_state` int(1) DEFAULT NULL,
  `winning_number` varchar(255) DEFAULT NULL,
  `additional_number` varchar(255) DEFAULT NULL,
  `lottery_day` varchar(55) DEFAULT NULL,
  `start_time` varchar(55) DEFAULT NULL,
  `lottery day` varchar(55) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 2023-12-11 10:49:42

體育
* 足球 Football 
* 藍球 Basketball
* 網球 Tennis
* 桌球  Table Tennis
* 曲棍球 Hockey
* 電競 Esports  
* 手球 Handball
* 排球 Volleyball
* 棒球 Baseball
* 橄欖球 American Football
* 綜合格鬥 MMA
* 賽車運動 Motorsport
* 撞球 snooker
* 室內足球 Futsal
* 迷你足球 Minifootball
* 羽毛球 Badminton
* 澳式足球 Aussie Rules
* 沙灘排球 Beach Volleyball
* 水球 Waterpolo
* 自行車 Cycling
* 軟式曲棍球 Floorball
* 俄式冰球 Bandy
高頻彩種
  時時彩
  分分彩
低頻彩種
  六合彩 

六和/台灣彩卷
  台灣四星彩
  台灣三星彩
  台灣威力彩
  台灣大樂透
  金彩539
  香港六合彩
美國天天樂
  加州天天樂
  密西根天天樂
  俄克拉河馬天天樂
  佛羅里達天天樂
菲律賓樂透
  菲律賓樂透49
  菲律賓樂透45
  菲律賓樂透55
  菲律賓樂透42
分分彩
  台灣快開
基諾彩
  台灣賓果
PC蛋蛋
  台灣28
越南大樂透
  河內(北部)
  富安(中部)
  順化(中部)
  同塔(南部)
  金甌(南部)
  胡志明(南部)
 

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


DROP TABLE IF EXISTS `Ticket_Numbers`;
CREATE TABLE `Ticket_Numbers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `lottery_type_id` int(11) NOT NULL,
  `check_state` int(1) DEFAULT NULL,
  `winning_number` varchar(255) DEFAULT NULL,
  `special_number` varchar(55) DEFAULT NULL,
  `additional_number` varchar(255) DEFAULT NULL,
  `lottery_day` varchar(55) DEFAULT NULL,
  `start_time` varchar(55) DEFAULT NULL,
  `original_number` varchar(55) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 2023-12-15 05:22:48

 
