package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ContentDir   string      `yaml:"contentDir"`
	SnapshotsDir string      `yaml:"snapshotsDir"`
	Retentions   []Retention `yaml:"retentions"`
}

type Retention struct {
	EveryNDays      int `yaml:"every"`
	SnapshotsToKeep int `yaml:"keep"`
}

func getConfig(filePath string) *Config {
	bb, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	var c Config
	err = yaml.Unmarshal(bb, &c)
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	return &c
}
