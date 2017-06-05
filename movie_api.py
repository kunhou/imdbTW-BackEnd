# coding: utf-8
import requests, os, time, csv, sys, json, mysql.connector, ConfigParser
sys.path.append("/usr/local/lib/python2.7/site-packages")
from flask import Flask, Response
from datetime import date, datetime

# Database config
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

def json_serial(obj):
    """JSON serializer for objects not serializable by default json code"""
    if isinstance(obj, (datetime, date)):
        serial = obj.isoformat()
        return serial
    raise TypeError ("Type %s not serializable" % type(obj))

def query_data(query):
    result = []
    try:
        cnx = mysql.connector.connect(**config)
        cursor = cnx.cursor()
        cursor.execute(query)
        rows = cursor.fetchall()

        columns = [desc[0] for desc in cursor.description]
        for row in rows:
            row = dict(zip(columns, row))
            result.append(row)

        cursor.close()
        cnx.close()
    except mysql.connector.Error as err:
        print("Something went wrong: {}".format(err))
    result = {
        'result': result
    }
    return result

def get_thieWeek_movie():
    query = ("SELECT * FROM movieList WHERE YEARWEEK(`releaseTime`, 1) = YEARWEEK(CURDATE(), 1)")
    result = query_data(query)
    return json.dumps(result, default=json_serial, ensure_ascii=False, encoding="utf-8")
    
def get_beforeThisWeek_movie():
    query = ("SELECT * FROM movieList WHERE YEARWEEK(`releaseTime`, 1) < YEARWEEK(CURDATE(), 1)")
    result = query_data(query)
    return json.dumps(result, default=json_serial, ensure_ascii=False, encoding="utf-8")
    
app = Flask(__name__)

@app.route('/this_week')
def thisWeek():
    resp = Response(get_thieWeek_movie())
    resp.headers['Access-Control-Allow-Origin'] = '*'
    resp.headers['Content-Type'] = 'application/json; charset=utf-8'
    return resp

@app.route('/other')
def other():
    resp = Response(get_beforeThisWeek_movie())
    resp.headers['Access-Control-Allow-Origin'] = '*'
    resp.headers['Content-Type'] = 'application/json; charset=utf-8'
    return resp

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)





