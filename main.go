package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"

	"github.com/pushm0v/golddigger"
)

func handleCount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		decoder := json.NewDecoder(r.Body)
		payload := struct {
			Wealth      int    `json:"wealth"`
			Name        string `json:"line_name"`
			Gold_amount int    `json:"gold_amount"`
		}{}
		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		price, _ := golddigger.DigHargaEmasOrg()
		nisab := 85 * price
		Total_wealth := payload.Wealth + (payload.Gold_amount * price)
		nama := payload.Name

		if Total_wealth <= nisab {
			data := []struct {
				Type         int
				Name         string
				Total_wealth int
				Nisab        int
			}{
				{0, nama, Total_wealth, nisab},
			}
			jsonInBytes, _ := json.Marshal(data)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonInBytes)
			return
		} else {
			zakat := Total_wealth * 25 / 1000
			data := []struct {
				Type         int
				Name         string
				Total_wealth int
				Nisab        int
				Zakat        int
			}{
				{1, nama, Total_wealth, nisab, zakat},
			}
			jsonInBytes, _ := json.Marshal(data)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonInBytes)
			return
		}
	}
	http.Error(w, "Only accept POST request", http.StatusBadRequest)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "index.html")
	tmpl := template.Must(template.ParseFiles(filepath))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/save", handleCount)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	fmt.Println("server started at localhost:5000")
	// err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}
