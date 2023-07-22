package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/pkg/errors"
	"bytes"
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

func parser(db *leveldb.DB, outPath string) error {
	data := make(map[string]any)
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		value := iter.Value()

		var valueData any
		err := json.Unmarshal(value, &valueData)
		if err != nil {
			fmt.Println("unmarshal data error:", err)
		}

		if key == "settings" || key == "state" {
			valueStr, ok := valueData.(string)
			if ok {
				err := json.Unmarshal([]byte(valueStr), &valueData)
				if err != nil {
					fmt.Println("unmarshal data error:", err)
				}
			}
		}

		data[key] = valueData
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return errors.Wrap(err, "error during iteration")
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return errors.Wrap(err, "marshal error")
	}
	err = os.WriteFile(outPath, buffer.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, "write file error")
	}

	return nil
}

func main() {
	// path to levelDB
	if dbPath == "" {
		flag.PrintDefaults()
		return
	}

	// open levelDB
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		fmt.Println("error opening database:", err)
		return
	}
	defer db.Close()

	// parser data
	err = parser(db, outPath)
	if err != nil {
		fmt.Println("parser data error:", err)
	}
}
