package columns

func (st *SessionToken) Set(v string) error {
	*st = SessionToken(v)
	return nil
}
