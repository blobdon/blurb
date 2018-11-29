package main

import (
	"encoding/json"
	"io/ioutil"
)

func bookToFile(b book, filename string) error {
	j, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, j, 0666)
}

func fileToBook(filename string) (book, error) {
	b := book{}
	j, err := ioutil.ReadFile(filename)
	if err != nil {
		return b, err
	}
	return b, json.Unmarshal(j, &b)
}

func mapToFile(source map[string]struct{}, filename string) error {
	j, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, j, 0666)
}

func fileToMap(filename string) (map[string]struct{}, error) {
	m := map[string]struct{}{}
	j, err := ioutil.ReadFile(filename)
	if err != nil {
		return m, err
	}
	return m, json.Unmarshal(j, &m)
}
