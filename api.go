package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

//estructura
type Alumno struct {
	Codigo string `json:"cod"`
	Nombre string `json:"nom"`
	Dni    int    `json:"dni"`
}

//global
var alumnos []Alumno

func resuelveBuscarAlumno(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /alumno")

	//recuperar los parametros x querystring
	sDni := req.FormValue("dni")

	//tipo de contenido de respuesta
	res.Header().Set("Content-Type", "application/json")

	//logica del endpoint
	iDni, _ := strconv.Atoi(sDni)
	for _, alumno := range alumnos {
		if alumno.Dni == iDni {
			//codificarlo
			jsonBytes, _ := json.MarshalIndent(alumno, "", "")
			io.WriteString(res, string(jsonBytes))
		}
	}
}

func resuelveCreditos(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /creditos")
	res.Header().Set("Content-Type", "text/html")
	io.WriteString(res,
		`<doctype html>
	<html>
	<head><title>API</title></head>
	<body>
	<h2>API desarrolado para el curso de programacion concurrente y distribuida</h2>
	</body>
	</html>
	`)
}

func resuelveDataSet(res http.ResponseWriter, req *http.Request) {
	log.Println("llamada al endpoint /dataset")
	res.Header().Set("Content-Type", "application/json")
	ds := DataSet{}
	ds = fillDataSet()
	jsonBytes, _ := json.MarshalIndent(ds.label, "", "")
	io.WriteString(res, string(jsonBytes))

}

func manejadorRequest() {
	//definir los endpoints de nuestro servicio
	http.HandleFunc("/alumno", resuelveBuscarAlumno)
	http.HandleFunc("/creditos", resuelveCreditos)
	http.HandleFunc("/dataset", resuelveDataSet)

	//establecer el puerto de servicio
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	//X := [][]float64{}
	//Y := []string{}
	manejadorRequest()
	/*ds := DataSet{}
	ds = fillDataSet()*/

}
