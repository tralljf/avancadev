package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/wesleywillians/go-rabbitmq/queue"
)


type Result struct {
	Status string
}

type Order struct {
	Coupon string
	CcNumber string
}

func init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}
}

func main(){
	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9090", nil)
}

func home(w http.ResponseWriter, r *http.Request){
	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, Result{})
}

func process(w http.ResponseWriter, r *http.Request){

	coupon := r.FormValue("coupon") 
	ccNumber := r.FormValue("cc-number")


	order := Order {
		Coupon: coupon,
		CcNumber: ccNumber,
	}

	jsonOrder, err := json.Marshal(order)

	if err != nil {
		log.Fatal("Erro ao converter o Json")
	}

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()

	defer ch.Close()

	err = rabbitMQ.Notify(string(jsonOrder), "application/json", "orders_ex", "")
	t := template.Must(template.ParseFiles("templates/process.html"))
	
	if err != nil {
		t.Execute(w, "Microserviço não conseguiu conectar na fila")	
		// log.Fatal("Erro ao adicionar a fila")
	}

	// result := makeHttpCall("http://localhost:9091", r.FormValue("coupon"), r.FormValue("cc-number"))

	t.Execute(w, "")
}


// func makeHttpCall(urlMicroservice string, coupon string, ccNumber string  ) Result {

// 	values := url.Values{}
// 	values.Add("coupon", coupon)
// 	values.Add("ccNumber", ccNumber)

// 	retryClient := retryablehttp.NewClient()
// 	retryClient.RetryMax = 5

// 	res, err := retryClient.PostForm(urlMicroservice, values)

// 	if err != nil {
// 		// log.Fatal("Microserviço de pagamento retornou erro")
// 		result := Result{Status: "Servidor fora do ar"}
// 		return result
// 	}

// 	defer res.Body.Close()

// 	data, err := ioutil.ReadAll(res.Body)

// 	if err != nil {
// 		// log.Fatal("Erro ao processar o resultado")
// 		result := Result{Status: "Erro ao processar o resultado"}
// 		return result
// 	}

// 	result := Result{}

// 	json.Unmarshal(data, &result)


// 	return result


// }