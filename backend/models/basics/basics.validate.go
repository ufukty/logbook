// Code generated by govalid. DO NOT EDIT.

package basics

import "fmt"

func (b Boolean) Validate() error {
	switch b {
	case False:
		return nil
	case True:
		return nil
	}
	return fmt.Errorf("invalid value")
}
