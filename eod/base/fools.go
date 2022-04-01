package base

import (
	"fmt"
	"strings"
)

// FOOLS
var fools []string

var IsFoolsMode = true

func IsFool(inp string) bool {
	for _, val := range fools {
		if len(val) > 0 && strings.Contains(strings.ToLower(inp), strings.ToLower(val)) {
			return true
		}
	}
	return false
}

func MakeFoolResp(val string) string {
	return fmt.Sprintf("**%s** doesn't satisfy me!", val)
}

func (b *Base) InitFools(foolsRaw string) {
	fools = strings.Split(foolsRaw, "\n")
	for i, val := range fools {
		fools[i] = strings.TrimSpace(val)
	}
}
