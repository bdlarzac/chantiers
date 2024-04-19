/*
Basic equivalent to php function print_r().
To have an output more readable than fmt.Printf("%+v\n", myVar).

@copyright  Thierry Graff
@license    GPL - conforms to file LICENCE located in root directory of current repository.
@history    2024-04-04 12:17:55+02:00, Thierry Graff : Creation
*/
package tiglib

import (
	"github.com/kortschak/utter"
)

func Print_r(x interface{}) {
	utter.Config.SortKeys = true
	utter.Config.Indent = "    "
	utter.Config.ElideType = true
	utter.Config.Dump(x)
}
