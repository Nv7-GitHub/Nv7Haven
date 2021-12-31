package ai

import "github.com/Nv7-Github/Nv7Haven/eod/types"

func (a *AI) PredictCombo() []int {
	elem := a.Starters.Predict()
	combo := []int{elem}
	a.lock.RLock()
	for {
		if len(combo) == types.MaxComboLength {
			break
		}
		next := a.Links[elem].Predict()
		if next == -1 {
			// If length is too short, get more
			if len(combo) == 1 {
				next = a.Starters.Predict()
			} else {
				break
			}
		}
		combo = append(combo, next)
		elem = next
	}
	a.lock.RUnlock()

	return combo
}
