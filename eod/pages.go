package eod

const ldbQuery = `
SELECT rw
FROM (
    SELECT 
         ROW_NUMBER() OVER (ORDER BY ` + "count" + ` DESC) AS rw,
         ` + "user" + `
    FROM eod_inv WHERE guild=?
) sub
WHERE sub.user=?
`
