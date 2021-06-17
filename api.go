package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var metodoData DataSet

func resuelveDataSet(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /dataset")

	jsonBytes, _ := json.Marshal(metodoData.Metodos)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)

}

func resuelveKNN(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /knn")

	bodyBytes, _ := ioutil.ReadAll(req.Body)

	metodoJSON := []Metodo{}
	json.Unmarshal(bodyBytes, &metodoJSON)

	log.Println(metodoJSON)
	metodoMap := [][]float64{}

	for i := range metodoJSON {
		metodo := []float64{metodoJSON[i].Edad, metodoJSON[i].Tipo, metodoJSON[i].Actividad, metodoJSON[i].Insumo}
		metodoMap = append(metodoMap, metodo)
	}

	predicciones := knn(metodoData.Data, metodoData.Labels, metodoMap)

	jsonBytes, _ := json.Marshal(predicciones)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)
}

func manejadorRequest() {
	// Definir los endpoints de nuestro servicio
	http.HandleFunc("/dataset", resuelveDataSet)
	http.HandleFunc("/knn", resuelveKNN)

	// Establecer el puerto de servicio
	router := mux.NewRouter()
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}

func main() {
	metodoData.loadData()
	manejadorRequest()
}
