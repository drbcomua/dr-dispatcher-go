package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	hash "github.com/theTardigrade/golang-hash"
)

type name struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var NameStorage = map[uint64]name{}

var mutex = &sync.Mutex{}

func nameToHash(n name) uint64 {
	return hash.Uint64String(n.FirstName) ^ hash.Uint64String(n.LastName)
}

func postNames(c *gin.Context) {
	var Names []name
	var resp = map[uint64]name{}
	var temp = map[uint64]name{}

	if err := c.BindJSON(&Names); err != nil {
		return
	}

	mutex.Lock()
	for _, element := range Names {
		if _, ok := NameStorage[nameToHash(element)]; ok {
			resp[nameToHash(element)] = element
		} else {
			temp[nameToHash(element)] = element
		}
	}

	for index, element := range temp {
		NameStorage[index] = element
	}
	mutex.Unlock()

	respSlice := make([]name, 0, len(resp))

	for _, value := range resp {
		respSlice = append(respSlice, value)
	}

	c.IndentedJSON(http.StatusOK, respSlice)
}

func deleteNames(c *gin.Context) {
	var Names []name

	if err := c.BindJSON(&Names); err != nil {
		return
	}

	mutex.Lock()
	for _, element := range Names {
		delete(NameStorage, nameToHash(element))
	}
	mutex.Unlock()
}

func main() {
	router := gin.Default()
	router.POST("/names", postNames)
	router.DELETE("/names", deleteNames)

	router.Run(":9000")
}
