from joblib import Parallel, delayed
import zipfile
from random import choice, uniform
import sys
import subprocess
from fake_headers import Headers
import undetected_chromedriver as uc
import warnings
from plugin_config import *
import selenium
from selenium import webdriver
from selenium.webdriver.common import by
from selenium.webdriver.common.keys import Keys
import indian_names
import random
from password_generator import PasswordGenerator
import requests
from selenium.webdriver.support.ui import Select
import time
from selenium.webdriver.remote.webelement import WebElement
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import pandas as pd
from selenium.webdriver.common.proxy import Proxy, ProxyType
from datetime import datetime, timedelta
import platform
from proxy_config import proxy
import traceback


def slow_type(element: WebElement, text: str, delay: float = 0.3):
    """Kirim teks ke elemen satu karakter pada satu waktu dengan penundaan."""
    for character in text:
        element.send_keys(character)
        time.sleep(uniform(.1, .3))


def phone_number_gen():
    api_key = "DUP-hxVS7EhQxS8SEDeKuPdXHfg-xyB8"
    phn_url = f"http://api.sms-man.com/stubs/handler_api.php?action=getNumber&api_key={api_key}&service=go&country=22"
    response = requests.get(phn_url)
    try:
        _id = response.text.split(":")[1]
        _number = response.text.split(":")[2]
    except IndexError:
        _id, _number = phone_number_gen()

    return _id.strip(), _number.strip()


def check_otp(id):
    api_key = "DUP-hxVS7EhQxS8SEDeKuPdXHfg-xyB8"
    status_url = f"http://api.sms-man.com/stubs/handler_api.php?action=getStatus&api_key={api_key}&id={id}"
    response = requests.get(status_url).text

    start_time = datetime.now()

    while(True):
        if datetime.now()-start_time > timedelta(minutes=2):
            return None

        if "STATUS_OK" in response:
            _otp = response.split(":")[1]
            print("OTP is", _otp)
            return _otp.strip()
        else:
            status_url = f"http://api.sms-man.com/stubs/handler_api.php?action=getStatus&api_key={api_key}&id={id}"
            response = requests.get(status_url).text


def enter_name_details(driver: webdriver.Chrome, f_name, l_name, username, password, wait):
    first_name = wait.until(EC.element_to_be_clickable(
        (By.XPATH, "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[1]/div/form/span/section/div/div/div[1]/div[1]/div[1]/div/div[1]/div/div[1]/input")))
    slow_type(first_name, f_name)

    last_name = wait.until(EC.element_to_be_clickable(
        (By.XPATH, "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[1]/div/form/span/section/div/div/div[1]/div[1]/div[2]/div/div[1]/div/div[1]/input")))
    slow_type(last_name, l_name)

    user_name = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                      "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[1]/div/form/span/section/div/div/div[2]/div[1]/div/div[1]/div/div[1]/input")
                                                      ))
    slow_type(user_name, username)

    pass1 = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                  "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[1]/div/form/span/section/div/div/div[3]/div[1]/div/div/div[1]/div/div[1]/div/div[1]/input")
                                                  ))
    slow_type(pass1, password)

    pass2 = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                  "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[1]/div/form/span/section/div/div/div[3]/div[1]/div/div/div[2]/div/div[1]/div/div[1]/input")
                                                  ))
    slow_type(pass2, password)

    next_button = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                        "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[2]/div/div[1]/div/div/button"))).click()


def Enter_Phone_Details(driver: webdriver.Chrome, phn_num, phn_id, wait):
    enter_phone_number = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                               "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[1]/div/form/span/section/div/div/div[2]/div/div[2]/div[1]/label/input")
                                                               ))
    enter_phone_number.clear()
    slow_type(enter_phone_number, phn_num)

    next_button_2 = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                          "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[2]/div/div[1]/div/div/button"))).click()

    time.sleep(3)

    try:
        check_phn_validity = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                                   "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[1]/div/form/span/section/div/div/div[2]/div/div[2]/div[2]/div")
                                                                   ))
        if "used too many times" in check_phn_validity.text:
            phn_id, phn_num = phone_number_gen()
            phn_id, phn_num = Enter_Phone_Details(
                driver=driver, phn_num=phn_num, phn_id=phn_id, wait=wait)

    except:
        pass

    return phn_num, phn_id


