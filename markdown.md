# 使用GO撰寫網路爬蟲爬取各站點比分
 
 // 藉由mod取得各項依賴選項
 go mod init getscore
 // 安裝依賴
 go mod tidy
 運行程式碼 go run main.go 
# python 方式
 // py or py3 go.py 

## 功能檢索與比對
功能在連續請求時候被偵測回傳429  處理方案: 目前想到先使用freeproxy 爬取一個免費列表 在發送請求時輪循此列表帶入資訊之後發送請求