package database

func (o Objective) Clone() *Objective {
	return &Objective{
		Oid:       o.Oid,
		Vid:       o.Vid,
		Based:     o.Based,
		Content:   o.Content,
		Creator:   o.Creator,
		CreatedAt: o.CreatedAt,
	}
}
