package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()

	router.GET("/", Home)
	fmt.Println("start...")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	fmt.Println("HOME")
}
