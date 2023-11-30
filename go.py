import time
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from bs4 import BeautifulSoup
from selenium.webdriver.chrome.service import Service
# Constants
port = 8080
max_attempts = 5
retry_interval = 11

# Chrome options
chrome_options = Options()
# Uncomment the line below to run headless
chrome_options.add_argument('--headless')
chrome_options.add_argument('--disable-gpu') #關閉GPU 避免某些系統或是網頁出錯

# Initialize WebDriver
driver = webdriver.Chrome()
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
        ip_addresses = [td.text.strip() for td in td_elements if td.text.strip()]

        return ip_addresses
    else:
        print(f"Failed to retrieve content. Status code: {response.status_code}")
        return []

if __name__ == "__main__":
    url = "https://free-proxy-list.net/"
    proxy_ips = scrape_proxy_ips(url)

    print("Proxy IPs:",proxy_ips)

# Get the game list page
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