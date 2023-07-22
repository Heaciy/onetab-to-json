package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	dbPath  string
	outPath string
)

func init() {

	flag.StringVar(&dbPath, "i", "", "leveldb path")
	flag.StringVar(&outPath, "o", "tabs.json", "output json file")
	flag.Parse()
}

func decode(key, value []byte) (string, interface{}) {
	k := string(key)
	v := string(value)
	var decodedValue interface{}
	var err error

	if k == "settings" || k == "state" {
		// For "settings" and "state" keys, perform two JSON decoding operations.
		var intermediateValue string
		if err = json.Unmarshal([]byte(v), &intermediateValue); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return k, nil
		}
		if err = json.Unmarshal([]byte(intermediateValue), &decodedValue); err != nil {
			fmt.Println("Error decoding nested JSON:", err)
			return k, nil
		}
	} else {
		err = json.Unmarshal([]byte(v), &decodedValue)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return k, nil
		}
	}

	return k, decodedValue
}

func encodeJSON(v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func parser(db *leveldb.DB, outPath string) {
	data := make(map[string]interface{})
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		k, v := decode(key, value)
		data[k] = v
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		fmt.Println("Error during iteration:", err)
		return
	}

	jsonData, err := encodeJSON(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = ioutil.WriteFile(outPath, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func main() {
	// path to levelDB
	if dbPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// open levelDB
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// parser data
	parser(db, outPath)
}
