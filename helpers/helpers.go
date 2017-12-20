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
