package filter

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Filter struct {
	Id          string `json:"id"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

type Cache struct {
	Filter []Filter
}

func ReadFilterFile(filename string) *Cache {
	cache := new(Cache)

	filterFile, err := ioutil.ReadFile("configs/" + filename)
	if err != nil {
		log.Printf("Reading file error %v\n", err)
	}
	err = json.Unmarshal(filterFile, cache)

	if err != nil {
		log.Fatalf("Unmarshal: %v\n", err)
	}

	return cache
}
