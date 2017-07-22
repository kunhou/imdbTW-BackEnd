package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	var (
		server      string
		port        int
		database    string
		user        string
		password    string
		server_port int
	)
	viper.SetConfigName("movie")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...")
	} else {
		server = viper.GetString("db.server")
		port = viper.GetInt("db.port")
		database = viper.GetString("db.database")
		user = viper.GetString("db.user")
		password = viper.GetString("db.password")
		server_port = viper.GetInt("server.port")
	}
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, server, port, database)
	fmt.Print(connString)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		fmt.Print(err.Error())
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	type Movie struct {
		ID          int    `json:"id"`
		ChName      string `json:"cname"`
		EnName      string `json:"ename"`
		ReleaseTime string `json:"releaseTime"`
		Type        string `json:"type"`
		Duration    string `json:"duration"`
		Director    string `json:"director"`
		Actor       string `json:"actor"`
		Company     string `json:"company"`
		Website     string `json:"website"`
		Score       string `json:"score"`
		Intro       string `json:"intro"`
		ImgPath     string `json:"imgPath"`
	}
	router := gin.Default()
	// Add API handlers here
	router.GET("/this_week", func(c *gin.Context) {
		var (
			movie Movie
		)
		rows, err := db.Query("select id, cname, ename, releaseTime, type, duration, director, actor, company, website, score, intro, imgPath from movieList WHERE YEARWEEK(`releaseTime`, 1) = YEARWEEK(CURDATE(), 1);")
		if err != nil {
			fmt.Print(err.Error())
		}
		movies := make([]Movie, 0)
		for rows.Next() {
			err = rows.Scan(&movie.ID, &movie.ChName, &movie.EnName, &movie.ReleaseTime, &movie.Type, &movie.Duration, &movie.Director, &movie.Actor, &movie.Company, &movie.Website, &movie.Score, &movie.Intro, &movie.ImgPath)
			movies = append(movies, movie)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		jsonGoesToHTML(c, http.StatusOK, gin.H{
			"result": movies,
			"count":  len(movies),
		})
	})
	router.GET("/other", func(c *gin.Context) {
		var (
			movie Movie
		)
		rows, err := db.Query("select id, cname, ename, releaseTime, type, duration, director, actor, company, website, score, intro, imgPath from movieList WHERE YEARWEEK(`releaseTime`, 1) < YEARWEEK(CURDATE(), 1);")
		if err != nil {
			fmt.Print(err.Error())
		}
		movies := make([]Movie, 0)
		for rows.Next() {
			err = rows.Scan(&movie.ID, &movie.ChName, &movie.EnName, &movie.ReleaseTime, &movie.Type, &movie.Duration, &movie.Director, &movie.Actor, &movie.Company, &movie.Website, &movie.Score, &movie.Intro, &movie.ImgPath)
			movies = append(movies, movie)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		jsonGoesToHTML(c, http.StatusOK, gin.H{
			"result": movies,
			"count":  len(movies),
		})
	})
	router.GET("/movie/:id", func(c *gin.Context) {
		id := c.Param("id")
		var (
			movie Movie
		)
		rows, err := db.Query("select id, cname, ename, releaseTime, type, duration, director, actor, company, website, score, intro, imgPath from movieList WHERE id = ?", id)
		if err != nil {
			fmt.Print(err.Error())
		}
		movies := make([]Movie, 0)
		for rows.Next() {
			err = rows.Scan(&movie.ID, &movie.ChName, &movie.EnName, &movie.ReleaseTime, &movie.Type, &movie.Duration, &movie.Director, &movie.Actor, &movie.Company, &movie.Website, &movie.Score, &movie.Intro, &movie.ImgPath)
			movies = append(movies, movie)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		jsonGoesToHTML(c, http.StatusOK, gin.H{
			"result": movies,
			"count":  len(movies),
		})
	})
	router.Static("/static", "../static")
	router.Run(":" + strconv.Itoa(server_port))
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func jsonGoesToHTML(c *gin.Context, code int, obj interface{}) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Status(code)
	writeContentType(c.Writer, []string{"application/json; charset=utf-8"})
	enc := json.NewEncoder(c.Writer)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(obj); err != nil {
		panic(err)
	}
}
