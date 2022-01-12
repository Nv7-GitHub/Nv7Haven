package eodb

func (b *DB) Optimize() error {
	b.Lock()
	defer b.Unlock()

	// Delete existing data
	_, err := b.elemFile.Seek(0, 0)
	if err != nil {
		return err
	}
	err = b.elemFile.Truncate(0)
	if err != nil {
		return err
	}

	// Rewrite elements
	for _, elem := range b.Elements {
		dat, err := json.Marshal(elem)
		if err != nil {
			return err
		}
		_, err = b.elemFile.WriteString(string(dat) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
