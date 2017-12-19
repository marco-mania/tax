// TAX
// Copyright (C) 2017 Marco Nelles, credativ GmbH
// <https://github.com/marco-mania/tax>

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marco-mania/tax/db"
	"github.com/marco-mania/tax/handlers"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./app/web")))

	db.Init()

	http.HandleFunc("/api/v1/", handlers.TaxHandler)

	fmt.Println("Service starting: Listen on http://localhost:8001/ ...")

	err := http.ListenAndServe("localhost:8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
