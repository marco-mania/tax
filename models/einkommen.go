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

package models

import (
	"encoding/json"
	"math"

	"github.com/marco-mania/tax/helpers"
)

var parameters = `[
    {
        "Jahr": 2016,
        "KrankenversicherungArbeitnehmeranteil": 0.073,
        "KrankenversicherungArbeitgeberanteil": 0.073,
        "KrankenversicherungArbeitnehmeranteilErmaessigt": 0.07,
        "KrankenversicherungArbeitgeberanteilErmaessigt": 0.07,
        "RentenversicherungArbeitnehmeranteil": 0.0935,
        "RentenversicherungArbeitgeberanteil": 0.0935,
        "KorrekturfaktorTeilVorsorgepauschaleRentenversicherung": 0.64,
        "ArbeitslosenversicherungArbeitnehmeranteil": 0.015,
        "ArbeitslosenversicherungArbeitgeberanteil": 0.015,
        "PflegeversicherungArbeitnehmeranteil": 0.01175,
        "PflegeversicherungArbeitgeberanteil": 0.01175,
        "PflegeversicherungZuschlagFuerKinderlose": 0.0025,
        "BeitragsbemessungsgrenzeKrankenPflegeversicherung": 4275.0,
        "BeitragsbemessungsgrenzeArbeitslosenRentenversicherung": 6200.0,
        "z2": 8652.0,
        "z2a": 0.0000099362,
        "z2b": 0.14,
        "z2c": 0.0,
        "z3": 13669.0,
        "z3a": 0.000002254,
        "z3b": 0.2397,
        "z3c": 952.48,
        "z4": 53665.0,
        "z4a": 0.0,
        "z4b": 0.42,
        "z4c": -8394.14,
        "z5": 254446.0,
        "z5a": 0.0,
        "z5b": 0.45,
        "z5c": -16027.52,
        "Solidaritaetszuschlag": 0.055,
        "Kirchensteuer": 0.09
    },
    {
        "Jahr": 2017,
        "KrankenversicherungArbeitnehmeranteil": 0.073,
        "KrankenversicherungArbeitgeberanteil": 0.073,
        "KrankenversicherungArbeitnehmeranteilErmaessigt": 0.07,
        "KrankenversicherungArbeitgeberanteilErmaessigt": 0.07,
        "RentenversicherungArbeitnehmeranteil": 0.0935,
        "RentenversicherungArbeitgeberanteil": 0.0935,
        "KorrekturfaktorTeilVorsorgepauschaleRentenversicherung": 0.68,
        "ArbeitslosenversicherungArbeitnehmeranteil": 0.015,
        "ArbeitslosenversicherungArbeitgeberanteil": 0.015,
        "PflegeversicherungArbeitnehmeranteil": 0.01275,
        "PflegeversicherungArbeitgeberanteil": 0.01275,
        "PflegeversicherungZuschlagFuerKinderlose": 0.0025,
        "BeitragsbemessungsgrenzeKrankenPflegeversicherung": 4350.0,
        "BeitragsbemessungsgrenzeArbeitslosenRentenversicherung": 6350.0,
        "z2": 8820.0,
        "z2a": 0.0000100727,
        "z2b": 0.14,
        "z2c": 0.0,
        "z3": 13769.0,
        "z3a": 0.0000022376,
        "z3b": 0.2397,
        "z3c": 939.57,
        "z4": 54057.0,
        "z4a": 0.0,
        "z4b": 0.42,
        "z4c": -8475.44,
        "z5": 256303.0,
        "z5a": 0.0,
        "z5b": 0.45,
        "z5c": -16164.53,
        "Solidaritaetszuschlag": 0.055,
        "Kirchensteuer": 0.09
    }
]`

// EinkommensFallBrutto fasst die Eingabe-Parameter zusammen
type EinkommensFallBrutto struct {
	ID                               string  `json:"ID"`
	Jahr                             int     `json:"Jahr"`
	Steuerklasse                     int     `json:"Steuerklasse"` // 1 oder 4
	EinkommenBruttoMonEur            float64 `json:"EinkommenBruttoMonEur"`
	FirmenwagenBruttolistenpreisEur  int     `json:"FirmenwagenBruttolistenpreisEur"` // wird abgerundet auf "ganze 100"
	EntfernungWohnortArbeitsortKM    int     `json:"EntfernungWohnortArbeitsortKM"`   // wird abgerundet auf ganze KM
	CarAllowanceEur                  float64 `json:"CarAllowanceEur"`
	ZusatzbeitragKrankenkasseProzent float64 `json:"ZusatzbeitragKrankenkasseProzent"`
	Kinder                           bool    `json:"Kinder"`
	Kirche                           bool    `json:"Kirche"`
}

