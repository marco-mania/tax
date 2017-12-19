// TAX
// Copyright (C) 2017 Marco Nelles, credativ GmbH
// <https://github.com/marco-mania/tax>

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

package helpers

import (
	"math"
)

// Round - Auf ganze Zahl runden
func Round(f float64) float64 {

	return math.Floor(f + .5)

}

// Round2 - Runden auf zwei Nachkommastellen
func Round2(f float64) float64 {

	return Round(f*100.0) / 100.0

}

// Floor2 - Alles nach der zweiten Nachkommastelle verwerfen
func Floor2(f float64) float64 {

	return math.Floor(f*100.0) / 100.0

}
