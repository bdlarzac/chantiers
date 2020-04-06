/******************************************************************************
    Renvoie la liste des clés d'une map.
    Les clés doivent être des strings

    PAS UTILISE parcequ'on ne peut pas passer n'importe quoi à la place de interface{}

    @license    GPL
    @history    2020-01-16 15:19:29+01:00, Thierry Graff : Creation
********************************************************************************/

package tiglib

func Keys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	return keys
}
