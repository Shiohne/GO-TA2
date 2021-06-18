package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var pacienteData DataSet

func jsonToSlice(pacientesJSON []Paciente) [][]float64 {
	pacientesSlice := [][]float64{}
	for i := range pacientesJSON {
		paciente := []float64{pacientesJSON[i].Edad, pacientesJSON[i].Tipo, pacientesJSON[i].Actividad, pacientesJSON[i].Insumo}
		pacientesSlice = append(pacientesSlice, paciente)
	}
	return pacientesSlice

}

func resuelveDataSet(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /dataset")

	jsonBytes, _ := json.Marshal(pacienteData.Pacientes)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(jsonBytes)

}

func resuelveKNN(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /knn")

	bodyBytes, _ := ioutil.ReadAll(req.Body)

	pacientesJSON := []Paciente{}
	json.Unmarshal(bodyBytes, &pacientesJSON)

	log.Println(pacientesJSON)

	// Transformar pacientes a Slice de Slices para pasar al KNN
	pacientesSlice := jsonToSlice(pacientesJSON)
	predicciones := knn(pacienteData.Data, pacienteData.Labels, pacientesSlice)
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
	pacienteData.loadData()
	manejadorRequest()
}
