package main

import "testing"

func TestExtractNameValue(t *testing.T) {
	input := `{__name__="my_metric", job="test"}`
	result := ExtractNameValue(input)
	if result != "my_metric" {
		t.Errorf("expected %q, got %q", "my_metric", result)
	}
}

func TestExtractNameValueNotFound(t *testing.T) {
	input := `{job="test"}`
	result := ExtractNameValue(input)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestExtractValueByParameter(t *testing.T) {
	input := `{__name__="metric", job="myapp"}`
	result := ExtractValueByParameter(input, "job")
	if result != "myapp" {
		t.Errorf("expected %q, got %q", "myapp", result)
	}
}

func TestExtractValueByParameterRegex(t *testing.T) {
	input := `{__name__="metric", job=~".*myapp.*"}`
	result := ExtractValueByParameter(input, "job")
	if result != "myapp" {
		t.Errorf("expected %q, got %q", "myapp", result)
	}
}