// EinkommensFallNetto fasst die Ausgabe-Parameter zusammen
type EinkommensFallNetto struct {
	EinkommenNettoMonEur                float64 `json:"EinkommenNettoMonEur"`
	EinkommenNettoOhneFirmenwagenMonEur float64 `json:"EinkommenNettoOhneFirmenwagenMonEur"`
	EinkommenDifferenzEur               float64 `json:"EinkommenDifferenzEur"` // Differenz Netto mit/ohne Firmenwagen
	GeldwerterVorteilFirmenwagenMonEur  float64 `json:"GeldwerterVorteilFirmenwagenMonEur"`
	KVBeitragMonEur                     float64 `json:"KVBeitragMonEur"`
	RVBeitragMonEur                     float64 `json:"RVBeitragMonEur"`
	AVBeitragMonEur                     float64 `json:"AVBeitragMonEur"`
	PVBeitragMonEur                     float64 `json:"PVBeitragMonEur"`
	LohnsteuerMonEur                    float64 `json:"LohnsteuerMonEur"`
	KirchensteuerMonEur                 float64 `json:"KirchensteuerMonEur"`
	SolZuschlagMonEur                   float64 `json:"SolZuschlagMonEur"`
}

// EinkommensFall definiert den vollständigen Einkommensfall
type EinkommensFall struct {
	Brutto EinkommensFallBrutto
	Netto  EinkommensFallNetto
}

// Parameter sind alle Sozialabgaben- und Steuer-Parameter für das lfd. Jahr
type Parameter struct {
	Jahr int `json:"jahr"`

	KrankenversicherungArbeitnehmeranteil                  float64 `json:"KrankenversicherungArbeitnehmeranteil"`
	KrankenversicherungArbeitgeberanteil                   float64 `json:"KrankenversicherungArbeitgeberanteil"`
	KrankenversicherungArbeitnehmeranteilErmaessigt        float64 `json:"KrankenversicherungArbeitnehmeranteilErmaessigt"`
	KrankenversicherungArbeitgeberanteilErmaessigt         float64 `json:"KrankenversicherungArbeitgeberanteilErmaessigt"`
	RentenversicherungArbeitnehmeranteil                   float64 `json:"RentenversicherungArbeitnehmeranteil"`
	RentenversicherungArbeitgeberanteil                    float64 `json:"RentenversicherungArbeitgeberanteil"`
	KorrekturfaktorTeilVorsorgepauschaleRentenversicherung float64 `json:"KorrekturfaktorTeilVorsorgepauschaleRentenversicherung"`
	ArbeitslosenversicherungArbeitnehmeranteil             float64 `json:"ArbeitslosenversicherungArbeitnehmeranteil"`
	ArbeitslosenversicherungArbeitgeberanteil              float64 `json:"ArbeitslosenversicherungArbeitgeberanteil"`
	PflegeversicherungArbeitnehmeranteil                   float64 `json:"PflegeversicherungArbeitnehmeranteil"`
	PflegeversicherungArbeitgeberanteil                    float64 `json:"PflegeversicherungArbeitgeberanteil"`
	PflegeversicherungZuschlagFuerKinderlose               float64 `json:"PflegeversicherungZuschlagFuerKinderlose"`
	BeitragsbemessungsgrenzeKrankenPflegeversicherung      float64 `json:"BeitragsbemessungsgrenzeKrankenPflegeversicherung"`
	BeitragsbemessungsgrenzeArbeitslosenRentenversicherung float64 `json:"BeitragsbemessungsgrenzeArbeitslosenRentenversicherung"`

	Z2                    float64 `json:"z2"`
	Z2a                   float64 `json:"z2a"`
	Z2b                   float64 `json:"z2b"`
	Z2c                   float64 `json:"z2c"`
	Z3                    float64 `json:"z3"`
	Z3a                   float64 `json:"z3a"`
	Z3b                   float64 `json:"z3b"`
	Z3c                   float64 `json:"z3c"`
	Z4                    float64 `json:"z4"`
	Z4a                   float64 `json:"z4a"`
	Z4b                   float64 `json:"z4b"`
	Z4c                   float64 `json:"z4c"`
	Z5                    float64 `json:"z5"`
	Z5a                   float64 `json:"z5a"`
	Z5b                   float64 `json:"z5b"`
	Z5c                   float64 `json:"z5c"`
	Solidaritaetszuschlag float64 `json:"Solidaritaetszuschlag"`
	Kirchensteuer         float64 `json:"Kirchensteuer"`
}

