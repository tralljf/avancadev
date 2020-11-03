package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Discount struct {
	Discount string
}


func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	result := Discount{Discount: "8"}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}
