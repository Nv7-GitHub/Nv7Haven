package base

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (b *Base) TreeSize(tx *sqlx.Tx, id int, els []int, guild string) (int, bool, error) {
	var treesize int
	var loop bool
	err := tx.QueryRow(`WITH RECURSIVE parents(els, id) AS (
			VALUES($2::integer[], 1)
	 	UNION
			(SELECT b.parents els, b.id id FROM elements b INNER JOIN parents p ON b.id=ANY(p.els) where guild=$1)
	 	) SELECT COUNT(*), EXISTS(SELECT 1 FROM parents WHERE id=$3) FROM parents WHERE id>0`, guild, pq.Array(els), id).Scan(&treesize, &loop)
	if err != nil {
		return 0, false, err
	}
	return treesize, loop, nil
}