func steuerFunktion(a float64, b float64, c float64, _b float64) float64 {

	return math.Floor(_b*(a*_b+b) + c)

}

func berechne(fall EinkommensFallBrutto, param Parameter) EinkommensFallNetto {

	var ergebnis EinkommensFallNetto

	// Gesamt-Brutto:
	var einkommenGesamtBruttoMonEur = fall.EinkommenBruttoMonEur
	//var einkommenGesamtBruttoJhrEur = fall.einkommenBruttoMonEur * 12.0

	// Steuer-Brutto:
	var einkommenSteuerBruttoMonEur = fall.EinkommenBruttoMonEur
	var einkommenSteuerBruttoJhrEur = einkommenSteuerBruttoMonEur * 12.0

	// SV-Brutto:
	var einkommenSozialversicherungBruttoMonEur = fall.EinkommenBruttoMonEur
	//var einkommenSozialversicherungBruttoJhrEur = einkommenSozialversicherungBruttoMonEur * 12.0

	// Falls ein Firmenwagen auch zur privaten Nutzung überlassen, erhöht sich das
	// steuer- und sozialabgabepflichtige Einkommen
	var geldwerterVorteilFirmenwagenMonEur float64
	if fall.FirmenwagenBruttolistenpreisEur > 0 {

		var firmenwagenNeupreisAbgerundetEur = float64(fall.FirmenwagenBruttolistenpreisEur - (fall.FirmenwagenBruttolistenpreisEur % 100)) // abgerundet auf "ganze 100"

		var a = firmenwagenNeupreisAbgerundetEur * 0.01                                                 // Lohnart 873: Privatfahrten, 1%
		var b = firmenwagenNeupreisAbgerundetEur * 0.0003 * float64(fall.EntfernungWohnortArbeitsortKM) // 0,03%
		var pendlerpauschbetragMonEur = 15.0 * 0.3 * float64(fall.EntfernungWohnortArbeitsortKM)        // Lohnart 875: Whg/Arbeit p. St.
		geldwerterVorteilFirmenwagenMonEur = a + b

		ergebnis.GeldwerterVorteilFirmenwagenMonEur = geldwerterVorteilFirmenwagenMonEur

		einkommenGesamtBruttoMonEur += geldwerterVorteilFirmenwagenMonEur
		//einkommenGesamtBruttoJhrEur = einkommenGesamtBruttoMonEur * 12.0

		einkommenSteuerBruttoMonEur += (geldwerterVorteilFirmenwagenMonEur - pendlerpauschbetragMonEur)
		einkommenSteuerBruttoJhrEur = einkommenSteuerBruttoMonEur * 12.0

		einkommenSozialversicherungBruttoMonEur += geldwerterVorteilFirmenwagenMonEur - pendlerpauschbetragMonEur

	}

	// Bruttobeträge für Sozialversicherungsbeiträge berechnen. Beitragsbemessungsgrenzen überschritten? Falls ja, dann deckeln.

	var einkommenKVBruttoMonEur = einkommenSozialversicherungBruttoMonEur
	if einkommenSozialversicherungBruttoMonEur > param.BeitragsbemessungsgrenzeKrankenPflegeversicherung {
		einkommenKVBruttoMonEur = param.BeitragsbemessungsgrenzeKrankenPflegeversicherung
	}
	var einkommenKVBruttoJhrEur = einkommenKVBruttoMonEur * 12.0

	var einkommenPVBruttoMonEur = einkommenSozialversicherungBruttoMonEur
	if einkommenSozialversicherungBruttoMonEur > param.BeitragsbemessungsgrenzeKrankenPflegeversicherung {
		einkommenPVBruttoMonEur = param.BeitragsbemessungsgrenzeKrankenPflegeversicherung
	}
	//var einkommenPVBruttoJhrEur = einkommenPVBruttoMonEur * 12.0

	var einkommenAVBruttoMonEur = einkommenSozialversicherungBruttoMonEur
	if einkommenSozialversicherungBruttoMonEur > param.BeitragsbemessungsgrenzeArbeitslosenRentenversicherung {
		einkommenAVBruttoMonEur = param.BeitragsbemessungsgrenzeArbeitslosenRentenversicherung
	}
	//var einkommenAVBruttoJhrEur = einkommenAVBruttoMonEur * 12.0

	var einkommenRVBruttoMonEur = einkommenSozialversicherungBruttoMonEur
	if einkommenSozialversicherungBruttoMonEur > param.BeitragsbemessungsgrenzeArbeitslosenRentenversicherung {
		einkommenRVBruttoMonEur = param.BeitragsbemessungsgrenzeArbeitslosenRentenversicherung
	}
	var einkommenRVBruttoJhrEur = einkommenRVBruttoMonEur * 12.0

	// Vorsorgepauschbetrag berechnen
	var teilVorsorgepauschaleRentenversicherung = einkommenRVBruttoJhrEur * param.RentenversicherungArbeitnehmeranteil * param.KorrekturfaktorTeilVorsorgepauschaleRentenversicherung

	var teilVorsorgepauschaleKrankenPflegeversicherungA = einkommenKVBruttoJhrEur // Mindestsatz 12% des Arbeitslohns
	if teilVorsorgepauschaleKrankenPflegeversicherungA > 1900.0 {
		teilVorsorgepauschaleKrankenPflegeversicherungA = 1900.0
	}

	var teilVorsorgepauschaleKrankenPflegeversicherungB float64
	if fall.Kinder {
		teilVorsorgepauschaleKrankenPflegeversicherungB = einkommenKVBruttoJhrEur * (param.KrankenversicherungArbeitnehmeranteilErmaessigt + (fall.ZusatzbeitragKrankenkasseProzent / 100.0) + param.PflegeversicherungArbeitnehmeranteil)
	} else {
		teilVorsorgepauschaleKrankenPflegeversicherungB = einkommenKVBruttoJhrEur * (param.KrankenversicherungArbeitnehmeranteilErmaessigt + (fall.ZusatzbeitragKrankenkasseProzent / 100.0) + param.PflegeversicherungArbeitnehmeranteil + param.PflegeversicherungZuschlagFuerKinderlose)
	}

	var vorsorgePauschbetrag float64

	if (teilVorsorgepauschaleRentenversicherung + teilVorsorgepauschaleKrankenPflegeversicherungB) > (teilVorsorgepauschaleRentenversicherung + teilVorsorgepauschaleKrankenPflegeversicherungA) {
		vorsorgePauschbetrag = math.Ceil(teilVorsorgepauschaleRentenversicherung + teilVorsorgepauschaleKrankenPflegeversicherungB)
	} else {
		vorsorgePauschbetrag = math.Ceil(teilVorsorgepauschaleRentenversicherung + teilVorsorgepauschaleKrankenPflegeversicherungA)
	}

	// Sozialabgaben berechnen: Bei Sozialabgaben wird immer auf zwei Nachkommastellen *gerundet*

	// Krankenkassen können seit 2015 einen Zusatzbeitrag erheben
	var krankenversicherungMon = helpers.Round2(einkommenKVBruttoMonEur * (param.KrankenversicherungArbeitnehmeranteil + (fall.ZusatzbeitragKrankenkasseProzent / 100.0)))
	ergebnis.KVBeitragMonEur = krankenversicherungMon

	var pflegeversicherungMon float64
	if fall.Kinder {
		pflegeversicherungMon = helpers.Round2(einkommenPVBruttoMonEur * param.PflegeversicherungArbeitnehmeranteil)
	} else {
		pflegeversicherungMon = helpers.Round2(einkommenPVBruttoMonEur * (param.PflegeversicherungArbeitnehmeranteil + param.PflegeversicherungZuschlagFuerKinderlose))
	}
	ergebnis.PVBeitragMonEur = pflegeversicherungMon

	var rentenversicherungMon = helpers.Round2(einkommenRVBruttoMonEur * param.RentenversicherungArbeitnehmeranteil)
	ergebnis.RVBeitragMonEur = rentenversicherungMon

	var arbeitslosenversicherungMon = helpers.Round2(einkommenAVBruttoMonEur * param.ArbeitslosenversicherungArbeitnehmeranteil)
	ergebnis.AVBeitragMonEur = arbeitslosenversicherungMon

	var werbungskostenPauschbetrag = 1000.0
	var sonderausgabenPauschbetrag = 36.0

	// zu versteuerndes Einkommen nach Abzug der Pauschalen
	var _einkommenVorsteuerSteuerJhrEur = math.Floor(einkommenSteuerBruttoJhrEur - (vorsorgePauschbetrag + werbungskostenPauschbetrag + sonderausgabenPauschbetrag))

	// Lohnsteuer

	// Maximale Steuerbeträge der Zonen
	//var Z2_m float64 = steuerFunktion(param.Z2a, param.Z2b, param.Z2c, param.Z3-param.Z2)
	//var Z3_m float64 = Z2_m + steuerFunktion(param.Z3a, param.Z3b, param.Z3c, param.Z4-param.Z3)
	//var Z4_m float64 = steuerFunktion(param.Z4a, param.Z4b, param.Z4c, param.Z5)

	var einkommensteuerJahr float64

	if _einkommenVorsteuerSteuerJhrEur > param.Z5 {

		einkommensteuerJahr = steuerFunktion(param.Z5a, param.Z5b, param.Z5c, _einkommenVorsteuerSteuerJhrEur)

	} else if _einkommenVorsteuerSteuerJhrEur > param.Z4 {

		einkommensteuerJahr = steuerFunktion(param.Z4a, param.Z4b, param.Z4c, _einkommenVorsteuerSteuerJhrEur)

	} else if _einkommenVorsteuerSteuerJhrEur > param.Z3 {

		einkommensteuerJahr = steuerFunktion(param.Z3a, param.Z3b, param.Z3c, _einkommenVorsteuerSteuerJhrEur-param.Z3)

	} else if _einkommenVorsteuerSteuerJhrEur > param.Z2 {

		einkommensteuerJahr = steuerFunktion(param.Z2a, param.Z2b, param.Z2c, _einkommenVorsteuerSteuerJhrEur-param.Z2)

	} else {

		einkommensteuerJahr = _einkommenVorsteuerSteuerJhrEur

	}

	// Bei Steuern wird immer auf zwei Nachkommastellen *ab*gerundet
	var einkommensteuerMon = helpers.Floor2(einkommensteuerJahr / 12.0)
	ergebnis.LohnsteuerMonEur = einkommensteuerMon

	var solidaritaetszuschlagJahr = einkommensteuerJahr * param.Solidaritaetszuschlag
	var solidaritaetszuschlagMon = helpers.Floor2(solidaritaetszuschlagJahr / 12.0)
	ergebnis.SolZuschlagMonEur = solidaritaetszuschlagMon

	var _kirchensteuer = param.Kirchensteuer
	if !fall.Kirche {
		_kirchensteuer = 0.0
	}
	var kirchensteuerJahr = einkommensteuerJahr * _kirchensteuer
	var kirchensteuerMon = helpers.Floor2(kirchensteuerJahr / 12.0)
	ergebnis.KirchensteuerMonEur = kirchensteuerMon

	ergebnis.EinkommenNettoMonEur = helpers.Round2(einkommenGesamtBruttoMonEur - (einkommensteuerMon + solidaritaetszuschlagMon + kirchensteuerMon + krankenversicherungMon + rentenversicherungMon + arbeitslosenversicherungMon + pflegeversicherungMon + geldwerterVorteilFirmenwagenMonEur))

	return ergebnis

}

