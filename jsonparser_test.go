package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseInputJSON_ValidFile(t *testing.T) {
	expected := &InputJSON{Data: map[string]interface{}{"key": "value"}}
	content := []byte(`{"key": "value"}`)
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := parseInputJSON(tmpfile.Name())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestParseInputJSON_InvalidFile(t *testing.T) {
	content := []byte(`{"key": "value"`) // missing closing bracket
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = parseInputJSON(tmpfile.Name())

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestParseInputJSON_NonExistentFile(t *testing.T) {
	_, err := parseInputJSON("nonexistent.json")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
