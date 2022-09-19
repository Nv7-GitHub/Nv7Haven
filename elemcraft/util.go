package elemcraft

func StripRecipe(recipe [][]int) [][]int {
	newR := recipe
	// Forwards
	for _, row := range recipe {
		allZeros := true
		for _, i := range row {
			if i != -1 {
				allZeros = false
				break
			}
		}
		if allZeros {
			newR = newR[1:]
		} else {
			break
		}
	}
	recipe = newR

	// Backwards
	newR = recipe
	for i := len(recipe) - 1; i >= 0; i-- {
		allZeros := true
		for _, i := range recipe[i] {
			if i != -1 {
				allZeros = false
				break
			}
		}
		if allZeros {
			newR = newR[:len(newR)-1]
		} else {
			break
		}
	}

	// TODO: Strip cols

	return newR
}
