package parser

import (
	"encoding/json"
	"io/ioutil"
)

type secrets map[string]string

var (
	s secrets
)

func loadSecrets() (secrets, error) {
	file, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
