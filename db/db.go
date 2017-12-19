// TAX
// Copyright (C) 2017 Marco Nelles, credativ GmbH
// <https://github.com/marco-mania/tax>

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package db

import (
	"encoding/json"
	"os"

	"github.com/marco-mania/tax/models"
)

var database = make(map[string]models.EinkommensFallBrutto)
var jsonfile = "db.json"

// Init ...
func Init() {
	readJSONDBFile()
}

// FindAll ...
func FindAll() []models.EinkommensFallBrutto {
	items := make([]models.EinkommensFallBrutto, 0, len(database))
	for _, v := range database {
		items = append(items, v)
	}
	return items
}

// FindBy ...
func FindBy(key string) (models.EinkommensFallBrutto, bool) {
	com, ok := database[key]
	return com, ok
}

func readJSONDBFile() {

	if _, err := os.Stat(jsonfile); os.IsNotExist(err) {
		return
	}

	f, err := os.Open(jsonfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e := json.NewDecoder(f)
	if err := e.Decode(&database); err != nil {
		panic(err)
	}

}

func updateJSONDBFile() {

	f, err := os.Create(jsonfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	e := json.NewEncoder(f)
	if err := e.Encode(database); err != nil {
		panic(err)
	}

}

// Save ...
func Save(key string, fall models.EinkommensFallBrutto) {
	database[key] = fall
	updateJSONDBFile()
}

// Remove ...
func Remove(key string) {
	delete(database, key)
	updateJSONDBFile()
}
