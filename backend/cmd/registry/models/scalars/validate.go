package scalars

import "logbook/models/validators"

func (v InstanceId) Validate() any { return validators.Uuid.Validate(string(v)) }
