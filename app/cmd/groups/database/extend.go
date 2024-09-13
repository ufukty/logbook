package database

func lastsix[S ~string](id S) S {
	return id[max(0, len(id)-6):]
}
