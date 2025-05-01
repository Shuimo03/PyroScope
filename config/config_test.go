package config

import (
	"fmt"
	"os"
	"testing"
)

func TestYAML(t *testing.T) {
	yaml, err := os.ReadFile("../example/yaml/node-exporter.yaml")
	if err != nil {
		t.Fatal(err)
	}
	cnf, err := Load(string(yaml))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cnf)
}
