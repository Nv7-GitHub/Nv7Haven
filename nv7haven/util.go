package nv7haven

import "fmt"

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

func FormatByteSize(b int) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
