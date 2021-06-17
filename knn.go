package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-gota/gota/dataframe"
)

//argument sort
type Slice struct {
	sort.Interface
	idx []int
}

func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(n sort.Interface) *Slice {
	s := &Slice{Interface: n, idx: make([]int, n.Len())}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

func NewFloat64Slice(n []float64) *Slice { return NewSlice(sort.Float64Slice(n)) }

//map sort
type Entry struct {
	name  string
	value int
}
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return l[i].name < l[j].name
	} else {
		return l[i].value > l[j].value
	}
}

//count item frequence in slice
func Counter(target []string) map[string]int {
	counter := map[string]int{}
	for _, elem := range target {
		counter[elem] += 1
	}
	return counter
}

//calculate euclidean distance betwee two slices
func Euclidian(source, dest []float64) float64 {
	distance := 0.0
	for i := range source {
		distance += math.Pow(source[i]-dest[i], 2)
	}
	return math.Sqrt(distance)
}

type KNN struct {
	k      int
	data   [][]float64
	labels []string
}

func (knn *KNN) predict(X [][]float64) []string {

	predictedLabel := []string{}
	for _, source := range X {
		var (
			distances  []float64
			nearLabels []string
		)
		//calculate distance between predict target data and surpervised data
		for _, dest := range knn.data {
			distances = append(distances, Euclidian(source, dest))
		}
		//take top k nearest item's index
		s := NewFloat64Slice(distances)
		sort.Sort(s)
		targetIndex := s.idx[:knn.k]

		//get the index's label
		for _, ind := range targetIndex {
			nearLabels = append(nearLabels, knn.labels[ind])
		}

		//get label frequency
		labelFreq := Counter(nearLabels)

		//the most frequent label is the predict target label
		a := List{}
		for k, v := range labelFreq {
			e := Entry{k, v}
			a = append(a, e)
		}
		sort.Sort(a)
		predictedLabel = append(predictedLabel, a[0].name)
		//fmt.Println(len(predictedLabel))
	}
	return predictedLabel

}

// Funcion con la que se realizaron pruebas para hallar el K más óptimo y como resultado fue el 7
func knnDemo(dataX [][]float64, dataY []string) {
	//split data into training and test
	var (
		trainX [][]float64
		trainY []string
		testX  [][]float64
		testY  []string
	)
	for i := 0.0; i < float64(len(dataX)); i++ {
		if i == 0 {
			fmt.Println(len(dataX))
			fmt.Println(float64(len(dataX)) * 0.005)
		}
		if i < float64(len(dataX))*0.005 {
			testX = append(testX, dataX[int(i)])
			testY = append(testY, dataY[int(i)])
		} else {
			trainX = append(trainX, dataX[int(i)])
			trainY = append(trainY, dataY[int(i)])
		}
	}

	//training
	knn := KNN{}
	knn.data = trainX
	knn.labels = trainY
	bestAcc := 0.0
	bestK := 0

	for i := 1; i < 42; i++ {
		knn.k = i
		predicted := knn.predict(testX)

		//check accuracy
		correct := 0
		for i := range predicted {
			if predicted[i] == testY[i] {
				correct += 1
			}
		}
		precision := float64(correct) / float64(len(predicted))

		if bestAcc < precision {
			bestK = knn.k
		}

		fmt.Printf("Usando K = %d vecinos\n", knn.k)
		fmt.Printf("Predicciones correctas: %d de %d \n", correct, len(predicted))
		fmt.Printf("Precisión de %0.10f%%\n", precision*100)
	}
	fmt.Printf("El mejor K es de %d", bestK)
}

