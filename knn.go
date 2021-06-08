package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

//calculate euclidean distance betwee two slices
func Euclidian(source, dest []float64) float64 {
	distance := 0.0
	for i := range source {
		distance += math.Pow(source[i]-dest[i], 2)
	}
	return math.Sqrt(distance)
}

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

type KNN struct {
	k      int
	data   [][]float64
	labels []string
}

func (knn *KNN) loadData(X [][]float64, Y []string) {
	//read data
	knn.data = X
	knn.labels = Y
}

func (knn *KNN) predict(X [][]float64) []string {

	predictedLabel := []string{}
	for _, source := range X {
		var (
			distList   []float64
			nearLabels []string
		)
		//calculate distance between predict target data and surpervised data
		for _, dest := range knn.data {
			distList = append(distList, Euclidian(source, dest))
		}
		//take top k nearest item's index
		s := NewFloat64Slice(distList)
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
	}
	return predictedLabel

}

func knn(X [][]float64, Y []string, K int) {
	//split data into training and test
	var (
		trainX [][]float64
		trainY []string
		testX  [][]float64
		testY  []string
	)
	for i := range X {
		if i%2 == 0 {
			trainX = append(trainX, X[i])
			trainY = append(trainY, Y[i])
		} else {
			testX = append(testX, X[i])
			testY = append(testY, Y[i])
		}
	}

	//training
	knn := KNN{}
	knn.k = K
	knn.loadData(trainX, trainY)
	predicted := knn.predict(testX)

	//check accuracy
	correct := 0
	for i := range predicted {
		if predicted[i] == testY[i] {
			correct += 1
		}
	}
	fmt.Printf("Usando K = %d vecinos\n", K)
	fmt.Printf("Predicciones correctas: %d de %d \n", correct, len(predicted))
	fmt.Printf("Precisión de %0.3f%%\n", (float64(correct)/float64(len(predicted)))*100)

}

func readDataSet() [][]string {
	irisMatrix := [][]string{}
	iris, err := os.Open("iris.csv")
	if err != nil {
		panic(err)
	}
	defer iris.Close()
	br := bufio.NewReader(iris)
	r, _, err := br.ReadRune()
	if err != nil {
		panic(err)
	}
	if r != '\uFEFF' {
		br.UnreadRune()
	}

	reader := csv.NewReader(br)
	reader.Comma = ','
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		irisMatrix = append(irisMatrix, record)
	}

	return irisMatrix
}

type Iris struct {
	SepalLength float64 `json:"sepalLength"`
	SepalWidth  float64 `json:"sepalWidth"`
	PetalLength float64 `json:"petalLength"`
	PetalWidth  float64 `json:"petalWidth"`
	Species     string  `json:"species"`
}

type DataSet struct {
	Irises []Iris
	Data   [][]float64
	Labels []string
}

func (ds *DataSet) loadData() {

	// Carga el DataSet desde su CSV
	irisMatrix := readDataSet()

	// Se inicializa el Iris Struct para llenarlo con datos
	iris := Iris{}

	// X para la data del DataSet y Y para el Label

	for i, data := range irisMatrix {
		// Si es que el DataSet contiene una primera fila de títulos
		if i == 0 {
			continue
		}

		temp := []float64{}
		// Convertimos los datos necesarios a floats para poder añadirlos
		for j, value := range data[:] {
			if j != 4 {
				parsedValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				if j == 0 {
					iris.SepalLength = parsedValue
				} else if j == 1 {
					iris.SepalWidth = parsedValue
				} else if j == 2 {
					iris.PetalLength = parsedValue
				} else if j == 3 {
					iris.PetalWidth = parsedValue
				}
				temp = append(temp, parsedValue)
			}

			iris.Species = value

		}
		ds.Data = append(ds.Data, temp)
		ds.Labels = append(ds.Labels, data[4])
		// Añadimos los datos al DataSet struct ahora convertidos
		ds.Irises = append(ds.Irises, iris)

	}

}
