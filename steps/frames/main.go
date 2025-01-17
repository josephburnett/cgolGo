package main

import (
	"encoding/json"
	"flag"
	"os"
)

var (
	name       = flag.String("name", "", "String name of output")
	valueFile  = flag.String("value-file", "", "File containing value JSON")
	outputFile = flag.String("output-file", "", "Step output file")
)

func main() {
	validate()
	data, err := os.ReadFile(*valueFile)
	if err != nil {
		panic(err)
	}
	var value []any
	err = json.Unmarshal(data, &value)
	if err != nil {
		panic(err)
	}
	filenames := []string{}
	for _, frame := range value {
		func() {
			file, err := os.CreateTemp("", "frame-*.json")
			if err != nil {
				panic(err)
			}
			filenames = append(filenames, file.Name())
			defer file.Close()
			encoder := json.NewEncoder(file)
			err = encoder.Encode(frame)
			if err != nil {
				panic(err)
			}
		}()
	}
	output := map[string]any{
		"name":  *name,
		"value": filenames,
	}
	data, err = json.Marshal(output)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(*outputFile, data, 0644)
	if err != nil {
		panic(err)
	}
}

func validate() {
	flag.Parse()
	if *name == "" {
		panic("name is required")
	}
	if *valueFile == "" {
		panic("value-file is required")
	}
	if *outputFile == "" {
		panic("output-file is required")
	}
}
