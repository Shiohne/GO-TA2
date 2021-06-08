package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var irisData DataSet

func resuelveDataSet(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /dataset")

	jsonBytes, _ := json.Marshal(irisData.Irises)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)

}
func resuelveData(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /data")

	jsonBytes, _ := json.Marshal(irisData.Data)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)

}
func resuelveLabel(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /labels")

	jsonBytes, _ := json.Marshal(irisData.Labels)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)

}

func resuelveKNN(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /knn")

	bodyBytes, _ := ioutil.ReadAll(req.Body)

	//res.Header().Set("Content-Type", "application/json")
	var iris []Iris
	json.Unmarshal(bodyBytes, &iris)
	log.Println(iris[1].PetalLength, iris[1].PetalWidth, iris[1].SepalLength, iris[0].SepalWidth, iris[0].Species)

}

func manejadorRequest() {
	// Definir los endpoints de nuestro servicio
	http.HandleFunc("/dataset", resuelveDataSet)
	http.HandleFunc("/data", resuelveData)
	http.HandleFunc("/labels", resuelveLabel)
	http.HandleFunc("/knn", resuelveKNN)

	// Establecer el puerto de servicio
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	irisData.loadData()
	manejadorRequest()

}
