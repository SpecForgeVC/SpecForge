package db

// DB exposes the underlying DBTX for custom queries
func (q *Queries) DB() DBTX {
	return q.db
}
