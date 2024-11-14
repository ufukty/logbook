package endpoints

import (
	"fmt"
	"strings"
)

func (c CreateAccountRequest) Validate() error {
	errs := []string{}
	err := c.AntiCsrfToken.Validate()
	if err != nil {
		errs = append(errs, fmt.Sprintf("acsrft (%s)", err.Error()))
	}
	err = c.Firstname.Validate()
	if err != nil {
		errs = append(errs, fmt.Sprintf("firstname (%s)", err.Error()))
	}
	err = c.Lastname.Validate()
	if err != nil {
		errs = append(errs, fmt.Sprintf("lastname (%s)", err.Error()))
	}
	err = c.Birthday.Validate()
	if err != nil {
		errs = append(errs, fmt.Sprintf("birthday (%s)", err.Error()))
	}
	err = c.Country.Validate()
	if err != nil {
		errs = append(errs, fmt.Sprintf("country (%s)", err.Error()))
	}
	err = c.EmailGrant.Validate()
	if err != nil {
		errs = append(errs, fmt.Sprintf("email (%s)", err.Error()))
	}
	err = c.PasswordGrant.Validate()
	if err != nil {
		errs = append(errs, fmt.Sprintf("password (%s)", err.Error()))
	}
	if len(errs) > 0 {
		return fmt.Errorf("issues: %s", strings.Join(errs, ", "))
	}
	return nil
}
