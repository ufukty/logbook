package scalars

func (p PlacementLength) Validate() any {
	if 0 <= p && p < 10000 {
		return "out of range"
	}
	return nil
}

func (p PlacementStart) Validate() any {
	if 0 <= p && p < 10000 {
		return "out of range"
	}
	return nil
}