func knn(dataX [][]float64, dataY []string, testX [][]float64) Respuesta {

	// Insertar datos al knn
	knn := KNN{}
	knn.data = dataX
	knn.labels = dataY
	knn.k = 7

	// Predecir los métodos anticonceptivos
	predicted := knn.predict(testX)

	// Inicializar las estructuras que reciben los resultados
	resultado := Resultado{}
	respuesta := Respuesta{}

	for i, label := range predicted {
		resultado.Prediccion = fmt.Sprintf("Para la paciente %d recomiendo el método %s", i+1, label)
		respuesta.Resultados = append(respuesta.Resultados, resultado)
	}
	respuesta.Detalles = fmt.Sprintf("Usando K = %d vecinos", knn.k)
	return respuesta
}

func readDataSet() [][]string {
	// Obtiene el dataset desde el github
	url := "https://github.com/Shiohne/GO-TA2/raw/master/DAT%20PlaniFamiliar_01_Metodo.csv"
	dataset, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer dataset.Body.Close()

	// Maneja la codificación del archivo si es que hubiera
	br := bufio.NewReader(dataset.Body)
	r, _, err := br.ReadRune()
	if err != nil {
		panic(err)
	}
	if r != '\uFEFF' {
		br.UnreadRune()
	}

	// Lee el dataset
	reader := csv.NewReader(br)
	reader.Comma = ','
	reader.LazyQuotes = true

	// Almacena en un dataframe y separa en dataX(features) y dataY(labels)
	df := dataframe.ReadCSV(br)
	dfFilter := df.Select([]int{6, 9, 10, 11, 7})
	data := dfFilter.Records()
	return data

}

type Resultado struct {
	Prediccion string `json:"prediccion"`
}

type Respuesta struct {
	Detalles   string      `json:"detalles"`
	Resultados []Resultado `json:"resultados"`
}

type Metodo struct {
	Edad      float64 `json:"edad"`
	Tipo      float64 `json:"tipo"`
	Actividad float64 `json:"actividad"`
	Insumo    float64 `json:"insumo"`
	Metodo    string  `json:"metodo"`
}

type DataSet struct {
	Metodos []Metodo `json:"metodos"`
	Data    [][]float64
	Labels  []string
}

func (ds *DataSet) loadData() {

	// Cargar el DataSet desde su CSV
	data := readDataSet()

	// Inicializar el metodo Struct para llenarlo con datos
	metodo := Metodo{}

	// Almacenar los datos en las estructuras
	for i, metodos := range data {
		// Drop de la primera fila (titles)
		if i == 0 {
			continue
		}

		temp := []float64{}
		// Recorrer las columnas
		for j, value := range metodos {

			// Convertir los datos según su columna
			if j == 0 {
				// Sacar la media de las edades para estandarizar los datos
				switch value {
				case "12 a - 17 a":
					metodo.Edad = 14.5
				case "18 a - 29 a":
					metodo.Edad = 23.5
				case "30 a - 59 a":
					metodo.Edad = 44.5
				case "> 60 a":
					metodo.Edad = 65.0
				}
				// EDAD
				temp = append(temp, metodo.Edad)
			} else if j == 1 {
				// Si son Nuevas = 0 y si son Continuadoras = 1
				switch value {
				case "NUEVAS":
					metodo.Tipo = 0.0
				case "CONTINUADORAS":
					metodo.Tipo = 1.0
				}
				// TIPO DE USUARIA
				temp = append(temp, metodo.Tipo)
			} else if j == 2 {
				// int a float para facilitar operaciones
				parsedValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				// ACTIVIDAD
				metodo.Actividad = parsedValue
				temp = append(temp, metodo.Actividad)
			} else if j == 3 {
				// int a float para facilitar operaciones
				parsedValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				// INSUMO
				metodo.Insumo = parsedValue
				temp = append(temp, metodo.Insumo)
			} else if j == 4 {
				// METODO
				metodo.Metodo = value
			}

		}
		// Añadir los datos al DataSet struct ahora convertidos
		ds.Data = append(ds.Data, temp)
		ds.Labels = append(ds.Labels, metodos[4])
		ds.Metodos = append(ds.Metodos, metodo)
	}
}
