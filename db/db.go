/* TAX - Simple german tax calculator web app
 * Copyright (C) 2017-2018 Marco Nelles, credativ GmbH
 * <https://github.com/marco-mania/tax>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package db

import (
	"encoding/json"
	"os"
	"path"

	"github.com/marco-mania/tax/configdir"
	"github.com/marco-mania/tax/helpers"
	"github.com/marco-mania/tax/models"
)

var database = make(map[string]models.EinkommensFallBrutto)
var jsonfile = configdir.ConfigDir + "/tax/db.json"

func init() {
	configDirTax := path.Dir(jsonfile)
	exists, _ := helpers.Exists(configDirTax)
	if !exists {
		os.MkdirAll(configDirTax, 0755)
	}
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
