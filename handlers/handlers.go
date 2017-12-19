// TAX
// Copyright (C) 2017 Marco Nelles, credativ GmbH
// <https://github.com/marco-mania/tax>

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"tax/db"
	"tax/models"
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

			eingabefallbrutto, found := db.FindBy(id)
			if !found {
				http.Error(w, "Not Found", http.StatusNotFound)
			}

			ausgabefallnetto := models.BerechneEinkommenNettoEur(eingabefallbrutto)

			var ergebnis models.EinkommensFall
			ergebnis.Brutto = eingabefallbrutto
			ergebnis.Netto = ausgabefallnetto

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
		}

		db.Save(fall.ID, fall)

		w.Header().Set("Location", r.URL.Path+"/"+fall.ID)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		w.Write([]byte("{}"))

	case "PUT":

		id := r.URL.Path[len("/api/v1/"):]
		_, found := db.FindBy(id)
		if !found {
			http.Error(w, "Not Found", http.StatusNotFound)
		}

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

		w.Write([]byte("{}"))

	case "DELETE":

		id := r.URL.Path[len("/api/v1/"):]
		_, found := db.FindBy(id)
		if !found {
			http.Error(w, "Not Found", http.StatusNotFound)
		}

		db.Remove(id)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("{}"))

	default:

		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

	}

}
