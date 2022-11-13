package polls

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Polls) pollSuccess(p *types.Poll) {
	fmt.Println("Success", p)
}
