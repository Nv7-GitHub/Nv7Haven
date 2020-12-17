package nv7haven

func (n *Nv7Haven) query(query string, args []interface{}, out ...interface{}) error {
	res, err := n.sql.Query(query, args...)
	if err != nil {
		return err
	}
	defer res.Close()
	res.Next()
	err = res.Scan(out...)
	if err != nil {
		return err
	}
	return nil
}
