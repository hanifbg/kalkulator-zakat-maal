package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/json"
	"path"
)

func handleSave(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        decoder := json.NewDecoder(r.Body)
        payload := struct {
            Gold_value   int 	`json:"gold_value"`
            Total_wealth int 	`json:"wealth"`
            //Gender string `json:"gender"`
        }{}
        if err := decoder.Decode(&payload); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        nisab := 85*payload.Gold_value

        if payload.Total_wealth < nisab {
	        message := fmt.Sprintf(
	            "Total harta Anda %d, dan Nisab saat ini %d. Karena Total harta lebih kecil dari nisab anda tidak perlu membayar zakat", 
	            payload.Total_wealth,
	            nisab, 
	        )
        w.Write([]byte(message))
        return
    	} else {
    		zakat := payload.Total_wealth * 25 / 100
	        message := fmt.Sprintf(
	            "Total harta Anda %d, dan Nisab saat ini %d. Karena total harta lebih besar dari nisab jadi zakat maal yang harus anda bayarkan sebesar %d", 
	            payload.Total_wealth,
	            nisab,
	            zakat,
	        )
        w.Write([]byte(message))
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
    http.HandleFunc("/save", handleSave)

    http.Handle("/static/", 
        http.StripPrefix("/static/", 
            http.FileServer(http.Dir("assets"))))

    fmt.Println("server started at localhost:9000")
    http.ListenAndServe(":9000", nil)
}