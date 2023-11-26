package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"main.go/post05-main/post05-main"
)

var MIN = 0
var MAX = 26

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getString(length int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == length {
			break
		}
		i++
	}
	return temp
}

func handler(w http.ResponseWriter, r *http.Request) {
	post05.Hostname = "localhost"
	post05.Port = 5433
	post05.Username = "postgres"
	post05.Password = "root"
	post05.Database = "gocourses"

	data, err := post05.ListCourses()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range data {
		fmt.Println(v)
	}

	SEED := time.Now().Unix()
	rand.Seed(SEED)
	random_coursename := getString(5)

	t := post05.MSDSCourse{
		CID:     random_coursename,
		CNAME:   "Business Leadership",
		CPREREQ: "Business Fundamentals"}

	id := post05.AddCourse(t)
	if id == -1 {
		fmt.Println("There was an error adding course", t.CID)
	}

	err = post05.DeleteCourse(id)
	if err != nil {
		fmt.Println(err)
	}

	// Trying to delete it again!
	err = post05.DeleteCourse(id)
	if err != nil {
		fmt.Println(err)
	}

	id = post05.AddCourse(t)
	if id == -1 {
		fmt.Println("There was an error adding course", t.CID)
	}

	t = post05.MSDSCourse{
		CID:     random_coursename,
		CNAME:   "Business Leadership",
		CPREREQ: "Business Fundamentals"}

	err = post05.UpdateCourse(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, "Database is updated, added a new course in database")
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Server listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
