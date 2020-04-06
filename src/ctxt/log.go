/******************************************************************************

    @license    GPL
    @history    2019-12-13 15:06:25+01:00, Thierry Graff : Creation
********************************************************************************/
package ctxt

import (
	"bdl.local/bdl/generic/wilk/werr"
)

func LogError(err error) {
	werr.Print(err)
}
