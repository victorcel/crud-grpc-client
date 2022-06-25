package main

import "github.com/victorcel/crud-grpc-client/clients"

func main() {
	//quit := make(chan bool, 2)
	clients.ClientElasticSearch()
	clients.ClientMongoDB()
	//<-quit
}
