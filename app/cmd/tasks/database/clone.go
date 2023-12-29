package database

func (o Objective) Clone() *Objective {
	return &Objective{
		Oid:         o.Oid,
		ParentId:    o.ParentId,
		Vid:         o.Vid,
		Creator:     o.Creator,
		Text:        o.Text,
		CreatedAt:   o.CreatedAt,
		CompletedAt: o.CompletedAt,
		ArchivedAt:  o.ArchivedAt,
	}
}
