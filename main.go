// TAX
// Copyright (C) 2017 Marco Nelles, credativ GmbH
// <https://github.com/marco-mania/tax>

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"net/http"
	"tax/db"
	"tax/handlers"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./app/web")))

	db.Init()

	http.HandleFunc("/api/v1/", handlers.TaxHandler)

	http.ListenAndServe("localhost:8001", nil)

}