def Enter_Verify_Otp(driver: webdriver.Chrome, phn_id, phn_num, wait):

    get_code = check_otp(phn_id)

    while(get_code == None):
        back_button = wait.until(EC.element_to_be_clickable(
            (By.XPATH, "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[2]/div[1]/div/div/button"))).click()

        phn_id, phn_num = phone_number_gen()

        phn_num, phn_id = Enter_Phone_Details(
            driver, phn_num=phn_num, phn_id=phn_id, wait=wait)

        get_code = check_otp(phn_id)

    time.sleep(3)

    enter_otp = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                      "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[1]/div/form/span/section/div/div/div[2]/div/div[1]/div/div[1]/input")
                                                      ))
    slow_type(enter_otp, get_code)
    time.sleep(3)
    try:
        verify_button = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                              "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[2]/div[2]/div[1]/div/div/button"))).click()

    except:
        back_button = wait.until(EC.element_to_be_clickable(
            (By.XPATH, "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[2]/div[1]/div/div/button"))).click()

        phn_id, phn_num = phone_number_gen()

        phn_num, phn_id = Enter_Phone_Details(
            driver, phn_num=phn_num, phn_id=phn_id, wait=wait)

        phn_num, phn_id = Enter_Verify_Otp(
            driver, phn_num=phn_num, phn_id=phn_id, wait=wait)

    return phn_id, phn_num


def Enter_DOB(driver: webdriver.Chrome, day_bday, year_bday, month_bday, wait):
    bday_month = Select(wait.until(
        EC.element_to_be_clickable((By.ID, "month"))))
    bday_month.select_by_value(month_bday)
    time.sleep(3)

    bday_day = wait.until(EC.element_to_be_clickable((By.ID, "day")))
    slow_type(bday_day, day_bday)

    bday_year = wait.until(EC.element_to_be_clickable((By.ID, "year")))
    slow_type(bday_year, year_bday)

    gender = Select(wait.until(EC.element_to_be_clickable((By.ID, "gender"))))
    gender.select_by_value("1")
    time.sleep(3)

    next3_button = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                         "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[2]/div/div[1]/div/div/button"))).click()


def Generate_Details():
    pwo = PasswordGenerator()
    pwo.minlen = 9
    pwo.maxlen = 26

    password = pwo.generate()

    full_name = indian_names.get_full_name(gender='male')
    f_name = full_name.split()[0]
    l_name = full_name.split()[1]
    username = f_name+l_name+str(random.randint(10000, 6000000))

    phn_id, phn_num = phone_number_gen()

    month_bday = str(random.randint(1, 12))
    day_bday = str(random.randint(1, 30))
    year_bday = str(random.randint(1990, 1997))

    return password, username, phn_id, phn_num, month_bday, day_bday, year_bday, f_name, l_name


def prepare_env():
    OSNAME = platform.system()

    if OSNAME == 'Linux':
        OSNAME = 'lin'
        with subprocess.Popen(['google-chrome', '--version'], stdout=subprocess.PIPE) as proc:
            version = proc.stdout.read().decode('utf-8').replace('Google Chrome', '').strip()
    elif OSNAME == 'Darwin':
        OSNAME = 'mac'
        process = subprocess.Popen(
            ['/Applications/Google Chrome.app/Contents/MacOS/Google Chrome', '--version'], stdout=subprocess.PIPE)
        version = process.communicate()[0].decode(
            'UTF-8').replace('Google Chrome', '').strip()
    elif OSNAME == 'Windows':
        OSNAME = 'win'
        process = subprocess.Popen(
            ['reg', 'query', 'HKEY_CURRENT_USER\\Software\\Google\\Chrome\\BLBeacon', '/v', 'version'],
            stdout=subprocess.PIPE, stderr=subprocess.DEVNULL, stdin=subprocess.DEVNULL
        )
        version = process.communicate()[0].decode('UTF-8').strip().split()[-1]
    else:
        print('{} OS is not supported.'.format(OSNAME))
        sys.exit()

    major_version = version.split('.')[0]

    uc.TARGET_VERSION = major_version

    uc.install()

    return OSNAME


