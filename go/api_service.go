package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/movie")
	if err != nil {
		fmt.Print(err.Error())
	}
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
	router.Static("/static", "../static")
	router.Run(":8080")
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
