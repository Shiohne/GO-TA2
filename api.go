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
	irisJSON := []Iris{}
	json.Unmarshal(bodyBytes, &irisJSON)
	log.Println(irisJSON)
	irisX := [][]float64{}
	for i, _ := range irisJSON {
		irisI := []float64{irisJSON[i].SepalLength, irisJSON[i].SepalWidth, irisJSON[i].PetalLength, irisJSON[i].PetalWidth}
		irisX = append(irisX, irisI)
	}
	predicciones := knn(irisData.Data, irisData.Labels, irisX, 5)
	jsonBytes, _ := json.Marshal(predicciones)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)

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
