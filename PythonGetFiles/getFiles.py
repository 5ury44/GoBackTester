import argparse
import datetime
import os
import time

from selenium.webdriver.support import expected_conditions as EC
import requests
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.wait import WebDriverWait


def month_year_range(start_month, start_year, end_month, end_year):
    result = []
    start_date = datetime.date(start_year, start_month, 1)
    end_date = datetime.date(end_year, end_month, 1)
    current_date = start_date
    while current_date <= end_date:
        result.append(current_date)
        current_month = current_date.month + 1
        current_year = current_date.year
        if current_month > 12:
            current_month = 1
            current_year += 1
        current_date = datetime.date(current_year, current_month, 1)
    return result


parser = argparse.ArgumentParser(description="Add dates")
parser.add_argument("startDate", help="start date")
parser.add_argument("endDate", help="end date")
parser.add_argument("currency", help="curr")
# startDate = str(args.startDate).split("-")
# endDate = str(args.endDate).split("-")
# currency = str(args.currency)
startDate = str("10-2022").split("-")
endDate = str("11-2022").split("-")
currency = "eurusd"
res = month_year_range(int(startDate[0]), int(startDate[1]), int(endDate[0]), int(endDate[1]))

url = "https://www.truefx.com/truefx-login/"
data = {
    "log": "5ury44",
    "pwd": "suryaaisthegreatestguyik",
    "rememberme": "forever",
    "redirect_to": "/streaming-market-data-truefx/",
    "mepr_process_login_form": "true",
    "mepr_is_login_page": "true"
}
with requests.Session() as session:
    # Send the POST request
    responseErr = session.post(url, data=data)

    # Check the response status code to see if the request was successful
    if responseErr.status_code == 200:
        # The request was successful
        print("Login successful")
    else:
        # There was an error with the request
        print("Login failed")

    months = []
    year_adj = ''
    name_to_number = {}
    test_val = ""
    for i in range(1, 13):
        months.append(datetime.date(2008, i, 1).strftime('%B'))
    for i in range(1, 13):
        if len(str(i)) == 1:
            year_adj = "0" + str(i)
        else:
            year_adj = i
        name_to_number[str(datetime.date(2008, i, 1).strftime('%B'))] = year_adj

    if 'files' not in os.listdir():
        os.mkdir('files')

    response = session.get('https://www.truefx.com/truefx-historical-downloads/')
    html = BeautifulSoup(response.content, 'html.parser')

    idCats = []
    for month in html.find_all('a'):
        if months.__contains__(str(month.get('title'))):
            idCats.append(str(month.get('data-idcat')))
    response.close()

    driver = webdriver.Chrome()
    driver.get("https://www.truefx.com/truefx-login/")
    username_field = driver.find_element(By.NAME, "log")
    password_field = driver.find_element(By.NAME, "pwd")
    login_form = driver.find_element(By.ID, "wp-submit")
    username_field.send_keys("5ury44")
    password_field.send_keys("suryaaisthegreatestguyik")
    login_form.submit()

    for idCat in idCats:
        found = False
        driver.execute_script("window.open('');")
        driver.switch_to.window(driver.window_handles[-1])
        driver.get('https://www.truefx.com/truefx-historical-downloads/#93-' + str(idCat) + '-top')
        time.sleep(2)
        html = BeautifulSoup(driver.page_source, 'html.parser')

        for file in html.find_all('a'):

            if str(file.get('href')).__contains__(currency):
                for dt in res:
                    mstr = str(dt.month)
                    if len(mstr) == 1:
                        mstr = "0" + str(dt.month)
                    if str(file.get('href'))[45:].__contains__(mstr) and str(file.get('href'))[45:].__contains__(
                            str(dt.year)):
                        print(str(file.get('href')))
                        found = True
                        break
        if found:
            break
        driver.close()
        driver.switch_to.window(driver.window_handles[-1])
    driver.close()