func getParameters() []Parameter {

	var p []Parameter
	json.Unmarshal([]byte(parameters), &p)
	return p

}

// BerechneEinkommenNettoEur ...
func BerechneEinkommenNettoEur(eingabe EinkommensFallBrutto) EinkommensFallNetto {

	indexJahr := -1
	parameters := getParameters()
	for i, param := range parameters {
		if param.Jahr == eingabe.Jahr {
			indexJahr = i
		}
	}

	var ausgabe EinkommensFallNetto

	if indexJahr != -1 {

		ausgabe = berechne(eingabe, parameters[indexJahr])

		eingabeOhneFirmenwagen := eingabe
		eingabeOhneFirmenwagen.FirmenwagenBruttolistenpreisEur = 0
		eingabeOhneFirmenwagen.EinkommenBruttoMonEur = eingabe.EinkommenBruttoMonEur + eingabe.CarAllowanceEur

		ausgabeOhneFirmenwagen := berechne(eingabeOhneFirmenwagen, parameters[indexJahr])

		ausgabe.EinkommenNettoOhneFirmenwagenMonEur = ausgabeOhneFirmenwagen.EinkommenNettoMonEur
		ausgabe.EinkommenDifferenzEur = ausgabe.EinkommenNettoOhneFirmenwagenMonEur - ausgabe.EinkommenNettoMonEur

	}

	return ausgabe

}
