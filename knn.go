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
func (knn *KNN) predict(testX [][]float64) []string {

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

func (knn *KNN) findOptimalK(dataX [][]float64, dataY []string) {
	// Partir la data en entrenamiento y testing
	var (
		trainX [][]float64
		trainY []string
		testX  [][]float64
		testY  []string
	)
	for i := 0.0; i < float64(len(dataX)); i++ {
		if i < float64(len(dataX))*0.01 {
			testX = append(testX, dataX[int(i)])
			testY = append(testY, dataY[int(i)])
		} else {
			trainX = append(trainX, dataX[int(i)])
			trainY = append(trainY, dataY[int(i)])
		}
	}
	knn.data = trainX
	knn.labels = trainY
	knn.train(testX, testY, 1, 10)

}

func (knn *KNN) train(testX [][]float64, testY []string, start, end int) {
	optimalK := 1
	bestAccuracy := 0.0
	for i := start; i < end; i++ {
		if i == 1 || i%2 != 0 {

			knn.k = i

			// Separar las pruebas en 4 partes
			testXPart1 := testX[:len(testX)/4]
			testXPart2 := testX[len(testX)/4 : (2*len(testX))/4]
			testXPart3 := testX[(2*len(testX))/4 : (3*len(testX))/4]
			testXPart4 := testX[(3*len(testX))/4:]

			// Crear los canales para recibir las predicciones de ambas partes
			out1 := make(chan []string)
			out2 := make(chan []string)
			out3 := make(chan []string)
			out4 := make(chan []string)
			go func(out chan<- []string) { out <- knn.predict(testXPart1) }(out1)
			go func(out chan<- []string) { out <- knn.predict(testXPart2) }(out2)
			go func(out chan<- []string) { out <- knn.predict(testXPart3) }(out3)
			go func(out chan<- []string) { out <- knn.predict(testXPart4) }(out4)
			part1, part2, part3, part4 := <-out1, <-out2, <-out3, <-out4
			close(out1)
			close(out2)
			close(out3)
			close(out4)

			// Unir las predicciones
			predictions := []string{}
			predictions = append(predictions, part1...)
			predictions = append(predictions, part2...)
			predictions = append(predictions, part3...)
			predictions = append(predictions, part4...)

			// Probar precision
			correct := 0
			for i := range predictions {
				if predictions[i] == testY[i] {
					correct += 1
				}
			}

			fmt.Printf("Usando K = %d vecinos\n", knn.k)
			fmt.Printf("Predicciones correctas: %d de %d \n", correct, len(predictions))
			accuracy := (float64(correct) / float64(len(predictions)))
			fmt.Printf("Precisión de %0.3f%%\n", accuracy*100)

			if bestAccuracy < accuracy {
				bestAccuracy = accuracy
				optimalK = knn.k
			}

		}
	}
	fmt.Println("El k más óptimo es: ", optimalK)
	knn.k = optimalK
	knn.accuracy = bestAccuracy
}

func (knn *KNN) knn(dataX [][]float64, dataY []string, testX [][]float64) Respuesta {

	// Insertar datos al knn
	knn.data = dataX
	knn.labels = dataY

	// Separar las pruebas en 2 partes
	testXPart1 := testX[:len(testX)/2]
	testXPart2 := testX[len(testX)/2:]

	// Crear los canales para recibir las predicciones de ambas partes
	out1 := make(chan []string)
	out2 := make(chan []string)
	go func(out chan<- []string) { out <- knn.predict(testXPart1) }(out1)
	go func(out chan<- []string) { out <- knn.predict(testXPart2) }(out2)
	part1, part2 := <-out1, <-out2
	close(out1)
	close(out2)

	// Unir las predicciones
	predictions := []string{}
	predictions = append(predictions, part1...)
	predictions = append(predictions, part2...)

	// Inicializar las estructuras que reciben los resultados
	resultado := Resultado{}
	respuesta := Respuesta{}

	// Regresar respuesta con los resultados para mostrar
	for i, label := range predictions {
		resultado.Prediccion = fmt.Sprintf("Para la usuaria %d recomiendo el método %s", i+1, label)
		respuesta.Resultados = append(respuesta.Resultados, resultado)
	}
	respuesta.Detalles = fmt.Sprintf("Usando K = %d vecinos para las %d usuarias y con un porcentaje de precisión de  %0.3f%%", knn.k, len(predictions), knn.accuracy*100)
	return respuesta
}

func readDataSet() [][]string {
	// Obtener el dataset desde github
	metodoMatrix := [][]string{}
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

	// Leer el dataset
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
		metodoMatrix = append(metodoMatrix, record)
	}

	return metodoMatrix
}

type Resultado struct {
	Prediccion string `json:"prediccion"`
}

type Respuesta struct {
	Detalles   string      `json:"detalles"`
	Resultados []Resultado `json:"resultados"`
}

type Usuaria struct {
	Edad      float64 `json:"edad"`
	Tipo      float64 `json:"tipo"`
	Actividad float64 `json:"actividad"`
	Insumo    float64 `json:"insumo"`
	Metodo    string  `json:"metodo"`
}

type DataSet struct {
	Usuarias []Usuaria `json:"usuarias"`
	Data     [][]float64
	Labels   []string
}

func (ds *DataSet) loadData() {

	// Cargar el DataSet desde su CSV
	data := readDataSet()

	// Inicializar la usuaria Struct para llenarlo con datos
	usuaria := Usuaria{}

	// Almacenar los datos en las estructuras
	for i, metodos := range data {
		// Drop de la primera fila (titles)
		if i == 0 {
			continue
		}

		temp := []float64{}
		// Convertimos los datos necesarios a floats para poder añadirlos
		for j, value := range metodos[:] {

			if j == 6 {
				switch value {
				case "12 a - 17 a":
					usuaria.Edad = 14.5
				case "18 a - 29 a":
					usuaria.Edad = 23.5
				case "30 a - 59 a":
					usuaria.Edad = 44.5
				case "> 60 a":
					usuaria.Edad = 65.0
				}
				temp = append(temp, usuaria.Edad)
			} else if j == 7 {
				// METODO
				usuaria.Metodo = value
			} else if j == 9 {
				// Si son Nuevas = 0 y si son Continuadoras = 1
				switch value {
				case "NUEVAS":
					usuaria.Tipo = 0.0
				case "CONTINUADORAS":
					usuaria.Tipo = 1.0
				}
				// TIPO DE USUARIA
				temp = append(temp, usuaria.Tipo)
			} else if j == 10 {
				parsedValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				// ACTIVIDAD
				usuaria.Actividad = parsedValue
				temp = append(temp, usuaria.Actividad)
			} else if j == 11 {
				parsedValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				// INSUMO
				usuaria.Insumo = parsedValue
				temp = append(temp, usuaria.Insumo)
			}

		}
		// Filtramos todas las filas que contengan MELA ya que no es un Metodo anticonceptivo que se pueda recomendar normalmente
		if metodos[7] != "MELA" {

			// Añadir los datos al DataSet struct ahora convertidos
			ds.Data = append(ds.Data, temp)
			ds.Labels = append(ds.Labels, metodos[7])
			ds.Usuarias = append(ds.Usuarias, usuaria)
		}
	}
}

/*func main() {
	ds := DataSet{}
	ds.loadData()
	knn := KNN{}
	knn.findOptimalK(ds.Data, ds.Labels)

}*/
