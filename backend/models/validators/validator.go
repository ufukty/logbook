package validators

import (
	"github.com/ufukty/gohandlers/pkg/validator"
)

var (
	// CreditCard   = validator.ForStrings(`^(?:4[0-9]{12}(?:[0-9]{3})?)$`, 13, 19)
	// HtmlDate     = validator.ForStrings(`^\d{4}-\d{2}-\d{2}$`, 6, 8)
	// HtmlDatetime = validator.ForStrings(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$`, 9, 13)
	// HtmlTime     = validator.ForStrings(`^\d{2}:\d{2}$`, 3, 5)
	// Numeric      = validator.ForStrings(`^[1-9][0-9]*$`, 0, 100)
	// Text         = validator.ForStrings(`^[\p{L}0-9 ,.?!'’“”-]+$`, 0, 10000)
	// Url          = validator.ForStrings(`^[\p{L}0-9._%+-]+@[\p{L}0-9.-]+\.[\p{L}]{2,}$`, 0, 10000)
	AntiCsrfToken = validator.ForStrings(`[ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_]+$`, 32, 32) // defined in std lib [base64.URLEncoding]
	Email         = validator.ForStrings(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, 6, 150)
	GroupName     = validator.ForStrings(`[\p{L} ]+`, 2, 100)
	HumanName     = validator.ForStrings(`^\p{L}+([ '-]\p{L}+)*$`, 6, 100)
	PhoneNumber   = validator.ForStrings(`^\+?(\d{1,3})?[ -]?(\d{3})[ -]?(\d{3})[ -]?(\d{4})$`, 10, 15)
	SessionToken  = validator.ForStrings(`[A-Za-z0-9-_]+$`, 256, 256) // pattern is as defined in std lib base64.URLEncoding
	Username      = validator.ForStrings(`^[a-zA-Z]+[a-zA-Z0-9\_\.\-]*$`, 3, 50)
	Uuid          = validator.ForStrings(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, len("00000000-0000-0000-0000-000000000000"), len("00000000-0000-0000-0000-000000000000"))
)
