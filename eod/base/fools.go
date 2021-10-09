package base

import (
	"fmt"
	"strings"
	"time"
)

// FOOLS
var fools []string

var IsFoolsMode = time.Now().Month() == time.April && time.Now().Day() == 1

func IsFool(inp string) bool {
	for _, val := range fools {
		if strings.Contains(inp, val) {
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
