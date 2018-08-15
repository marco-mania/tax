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

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/marco-mania/tax/db"
	"github.com/marco-mania/tax/models"
)

// TaxHandler ...
func TaxHandler(w http.ResponseWriter, r *http.Request) {

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {

	case "GET":

		var output []byte
		var err error
		id := r.URL.Path[len("/api/v1/"):]
		if len(id) > 0 {

			var ergebnis models.EinkommensFall

			eingabefallbrutto, found := db.FindBy(id)
			if !found {

				http.Error(w, "Not Found", http.StatusNotFound)

			} else {

				ausgabefallnetto := models.BerechneEinkommenNettoEur(eingabefallbrutto)
				ergebnis.Brutto = eingabefallbrutto
				ergebnis.Netto = ausgabefallnetto

			}

			output, err = json.Marshal(ergebnis)

		} else {

			eingabefaellebrutto := db.FindAll()

			var ergebnis []models.EinkommensFall
			for _, f := range eingabefaellebrutto {
				var fall models.EinkommensFall
				ausgabefallnetto := models.BerechneEinkommenNettoEur(f)
				fall.Brutto = f
				fall.Netto = ausgabefallnetto
				ergebnis = append(ergebnis, fall)
			}

			output, err = json.Marshal(ergebnis)

		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		w.Write(output)

	case "POST":

		input, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var fall models.EinkommensFallBrutto
		err = json.Unmarshal(input, &fall)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, found := db.FindBy(fall.ID)
		if found {
			http.Error(w, "Already exists", http.StatusConflict)
		} else {
			db.Save(fall.ID, fall)
			w.Header().Set("Location", r.URL.Path+"/"+fall.ID)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
		}

		w.Write([]byte("{}"))

	case "PUT":

		id := r.URL.Path[len("/api/v1/"):]
		_, found := db.FindBy(id)
		if !found {

			http.Error(w, "Not Found", http.StatusNotFound)

		} else {

			input, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			var fall models.EinkommensFallBrutto
			err = json.Unmarshal(input, &fall)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			db.Save(fall.ID, fall)

			w.Header().Set("Location", r.URL.Path+"/"+fall.ID)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)

		}

		w.Write([]byte("{}"))

	case "DELETE":

		id := r.URL.Path[len("/api/v1/"):]
		_, found := db.FindBy(id)
		if !found {

			http.Error(w, "Not Found", http.StatusNotFound)

		} else {

			db.Remove(id)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)

		}

		w.Write([]byte("{}"))

	default:

		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

		w.Write([]byte("{}"))

	}

}
