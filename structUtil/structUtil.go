package structUtil

import "github.com/fatih/structs"

func ToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}
