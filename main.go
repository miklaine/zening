package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// InputJSON is a struct to hold the parsed JSON data.
// The json:"-" tag indicates that this field should be ignored
// by the JSON marshaller and unmarshaller.
type InputJSON struct {
	Data map[string]interface{} `json:"-"`
}

// parseInputJSON reads and parses a JSON file into an InputJSON struct.
func parseInputJSON(filename string) (*InputJSON, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return &InputJSON{Data: data}, nil
}

// main is the entry point of the application.
// It reads the input file, processes the data, and prints the result.
func main() {
	// Input JSON file is provided as the first command-line argument
	inputFile := os.Args[1]
	inputJSON, err := parseInputJSON(inputFile)
	if err != nil {
		fmt.Println("Error parsing input JSON:", err)
		return
	}

	startTime := time.Now()

	// Process data
	transformedData := make(map[string]interface{})
	for key, value := range inputJSON.Data {
		processedValue, err := processField(key, value)
		if err != nil {
			//fmt.Println("Error processing field:", key, err)
			continue
		}
		if processedValue != nil {
			transformedData[key] = processedValue
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// Convert the transformed data to JSON and print it
	outputJSON, err := json.MarshalIndent([]interface{}{transformedData}, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling output JSON:", err)
		return
	}

	fmt.Println(string(outputJSON))
	fmt.Println("Processing time:", duration)
}
