package models

import (
	"fmt"
)

func (t *VM) ValidateTransitionTo(to State) error {
	var valid bool
	from := t.State
	switch to {
	case StateFree:
		valid = ( from == StateUsing || from == StateProvisioning )
	case StateProvisioning:
		valid = ( from == StateFree || from == StateUsing )
	case StateUnknown:
		valid = true
	case StateUsing:
		valid = from == StateProvisioning
	}

	if !valid {
		return NewError(
			ErrorTypeInvalidStateTransition,
			fmt.Sprintf("Cannot transition from %v to %v", from, to),
		)
	}

	return nil
}

