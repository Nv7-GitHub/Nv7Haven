package eodb

func (d *DB) BeginTransaction() {
	d.inTransaction = true
}

func (d *DB) CommitTransaction() error {
	err := d.Optimize()
	if err != nil {
		return err
	}
	d.inTransaction = false
	return nil
}