def prepare_proxy(proxy):

    OSNAME = prepare_env()

    proxy_split = proxy.split(":")
    PROXY_HOST = proxy_split[0]
    PROXY_PORT = proxy_split[1]
    PROXY_USER = proxy_split[2]
    PROXY_PASS = proxy_split[3]

    header = Headers(
        browser="chrome",
        os=OSNAME,
        headers=False
    ).generate()
    agent = header['User-Agent']

    options = webdriver.ChromeOptions()
    pluginfile = 'proxy_auth_plugin.zip'
    with zipfile.ZipFile(pluginfile, 'w') as zp:
        zp.writestr("manifest.json", manifest_json)
        zp.writestr("background.js", background_js %
                    (PROXY_HOST, PROXY_PORT, PROXY_USER, PROXY_PASS))
    options.add_extension(pluginfile)
    prefs = {"profile.default_content_setting_values.notifications": 2}
    options.add_experimental_option("prefs", prefs)
    options.add_experimental_option('prefs', {
        'credentials_enable_service': False,
        'profile': {
            'password_manager_enabled': False
        }
    })
    options.add_argument('--no-sandbox')
    options.add_argument('--disable-dev-shm-usage')
    options.add_argument("--disable-web-security")
    # viewport = ['2560,1440', '1920,1080']
    # options.add_argument(f"--window-size={choice(viewport)}")
    options.add_argument("--log-level=3")
    options.add_argument('--disable-blink-features=AutomationControlled')
    options.add_experimental_option(
        "excludeSwitches", ["enable-automation", "enable-logging"])
    options.add_experimental_option('useAutomationExtension', False)
    options.add_argument(f"user-agent={agent}")
    driver = webdriver.Chrome(options=options)

    return driver


def main(_proxy):
    try:
        # proxy = "103.167.32.223:45785:Selsharmashubham859:F4s1LzV"

        driver = prepare_proxy(proxy=_proxy)
        wait = WebDriverWait(driver, 40)

        print('\033[92m' + f'Proxy: {_proxy}' + '\033[0m')

        warnings.filterwarnings("ignore", category=DeprecationWarning)

        driver.get(
            "https://accounts.google.com/signup/v2/webcreateaccount?flowName=GlifWebSignIn&flowEntry=SignUp")

        password, username, phn_id, phn_num, month_bday, day_bday, year_bday, f_name, l_name = Generate_Details()
        enter_name_details(driver, f_name=f_name, l_name=l_name,
                           username=username, password=password, wait=wait)

        time.sleep(3)

        phn_num, phn_id = Enter_Phone_Details(
            driver, phn_num=phn_num, phn_id=phn_id, wait=wait)

        time.sleep(3)

        phn_id, phn_num = Enter_Verify_Otp(
            driver=driver, phn_id=phn_id, phn_num=phn_num, wait=wait)

        time.sleep(3)

        Enter_DOB(driver=driver, day_bday=day_bday,
                  month_bday=month_bday, year_bday=year_bday, wait=wait)

        time.sleep(3)

        skip_button = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                            "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[2]/div[2]/div[2]/div/div/button"))).click()

        time.sleep(3)

        i_agree_button = wait.until(EC.element_to_be_clickable((By.XPATH,
                                                               "/html/body/div[1]/div[1]/div[2]/div[1]/div[2]/div/div/div[2]/div/div[2]/div/div[1]/div/div/button"))).click()

        dob = day_bday+"-"+month_bday+"-"+year_bday
        data = [f_name, l_name, username +
                "@gmail.com", password, dob, _proxy]
        print(data)

        time.sleep(10)
        driver.close()
        return data
    except Exception as e:
        print(e)
        traceback.print_exc()
        driver.close()


if __name__ == "__main__":
    result = Parallel(n_jobs=3)(delayed(main)(_proxy)
                                for _proxy in proxy)
    print(result)

    df = pd.read_csv(
        r"C:\Users\NamaPCMU\PATH\address.csv")

    for data in result:
        if data:
            df.loc[len(df)] = data

    df.to_csv(
        r"C:\Users\NamaPCMU\PATH\address.csv", index=False)

    # main()

