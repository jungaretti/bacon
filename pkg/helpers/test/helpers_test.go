package helpers

import (
	"bacon/pkg/helpers"
	"bufio"
	"fmt"
	"os"
	"testing"
)

type canto struct {
	Name string `yaml:"name"`
}

type profile struct {
	Name     string  `yaml:"name"`
	Age      int     `yaml:"age"`
	Cantiche []canto `yaml:"cantiche"`
}

func TestReadAndParseYamlFile(t *testing.T) {
	temp, err := os.CreateTemp("", "tmp")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	writer := bufio.NewWriter(temp)
	fmt.Fprintln(writer, "name: Dante Aligheri")
	fmt.Fprintln(writer, "age: 25")
	fmt.Fprintln(writer, "cantiche:")
	fmt.Fprintln(writer, "- name: inferno")
	fmt.Fprintln(writer, "- name: purgatorio")
	fmt.Fprintln(writer, "- name: paradiso")
	err = temp.Close()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	out := profile{}
	err = helpers.ReadAndParseYamlFile(temp.Name(), &out)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestWriteYamlFile(t *testing.T) {
	in := profile{
		Name: "Bruno Latini",
		Age:  42,
		Cantiche: []canto{
			{
				Name: "purgatorio",
			},
		},
	}

	temp, err := os.CreateTemp("", "tmp")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = helpers.WriteYamlFile(temp.Name(), &in)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestPostJson(t *testing.T) {
	body := profile{
		Name: "Virgilio",
		Age:  36,
		Cantiche: []canto{
			{
				Name: "inferno",
			},
			{
				Name: "purgatorio",
			},
		},
	}

	resp, err := helpers.PostJson("https://ptsv2.com/t/tg7o8-1651468588/post", body)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if resp.StatusCode != 200 {
		t.Log(resp.StatusCode)
		t.FailNow()
	}
}
