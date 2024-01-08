package main

import (
	"reflect"
	"testing"
	"time"
)

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  test  ", "test"},
		{"\ntest\n", "test"},
		{"\ttest\t", "test"},
		{"test", "test"},
		{"", ""},
	}

	for _, test := range tests {
		if result := sanitizeString(test.input); result != test.expected {
			t.Errorf("sanitizeString(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestProcessString(t *testing.T) {
	validTime := time.Now().Format(time.RFC3339)
	tests := []struct {
		input    string
		expected interface{}
		err      bool
	}{
		{"  test  ", "test", false},
		{"", nil, true},
		{validTime, time.Now().Unix(), false},
	}

	for _, test := range tests {
		result, err := processString(test.input)
		if (err != nil) != test.err {
			t.Errorf("processString(%q) error = %v, wantErr %v", test.input, err, test.err)
		}
		if !test.err && result != test.expected {
			t.Errorf("processString(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestProcessNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		err      bool
	}{
		{"  42  ", 42.0, false},
		{"-13.5", -13.5, false},
		{"invalid", nil, true},
		{"", nil, true},
	}

	for _, test := range tests {
		result, err := processNumber(test.input)
		if (err != nil) != test.err {
			t.Errorf("processNumber(%q) error = %v, wantErr %v", test.input, err, test.err)
		}
		if !test.err && result != test.expected {
			t.Errorf("processNumber(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestProcessBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		err      bool
	}{
		{"  true  ", true, false},
		{"f", false, false},
		{"invalid", nil, true},
		{"", nil, true},
	}

	for _, test := range tests {
		result, err := processBoolean(test.input)
		if (err != nil) != test.err {
			t.Errorf("processBoolean(%q) error = %v, wantErr %v", test.input, err, test.err)
		}
		if !test.err && result != test.expected {
			t.Errorf("processBoolean(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestProcessNull(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
		err      bool
	}{
		{"  true  ", nil, false},
		{"false", nil, true},
		{"false", nil, true},
		{"invalid", nil, true},
		{"", nil, true},
	}

	for _, test := range tests {
		result, err := processNull(test.input)
		if (err != nil) != test.err {
			t.Errorf("processNull(%q) error = %v, wantErr %v", test.input, err, test.err)
		}
		if !test.err && result != test.expected {
			t.Errorf("processNull(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestProcessList(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []interface{}
		err      bool
	}{
		{[]interface{}{map[string]interface{}{"N": "42"}}, []interface{}{42.0}, false},
		{[]interface{}{"test", "42"}, nil, true},
		{"not a list", nil, true},
	}

	for _, test := range tests {
		result, err := processList(test.input)
		if (err != nil) != test.err {
			t.Errorf("processList(%v) error = %v, wantErr %v", test.input, err, test.err)
		}
		if !test.err && !reflect.DeepEqual(result, test.expected) {
			t.Errorf("processList(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestProcessMap(t *testing.T) {
	tests := []struct {
		input    map[string]interface{}
		expected map[string]interface{}
		err      bool
	}{
		{map[string]interface{}{"string_1": map[string]interface{}{"S": "123"}}, map[string]interface{}{"string_1": "123"}, false},
		{map[string]interface{}{}, nil, true},
	}

	for _, test := range tests {
		result, err := processMap(test.input)
		if (err != nil) != test.err {
			t.Errorf("processMap(%v) error = %v, wantErr %v", test.input, err, test.err)
		}
		if !test.err && !reflect.DeepEqual(result, test.expected) {
			t.Errorf("processMap(%v) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestProcessField(t *testing.T) {
	tests := []struct {
		key      string
		value    interface{}
		expected interface{}
		err      bool
	}{
		{"key", "value", "value", true},
		{"", "value", nil, true},
		{"key", "", nil, true},
	}

	for _, test := range tests {
		result, err := processField(test.key, test.value)
		if (err != nil) != test.err {
			t.Errorf("processField(%q, %v) error = %v, wantErr %v", test.key, test.value, err, test.err)
		}
		if !test.err && result != test.expected {
			t.Errorf("processField(%q, %v) = %v, want %v", test.key, test.value, result, test.expected)
		}
	}
}

func TestProcessValue(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected interface{}
		err      bool
	}{
		{map[string]interface{}{"S": "string"}, "string", false},
		{map[string]interface{}{"N": "42"}, 42.0, false},
		{map[string]interface{}{"BOOL": "true"}, true, false},
		{map[string]interface{}{"NULL": "true"}, nil, false},
		{map[string]interface{}{"L": []interface{}{"string", "42"}}, []interface{}{}, true},
		{map[string]interface{}{"M": map[string]interface{}{"key": "value"}}, map[string]interface{}{}, true},
		{map[string]interface{}{"unknown": "value"}, nil, true},
	}

	for _, test := range tests {
		result, err := processValue(test.value)
		if (err != nil) != test.err {
			t.Errorf("processValue(%v) error = %v, wantErr %v", test.value, err, test.err)
		}
		if !test.err && !reflect.DeepEqual(result, test.expected) {
			t.Errorf("processValue(%v) = %v, want %v", test.value, result, test.expected)
		}
	}
}
