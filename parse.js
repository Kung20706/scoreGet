const formData = new FormData();
	formData.append('lotterycode', 'CQSSC');
	formData.append('lotteryname', 'CQSSC');

var headers = new Headers();

	// 添加需要的请求头信息// 添加请求头信息
	headers.append('Accept', 'application/json, text/javascript, */*; q=0.01');
	headers.append('Accept-Encoding', 'gzip, deflate, br');
	headers.append('Accept-Language', 'zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7');
	headers.append('Cache-Control', 'no-cache');
	headers.append('Content-Length', '35');
	headers.append('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
	headers.append('Origin', 'https://www.lkag3.com');
	headers.append('Pragma', 'no-cache');
	headers.append('User-Agent', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36');
	headers.append('X-Requested-With', 'XMLHttpRequest');
	// 添加其他请求头，根据需要添加更多
	
	// 创建一个包含请求头的 options 对象
	var requestOptions = {
	  method: 'POST',  // 设置请求的方法
	  headers: headers, // 将 Headers 对象传递给 headers 属性
	  body: 'lotterycode=CQSSC&lotteryname=CQSSC' // 请求的 body，可以是字符串、FormData 等
	};
	// 使用fetch发送POST请求
	fetch('https://www.lkag3.com/Issue/ajax_history', requestOptions	)
	  .then(response => response.json())
	  .then(data => {
		// 处理返回的数据
		console.log(data);
	  })
	  .catch(error => {
		console.error('Error:', error);
	  });
    
