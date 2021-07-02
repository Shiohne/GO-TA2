package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"sort"
)

// Estructura a ordenar
type Slice struct {
	sort.Interface
	idx []int
}

// Función facilitadora del sort
func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

// Función para ordenar el slice de distancias float64
func sortSliceDistances(distances []float64) *Slice {
	unsortSlice := sort.Float64Slice(distances)
	sortedSlice := &Slice{Interface: unsortSlice, idx: make([]int, unsortSlice.Len())}
	for i := range sortedSlice.idx {
		sortedSlice.idx[i] = i
	}
	return sortedSlice
}

// Facilitador para ordenar las predicciones
type Prediction struct {
	label string
	count int
}

// Calcular la distancia euclidiana de dos entradas
func Euclidian(source, dest []float64) float64 {
	distance := 0.0
	for i := range source {
		distance += math.Pow(source[i]-dest[i], 2)
	}
	return math.Sqrt(distance)
}

type KNN struct {
	k        int
	data     [][]float64
	labels   []string
	accuracy float64
}

func (knn *KNN) nearestNeighbors(source []float64) map[string]int {
	nearest := []string{}
	counter := map[string]int{}
	distances := []float64{}

	// Calcular distancia entre dato de entrada y los datos de entrenamiento
	for _, dest := range knn.data {
		distances = append(distances, Euclidian(source, dest))
	}

	// Tomar el índice de los vecinos más cercanos
	kNeighborsSlice := sortSliceDistances(distances)
	sort.Sort(kNeighborsSlice)
	neighbors := kNeighborsSlice.idx[:knn.k]

	// Listar los labels más cercanos según su índice
	for _, index := range neighbors {
		nearest = append(nearest, knn.labels[index])
	}

	// Contar cantidad de veces que se repite el label más cercano
	for _, elem := range nearest {
		counter[elem] += 1
	}
	// Devuelve el mapa[label]{contador}
	return counter
}

func sortHighestLabel(counter map[string]int) string {

	prediction := []Prediction{}

	for label, count := range counter {
		prediction = append(prediction, Prediction{label, count})
	}

	// Ordenar las predicciones según valor más alto
	sort.Slice(prediction, func(i, j int) bool {
		return prediction[i].count > prediction[j].count
	})

	// Regresar el label más repetido
	return prediction[0].label
}

// Predecir los labels
func predict(knn *KNN, out chan<- []string, testX [][]float64) []string {
	predictions := []string{}
	for _, source := range testX {

		// Contar vecinos cercanos y devolver en mapa[label]{contador}
		neighborsCounter := knn.nearestNeighbors(source)

		// Determinar label más repetido
		highestNeighbor := sortHighestLabel(neighborsCounter)

		// Agregar label más repetido a las predicciones
		predictions = append(predictions, highestNeighbor)
	}
	return predictions

}

type tmsg struct {
	Knn   KNN
	Out   chan<- []string
	TestX [][]float64
}
type preds struct {
	Predicts []string
}

var server string

func coms() {
	fmt.Print("Enter port: ")
	host := "localhost:8081"

	fmt.Print("Remote port: ")
	server = "localhost:8080"

	// Listener!
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var msg tmsg
	if err := dec.Decode(&msg); err != nil {
		log.Println("Can't decode from", conn.RemoteAddr())
	} else {
		fmt.Println(msg)

		predictions := preds{}

		predictions.Predicts = predict(&msg.Knn, msg.Out, msg.TestX)
		send(predictions)

	}
}

func send(predictions preds) {
	conn, _ := net.Dial("tcp", server)
	defer conn.Close()
	fmt.Println("Sending to", conn.RemoteAddr())
	enc := json.NewEncoder(conn)
	enc.Encode(predictions)
}
