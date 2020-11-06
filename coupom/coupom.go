package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}


type Result struct {
	Status string
}

var coupons Coupons

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}	

	return "invalid"
}


func main(){
	coupon := Coupon {
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)

}

func home(w http.ResponseWriter, r *http.Request){
	coupon := r.FormValue("coupon")

	valid := coupons.Check(coupon)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)

	if err != nil {
		log.Fatal("Erro ao converter string")
	}

	fmt.Fprint(w, string(jsonResult))
}