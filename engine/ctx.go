package engine

import "context"

const ctxDB = "db"

func CtxWithDB(ctx context.Context, db *DB) context.Context {
	return context.WithValue(ctx, ctxDB, db)
}

func CtxGetDB(ctx context.Context) *DB {
	v := ctx.Value(ctxDB)
	if v == nil {
		return nil
	}
	return v.(*DB)
}
