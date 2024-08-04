package database

func (o Objective) Clone() *Objective {
	return &Objective{
		Oid:       o.Oid,
		Vid:       o.Vid,
		Based:     o.Based,
		CreatedBy: o.CreatedBy,
		Props:     o.Props,
		CreatedAt: o.CreatedAt,
	}
}
