package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-gota/gota/dataframe"
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

func (knn *KNN) fit(X [][]float64, Y []string) {
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
		fmt.Println(len(predictedLabel))
	}
	return predictedLabel

}

func knnDemo(dataX [][]float64, dataY []string, K int) {
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
			fmt.Println(float64(len(dataX)) * 0.2)
		}
		if i < float64(len(dataX))*0.2 {
			testX = append(testX, dataX[int(i)])
			testY = append(testY, dataY[int(i)])
		} else {
			trainX = append(trainX, dataX[int(i)])
			trainY = append(trainY, dataY[int(i)])
		}
	}

	//training
	knn := KNN{}
	knn.k = K
	knn.fit(trainX, trainY)
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

func knn(dataX [][]float64, dataY []string, testX [][]float64, K int) []string {
	//split data into training and test
	var (
		trainX [][]float64
		trainY []string
	)
	for i := range dataX {
		trainX = append(trainX, dataX[i])
		trainY = append(trainY, dataY[i])
	}

	//training
	knn := KNN{}
	knn.k = K
	knn.fit(trainX, trainY)
	predicted := knn.predict(testX)

	predictions := []string{}

	fmt.Printf("Usando K = %d vecinos\n", K)
	fmt.Println("Predicciones:")
	for i, label := range predicted {
		predictions = append(predictions, fmt.Sprintf("Para la paciente %d recomiendo el método %s", i+1, label))
	}
	return predictions
}

func readDataSet() (dataX, dataY [][]string) {
	url := "https://github.com/Shiohne/GO-TA2/raw/master/DAT%20PlaniFamiliar_01_Metodo.csv"
	metodoMatrix := [][]string{}
	metodo, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer metodo.Body.Close()
	br := bufio.NewReader(metodo.Body)
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
	df := dataframe.ReadCSV(br)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		metodoMatrix = append(metodoMatrix, record)
	}
	X := df.Select([]int{6, 9, 10, 11})
	Y := df.Select(7)
	dataX := X.Records()
	dataY := Y.Records()
	fmt.Println(dataX)
	fmt.Println(dataY)

	return metodoMatrix

}

type Metodo struct {
	Edad      float64 `json:"edad"`
	Tipo      float64 `json:"tipo"`
	Actividad float64 `json:"actividad"`
	Insumo    float64 `json:"insumo"`
	Metodo    string  `json:"metodo"`
}

type DataSet struct {
	Metodos []Metodo
	Data    [][]float64
	Labels  []string
}

func (ds *DataSet) loadData() {

	// Carga el DataSet desde su CSV
	metodoMatrix := readDataSet()

	// Se inicializa el metodo Struct para llenarlo con datos
	metodo := Metodo{}

	// X para la data del DataSet y Y para el Label

	for i, data := range metodoMatrix {
		// Si es que el DataSet contiene una primera fila de títulos
		if i == 0 {
			continue
		}

		temp := []float64{}
		// Convertimos los datos necesarios a floats para poder añadirlos
		for j, value := range data[:] {

			if j == 6 {
				switch value {
				case "12 a - 17 a":
					metodo.Edad = 14.5
					break
				case "18 a - 29 a":
					metodo.Edad = 23.5
					break
				case "30 a - 59 a":
					metodo.Edad = 44.5
					break
				case "> 60 a":
					metodo.Edad = 65.0
					break
				}
				temp = append(temp, metodo.Edad)
			} else if j == 7 {
				// METODO
				metodo.Metodo = value
			} else if j == 9 {
				// TIPOS - NUEVAS O CONT
				switch value {
				case "NUEVAS":
					metodo.Tipo = 0.0
					break
				case "CONTINUADORAS":
					metodo.Tipo = 1.0
					break
				}
				temp = append(temp, metodo.Tipo)
			} else if j == 10 {
				parsedValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				// ACTIVIDAD
				metodo.Actividad = parsedValue
				temp = append(temp, metodo.Actividad)
			} else if j == 11 {
				parsedValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				// INSUMO
				metodo.Insumo = parsedValue
				temp = append(temp, metodo.Insumo)
			}

		}
		ds.Data = append(ds.Data, temp)
		ds.Labels = append(ds.Labels, data[7])
		// Añadimos los datos al DataSet struct ahora convertidos
		ds.Metodos = append(ds.Metodos, metodo)
	}
}

func main() {
	ds := DataSet{}
	ds.loadData()
	//knnDemo(ds.Data, ds.Labels, 5)
}
