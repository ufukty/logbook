package database

func (o Objective) Clone() *Objective {
	return &Objective{
		Oid:       o.Oid,
		Vid:       o.Vid,
		Based:     o.Based,
		Type:      o.Type,
		Content:   o.Content,
		Creator:   o.Creator,
		CreatedAt: o.CreatedAt,
	}
}
