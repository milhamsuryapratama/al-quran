package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Resdata struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   Data
}

type Data struct {
	Number        int    `json:"number"`
	Name          string `json:"name"`
	EnglishName   string `json:"englishName"`
	NumberOfAyahs int    `json:"numberOfAyahs"`
	Ayahs         []Ayahs
}

type Ayahs struct {
	Text          string `json:"text"`
	NumberInSurah int    `json:"numberInSurah"`
}

//var myClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	var tmpl, err = template.ParseGlob("views/*")
	if err != nil {
		panic(err.Error())
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response, err := http.Get("http://api.alquran.cloud/v1/surah/1/ar.alafasy")
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return
		}

		var res Resdata

		err = json.NewDecoder(response.Body).Decode(&res)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = tmpl.ExecuteTemplate(w, "index", res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	fmt.Println("server started")
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
