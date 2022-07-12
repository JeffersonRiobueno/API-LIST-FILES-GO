package config

import (
    "io/ioutil"
	"encoding/json"
    "log"
)
type DataConfig struct {
    Dfile  string `json:"file"`
    Password string `json:"password"`
}

func FilePath() string{
	jsonConfig, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatal(err)
    }
	c := DataConfig{}
    // decodificar el contenido del json sobre la estructura
    err = json.Unmarshal(jsonConfig, &c)
    if err != nil {
        log.Fatal(err)
    }
	return c.Dfile
}