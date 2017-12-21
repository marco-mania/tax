/* TAX - Simple german tax calculator web app
 * Copyright (C) 2017 Marco Nelles, credativ GmbH
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

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/marco-mania/tax/handlers"
)

func main() {

	box := packr.NewBox("./app/web")

	http.Handle("/", http.FileServer(box))

	http.HandleFunc("/api/v1/", handlers.TaxHandler)

	fmt.Println("Service TAX v2 (API v1) starting and listen on http://localhost:8001/ ...")

	err := http.ListenAndServe("localhost:8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
