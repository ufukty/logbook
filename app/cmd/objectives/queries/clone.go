package queries

func (o Objective) Clone() *Objective {
	return &Objective{
		Oid:       o.Oid,
		Vid:       o.Vid,
		Based:     o.Based,
		CreatedBy: o.CreatedBy,
		Pid:       o.Pid,
		Bupid:     o.Bupid,
		CreatedAt: o.CreatedAt,
	}
}
