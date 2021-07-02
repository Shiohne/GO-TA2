package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type tmsg struct {
	Knn   KNN
	Out   chan<- []string
	TestX [][]float64
}
type preds struct {
	Predicts []string
}

func (msg *tmsg) read(knn KNN, testX [][]float64) []string {
	msg.Knn = knn
	msg.TestX = testX
	coms()
	return predictiones

}

var client string
var predictiones []string

func coms() {

	fmt.Print("Enter port: ")
	client = "localhost:8081"

	fmt.Print("Remote port: ")
	host := "localhost:8080"

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
	//  RECIBE BACK
	out := make(chan []string)
	msg.Out = out
	send(msg)

	dec := json.NewDecoder(conn)
	var predictions preds
	if err := dec.Decode(&predictions); err != nil {
		log.Println("Can't decode from", conn.RemoteAddr())
	} else {
		fmt.Println(predictions)
		predictiones = predictions.Predicts

	}
}

func send(msg tmsg) {
	conn, _ := net.Dial("tcp", client)
	defer conn.Close()
	fmt.Println("Sending to", conn.RemoteAddr())
	enc := json.NewEncoder(conn)
	enc.Encode(
		msg)
}
