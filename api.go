package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var irisData DataSet

func resuelveDataSet(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /dataset")
	res.Header().Set("Content-Type", "application/json")

	jsonBytes, _ := json.MarshalIndent(irisData.Irises, "", "")
	io.WriteString(res, string(jsonBytes))

}
func resuelveData(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /data")
	res.Header().Set("Content-Type", "application/json")

	jsonBytes, _ := json.MarshalIndent(irisData.Data, "", "")
	io.WriteString(res, string(jsonBytes))

}
func resuelveLabel(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /labels")
	res.Header().Set("Content-Type", "application/json")

	jsonBytes, _ := json.MarshalIndent(irisData.Labels, "", "")
	io.WriteString(res, string(jsonBytes))

}

func manejadorRequest() {
	// Definir los endpoints de nuestro servicio
	http.HandleFunc("/dataset", resuelveDataSet)
	http.HandleFunc("/data", resuelveData)
	http.HandleFunc("/labels", resuelveLabel)

	// Establecer el puerto de servicio
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	irisData.loadData()
	manejadorRequest()

}
