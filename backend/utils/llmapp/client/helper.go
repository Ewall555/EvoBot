package client

import (
	"fmt"
)

func loadParamFromVars(key string, vars map[string]interface{}) string {
	if _, ok := vars[key]; !ok {
		return ""
	}

	return fmt.Sprintf("%v", vars[key])
}
