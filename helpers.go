package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// sanitizeString trims leading and trailing whitespace from a string.
func sanitizeString(str string) string {
	return strings.TrimSpace(str)
}

// processString processes a string value according to specified criteria:
// It trims the string, converts RFC3339 formatted dates to Unix epoch time,
// and returns the sanitized string or the converted time.
func processString(value string) (interface{}, error) {
	value = sanitizeString(value)
	if value == "" {
		return nil, fmt.Errorf("empty value")
	}

	// Check if the string is an RFC3339 formatted date
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t.Unix(), nil
	}
	return value, nil
}

// processNumber converts a string to a float64 after sanitizing it.
// It returns an error for invalid numeric values.
func processNumber(value string) (interface{}, error) {
	value = sanitizeString(value)
	if value == "" {
		return nil, fmt.Errorf("empty value")
	}

	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, nil
	}
	return nil, fmt.Errorf("invalid number: %s", value)
}

// processBoolean converts a string to a boolean.
// It returns true for '1', 't', 'true', and false for '0', 'f', 'false'.
// Any other value results in an error.
func processBoolean(value string) (interface{}, error) {
	value = sanitizeString(strings.ToLower(value))
	if value == "" {
		return nil, fmt.Errorf("empty value")
	}

	switch value {
	case "1", "t", "true":
		return true, nil
	case "0", "f", "false":
		return false, nil
	default:
		return nil, fmt.Errorf("invalid boolean: %s", value)
	}
}

// processNull determines if a string should represent a null value.
// Returns nil for '1', 't', 'true', and an error for any other value.
func processNull(value string) (interface{}, error) {
	value = sanitizeString(strings.ToLower(value))
	if value == "" {
		return nil, fmt.Errorf("empty value")
	}

	switch value {
	case "1", "t", "true":
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid null: %s", value)
	}
}

// processList processes a list of values, applying the appropriate
// transformations based on the contained types.
func processList(value interface{}) ([]interface{}, error) {
	list, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected a list, got: %v", value)
	}

	var resultList []interface{}
	for _, item := range list {
		processedItem, err := processValue(item)
		if err == nil {
			resultList = append(resultList, processedItem)
		}
	}
	if len(resultList) == 0 {
		return nil, fmt.Errorf("empty list")
	}
	return resultList, nil
}

// processMap processes a map of values, applying the appropriate
// transformations for each key-value pair.
func processMap(m map[string]interface{}) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	for key, value := range m {
		processedValue, err := processField(key, value)
		if err == nil {
			resultMap[key] = processedValue
		}
	}
	if len(resultMap) == 0 {
		return nil, nil
	}
	return resultMap, nil
}

// processField processes a key-value pair by applying the appropriate
// transformation based on the value's type.
func processField(key string, value interface{}) (interface{}, error) {
	key = sanitizeString(key)
	if key == "" {
		return nil, fmt.Errorf("empty key")
	}

	processedValue, err := processValue(value)
	if err != nil {
		return nil, fmt.Errorf("invalid field: %s", key)
	}

	return processedValue, nil
}

// processValue determines the type of value and applies the
// corresponding processing function.
func processValue(value interface{}) (interface{}, error) {
	switch val := value.(type) {
	case map[string]interface{}:
		for dataType, dataValue := range val {
			switch sanitizeString(dataType) {
			case "S":
				return processString(dataValue.(string))
			case "N":
				return processNumber(dataValue.(string))
			case "BOOL":
				return processBoolean(dataValue.(string))
			case "NULL":
				return processNull(dataValue.(string))
			case "L":
				return processList(dataValue)
			case "M":
				return processMap(dataValue.(map[string]interface{}))
			}
		}
	}
	return nil, fmt.Errorf("unknown type")
}
