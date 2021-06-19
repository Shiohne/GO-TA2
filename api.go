package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var usuariaData DataSet

func jsonToSlice(usuariasJSON []Usuaria) [][]float64 {
	usuariasSlice := [][]float64{}
	for i := range usuariasJSON {
		usuaria := []float64{usuariasJSON[i].Edad, usuariasJSON[i].Tipo, usuariasJSON[i].Actividad, usuariasJSON[i].Insumo}
		usuariasSlice = append(usuariasSlice, usuaria)
	}
	return usuariasSlice

}

func resuelveDataSet(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /dataset")

	jsonBytes, _ := json.Marshal(usuariaData.Usuarias)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)

}

func resuelveKNN(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /knn")

	bodyBytes, _ := ioutil.ReadAll(req.Body)

	usuariasJSON := []Usuaria{}
	json.Unmarshal(bodyBytes, &usuariasJSON)

	log.Println(usuariasJSON)

	// Transformar usuarias a Slice de Slices para pasar al KNN
	usuariasSlice := jsonToSlice(usuariasJSON)
	predicciones := knn(usuariaData.Data, usuariaData.Labels, usuariasSlice)
	jsonBytes, _ := json.Marshal(predicciones)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)
}

func manejadorRequest() {
	// Definir los endpoints del servicio
	http.HandleFunc("/api/dataset", resuelveDataSet)
	http.HandleFunc("/api/knn", resuelveKNN)

	// Establecer el puerto de servicio
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	usuariaData.loadData()
	manejadorRequest()
}
