import time
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from bs4 import BeautifulSoup
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
# Constants
port = 8080
max_attempts = 5
retry_interval = 11


import requests
from bs4 import BeautifulSoup

def scrape_proxy_ips(url):
    # 發送 GET 請求
    response = requests.get(url)

    if response.status_code == 200:
        # 解析 HTML
        soup = BeautifulSoup(response.text, 'html.parser')

        # 找到包含 IP 的 <td> 元素
        td_elements = soup.find_all('td')

        # 提取 IP 地址
        ip_addresses = [td.text.strip() for td in td_elements if td.text.strip().count('.') == 3]

        return ip_addresses
    else:
        print(f"Failed to retrieve content. Status code: {response.status_code}")
        return []

if __name__ == "__main__":
    url = "https://free-proxy-list.net/"
    proxy_ips = scrape_proxy_ips(url)
    
    print("Proxy IPs:",proxy_ips,proxy_ips[0])

# Get the game list page
# Chrome options
chrome_options = Options()
# set header 
headers = {
        "Accept": "application/json, text/javascript, */*; q=0.01",
        "Accept-Encoding": "gzip, deflate, br",
        "Accept-Language": "zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7",
        "Cache-Control": "no-cache",
        "Content-Length": "33",
        "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
        "Cookie": "visid_incap_1930594=luSvdSdLSdKH9+tglZ+uWp8BW2UAAAAAQUIPAAAAAAD1kv/T4c8KUS50vYwhwNyC; _gid=GA1.2.1888423171.1701050902; nlbi_1930594=7yWSXwv3uz0xSDNfT4OJ5QAAAADYK1LAJ24MJ56GYB03Xmlv; Hm_lvt_7a95cc9501d6fe6f73efcfe5575b4eef=1701229308,1701312313,1701396255,1701401096; incap_ses_936_1930594=xeCED5SX6jpvaxVRD1j9DH54aWUAAAAAIbj2RCoYi+xlvbHp+7uddg==; _ga_2Z09KKD993=GS1.2.1701410944.31.0.1701410944.0.0.0; Hm_lpvt_7a95cc9501d6fe6f73efcfe5575b4eef=1701410944; ci_session=afdf48dd2943b9651a88ed621264287f6d802cbd; _ga_7PNCWWB3KR=GS1.1.1701410944.33.0.1701412411.0.0.0; _ga=GA1.2.763795639.1700462919; _gat=1",
        "Origin": "https://www.lkag3.com",
        "Pragma": "no-cache",
        "Referer": "https://www.lkag3.com/Issue/history?lottername=GSUS",
        "Sec-Ch-Ua": "\"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"",
        "Sec-Ch-Ua-Mobile": "?0",
        "Sec-Ch-Ua-Platform": "\"Windows\"",
        "Sec-Fetch-Dest": "empty",
        "Sec-Fetch-Mode": "cors",
        "Sec-Fetch-Site": "same-origin",
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
        "X-Requested-With": "XMLHttpRequest"
    }
for key, value in headers.items():
    chrome_options.add_argument(f'--header={key}: {value}')
# Uncomment the line below to run headless
# chrome_options.add_argument('--headless')
chrome_options.add_argument('--disable-gpu') #關閉GPU 避免某些系統或是網頁出錯
chrome_options.add_argument(f'--proxy-server=https://34.77.56.122')
# Initialize WebDriver
driver = webdriver.Chrome(
    service=Service(ChromeDriverManager().install()),
    options=chrome_options)
driver.get("https://www.lkag3.com/index/lotterylist")

# Get the page source
source = driver.page_source

# Find lotternames
element_tag = "href"
element_title = "a"
contains = "lottername="
lotternames = []

soup = BeautifulSoup(source, 'html.parser')

for a_tag in soup.find_all(element_title, href=True):
    if element_tag in a_tag.attrs:
        if contains in a_tag[element_tag]:
            lottername = a_tag[element_tag].replace("https://www.lkag3.com/Issue/history?lottername=", "")
            lotternames.append(lottername)

print("Lotternames:")
score_list = []
for lottername in lotternames:
    print(lottername)
    score_list.append(lottername)

for i, lottername in enumerate(score_list[38:]):
    for attempt in range(max_attempts):
        driver.get("https://www.lkag3.com/Issue/history?lottername=" + lottername)

        print(lottername)
        time.sleep(2)

        page_source = driver.page_source
        soup = BeautifulSoup(page_source, 'html.parser')

        td_elements = soup.select("td.ball")
        if not td_elements:
            print("No matching <td> element found")
            time.sleep(10)
            continue

        spans = []
        for span in td_elements[0].select("div.b1 span, div.b2 span, div.b3 span, div.b4 span, td.v1 b1, div.gbs_bg span, tbody"):
            spans.append(span.text)

        print("Content of <td>: ", spans)

        if "429 Too Many Requests" in page_source:
            print("Received 429 Too Many Requests. Waiting for a while and retrying...")
            time.sleep(15)
            continue

        # Break out of the loop if successful
        break

time.sleep(55)
driver.quit()