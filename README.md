# TAX

Simple german tax calculator app with focus at the costs of a company car.

Features:

* Calculate "Steuern" and "Sozialabgaben"
* See the difference between your income with and without a company car

## Included

* JQuery - <https://jquery.com/>
* popper.js - <https://popper.js.org/>
* Bootstrap - <https://getbootstrap.com/>
* Knockout - <http://knockoutjs.com/>

## Build & Run on Linux

In the tax folder just execute

```Shell
packr
GOOS=linux go build -ldflags="-s -w" -o tax-linux-amd64
./tax-linux-amd64
```

Open <http://localhost:8001/>

## Licence

TAX itself is licensed under the GPLv3 (<https://www.gnu.org/licenses/gpl-3.0.de.html>).

All bundled components are Free/Open-Source software with a known and approved open source license.
