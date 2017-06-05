# coding: utf-8
import requests, os, time, sys, json, mysql.connector, logging, ConfigParser
sys.path.append("/usr/local/lib/python2.7/site-packages")
from apscheduler.schedulers.blocking import BlockingScheduler
from bs4 import BeautifulSoup
from PIL import Image
from io import BytesIO

cf = ConfigParser.ConfigParser()
cf.read("movie.conf")
config = {
    'user': cf.get("db", "db_user"),
    'password': cf.get("db", "db_pass"),
    'host': cf.get("db", "db_host"),
    'port': cf.get("db", "db_port"),
    'database': 'movie',
    'raise_on_warnings': True,
    'charset': 'utf8'
}

thisWeekUrl = 'https://tw.movies.yahoo.com/movie_thisweek.html'
intheaters = 'https://tw.movies.yahoo.com/movie_intheaters.html?p=%d'
urlObj = [[thisWeekUrl, False], [intheaters, True]]

titles = ['id', '中片名', '英片名', '上映日期', '類型', '片長', '導演', '演員', '發行公司', '官方網站', '評分', '簡介', 'img']
dbColName = ['id', 'cname', 'ename', 'releaseTime', 'type', 'duration', 'director', 'actor', 'company', 'website', 'score', 'intro', 'imgPath']
picPath = 'static/'
if not os.path.exists(picPath):
    os.makedirs(picPath)
    
logging.basicConfig(filename='movie_crawler.log', level=logging.DEBUG, format='%(asctime)s %(message)s', datefmt='%m/%d/%Y %I:%M:%S %p')
# keyMap = {
#     '上映日期' : 'release',
#     '類型' : 'type',
#     '片長' : 'duration',
#     '導演' : 'director',
#     '演員' : 'actor',
#     '發行公司' : 'company',
#     '官方網站' : 'website'
# }
def findIndexOfList(somelist, x):
    return somelist.index(x) if x in somelist else -1
def toUTF8(str):
    if isinstance(str, basestring):
        return str.encode("utf-8")
    else:
        return str
def getData(url, isMore):
    aData = []
    i = 1
    thisUrl = url
    while True:
        if isMore:
            thisUrl = url%i
        res = requests.get(thisUrl)
        soup = BeautifulSoup(res.text, "html.parser")
        moviesContainer = soup.select('.row-container')
        if len(moviesContainer) == 0:
            logging.info("No more information.")
            break
        for index, movie in enumerate(moviesContainer):
            movie_href = movie.select('a')[0]['href'].split('*')[1]
            if movie_href is not None and len(movie_href) > 0:
                while True:
                    resMovie = requests.get(movie_href)
                    mid = movie_href.split('/id=')[1]
                    info = findMovieInfo(resMovie.text, mid)
                    if info != False:
                        break
                    logging.info('fail' + mid)
                aData.append(tuple(info))
        if isMore == False:
            logging.info("No more page.")
            break
        i += 1
    return aData

def findMovieInfo(html, mid):
    dataArr = ["" for x in range(len(titles))]
    soupMovie = BeautifulSoup(html, "html.parser")
    # Information container
    if len(soupMovie.select('#ymvmvf')) == 0:
        return False
    infoMovie = soupMovie.select('#ymvmvf')[0]
    # Get information
    rawDat = infoMovie.select('p span.tit')
    cname = soupMovie.select('#ymvmvf .text.bulletin h4')[0].text.encode('utf8')
    ename = soupMovie.select('#ymvmvf .text.bulletin h5')[0].text.encode('utf8')
    imghref = soupMovie.select('#ymvmvf .bd-container a')[0]['href'].split('*')[1]
    imgPath = getPic(imghref, mid)
    pos = findIndexOfList(titles, 'id')
    if pos > -1:
        dataArr[pos] = mid
    pos = findIndexOfList(titles, '中片名')
    if pos > -1:
        dataArr[pos] = cname
    pos = findIndexOfList(titles, '英片名')
    if pos > -1:
        dataArr[pos] = ename
        
    for tit in rawDat:
        key = tit.text.encode('utf8').replace('：', '').replace('　', '')
        data = tit.parent.select('.dta')[0].text.encode('utf8').replace('\r', '').replace('\n', '')
        pos = findIndexOfList(titles, key)
        if pos > -1:
            dataArr[pos] = data
    rate = soupMovie.select('#ymvis .sum em')[0].text.encode('utf8')
    intro = soupMovie.select('#ymvs .text.full p')
    if len(intro) > 0:
        intro = intro[0]
        intro.find('p').decompose()
        intro = intro.encode('utf8').replace('<p>', '').replace('</p>', '')
    else:
        intro = ''
    pos = findIndexOfList(titles, '評分')
    if pos > -1:
        dataArr[pos] = rate
    pos = findIndexOfList(titles, '簡介')
    if pos > -1:
        dataArr[pos] = intro
    
    pos = findIndexOfList(titles, 'img')
    if pos > -1:
        dataArr[pos] = imgPath
    return dataArr

def getPic(url, fileName):
    path = picPath + fileName + '.jpg'
    import os.path
    if os.path.exists(path) == False:
        r = requests.get(url)
        i = Image.open(BytesIO(r.content))
        i.save(path)
    return path

def crawlerMovie():
    logging.info('Crawler Start!!!')
    collectData = []
    for item in urlObj:
        someData = getData(item[0], item[1])
        collectData.extend(someData)
    try:
        cnx = mysql.connector.connect(**config)
        cursor = cnx.cursor()
        colStr = ",".join(str(x) for x in dbColName)
        sArr = ["%s" for x in range(len(titles))]
        sStr = ",".join(str(x) for x in sArr)
        stmt = "INSERT INTO movieList (" + colStr + ") VALUES (" + sStr + ") ON DUPLICATE KEY UPDATE score = VALUES(score)"

        cursor.executemany(stmt, collectData)

        cnx.commit()
        cnx.close()
    except mysql.connector.Error as err:
        logging.error("Something went wrong: {}".format(err))

if __name__ == '__main__':
    crawlerMovie()
    scheduler = BlockingScheduler()
    scheduler.add_job(crawlerMovie, 'interval', days=1)
    scheduler.start()

