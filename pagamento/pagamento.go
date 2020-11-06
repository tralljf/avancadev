package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/wesleywillians/go-rabbitmq/queue"
)


type Result struct {
	Status string
}

type Order struct {
	ID       uuid.UUID
	Coupon string
	CcNumber string
}


const (
	InvalidCoupon = "invalid"
	ValidCoupon = "valid"
	ConnectionError = "connection error"
)

func NewOrder() Order {
	return Order{ID: uuid.NewV4()}
}

func init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}
}

//TRABALHANDO COM WEBSERVER REST 
// func main(){
// 	http.HandleFunc("/", home)
// 	http.ListenAndServe(":9091", nil)
// }


func main(){
	messageChannel := make(chan amqp.Delivery)

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	for msg := range messageChannel {
		process(msg)
	}
}


func process(msg amqp.Delivery) {
	order := NewOrder()

	json.Unmarshal(msg.Body, &order)
	resultCoupon := makeHttpCall("http://localhost:9092/", order.Coupon)


	switch resultCoupon.Status {
	case InvalidCoupon:
		log.Println("Order: ", order.ID, ": invalid coupon!")
	
	case ConnectionError:
		msg.Reject(false)
		log.Println("Order: ", order.ID, ": could not process!")

	case ValidCoupon:
		log.Println("Order: ", order.ID, ": Processed")
	}


}


//TRABALHANDO COM REST 
// func home(w http.ResponseWriter, r *http.Request){

// 	coupon := r.FormValue("coupon")
// 	ccNumber := r.FormValue("ccNumber")

// 	result := Result{Status: "declined"}

// 	if ccNumber != "" {
// 		result.Status = "approved"
// 	}

// 	if coupon != ""{
// 		resultCoupon := makeHttpCall("http://localhost:9092/", coupon)
// 		if resultCoupon.Status == "invalid" {
// 			result.Status = "invalid coupoun"
// 		}
// 	}
	
// 	jsonData, err := json.Marshal(result)
// 	if err != nil {
// 		log.Fatal("Error processar JSON")
// 	}

// 	fmt.Fprint(w, string(jsonData))
// }



func makeHttpCall(urlMicroservice string, coupon string  ) Result {

	values := url.Values{}
	values.Add("coupon", coupon)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)

	if err != nil {
		result := Result{Status: ConnectionError}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		result := Result{Status: "Erro ao processar o resultado"}
		return result
	}

	result := Result{}

	json.Unmarshal(data, &result)


	return result


}