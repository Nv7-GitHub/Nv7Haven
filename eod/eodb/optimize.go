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

func (b *DB) OptimizeCats() error {
	b.Lock()
	defer b.Unlock()

	for name, f := range b.catCacheFiles {
		// Delete existing data
		_, err := f.Seek(0, 0)
		if err != nil {
			return err
		}
		err = f.Truncate(0)
		if err != nil {
			return err
		}

		// Create entry
		dat := b.catCache[name]
		els := make([]int, len(dat))
		i := 0
		for k := range dat {
			els[i] = k
			i++
		}
		entry := catCacheEntry{
			Op:   catCacheOpAdd,
			Data: els,
		}
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		// Rewrite elements
		_, err = f.WriteString(string(data) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *DB) OptimizeInvs() error {
	b.Lock()
	defer b.Unlock()

	for name, f := range b.invDataFiles {
		// Delete existing data
		_, err := f.Seek(0, 0)
		if err != nil {
			return err
		}
		err = f.Truncate(0)
		if err != nil {
			return err
		}

		// Create entry
		dat := b.invData[name]
		els := make([]int, len(dat))
		i := 0
		for k := range dat {
			els[i] = k
			i++
		}
		entry := invOp{
			Kind: invOpAdd,
			Data: els,
		}
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		// Rewrite elements
		_, err = f.WriteString(string(data) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
