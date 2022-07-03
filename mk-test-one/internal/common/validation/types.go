package validation

import (
	"fmt"
	"regexp"
)

// TODO: move to any errors, not only validation?

const (
	notImplementedCode = "not_implemented"
	unknownField       = "unknown"
)

// structure of backend errors
// https://oua.atlassian.net/wiki/spaces/ORIENTCODE/pages/1535770629/Unified+server+error+reporting
type Result struct {
	Details string   `json:"details"`
	Code    string   `json:"code"`
	Errors  []*Error `json:"errors"`
}

type Error struct {
	Name  string         `json:"name"`
	Codes []ErrorDetails `json:"codes"`
}

type ErrorDetails struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func NewResult() *Result {
	return &Result{
		Details: "",
		Errors:  make([]*Error, 0),
	}
}

func UnmarshalError(err error) *Result {
	return &Result{
		Details: err.Error(),
		Code:    notImplementedCode,
	}
}

var (
	rxFieldName    = regexp.MustCompile("field=(.+?),")
	rxErrorDetails = regexp.MustCompile(`internal=(.+)`)
)

func UnmarshalDetailedError(err error) *Result {
	return unmarshalError(notImplementedCode, err)
}

func unmarshalError(code string, err error) *Result {
	out := NewResult()

	fieldNameGroup := rxFieldName.FindStringSubmatch(err.Error())
	errorDetailsGroup := rxErrorDetails.FindStringSubmatch(err.Error())

	fieldName := unknownField // TODO: ticket was created to not allow this, need some fancy parsing to always have the field name
	if len(fieldNameGroup) > 1 {
		fieldName = fieldNameGroup[1]
	}

	errorDetails := err.Error()
	if len(errorDetailsGroup) > 1 {
		errorDetails = errorDetailsGroup[1]
	}

	out.AddFieldError(fieldName, ErrorDetails{
		Message: errorDetails,
		Code:    code,
	})

	return out
}

func DBOperationError(err error) *Result {
	return &Result{
		Details: err.Error(),
		Code:    notImplementedCode,
	}
}

func BothEmailAndPhoneProvided() *Result {
	return &Result{
		Details: "can't provide both email and phone during registration",
		Code:    notImplementedCode,
	}
}

func CaptchaError(err error) *Result {
	return &Result{
		Details: err.Error(),
		Code:    notImplementedCode,
	}
}

// NoCodeError is a generic error, on request from FE will refactor responses
// that uses this if they need a code
func NoCodeError(err error) *Result {
	return &Result{
		Details: err.Error(),
	}
}

// CodeError is a generic error with code value
func CodeError(code string, err error) *Result {
	return &Result{
		Code:    code,
		Details: err.Error(),
	}
}

func (r *Result) IsValid() bool {
	return len(r.Errors) == 0 && r.Details == ""
}

func (r *Result) AddResult(result *Result) {
	r.Errors = append(r.Errors, result.Errors...)
}

func (r *Result) AddFieldError(field string, ed ErrorDetails) *Result {
	for _, e := range r.Errors {
		if e.Name == field {
			e.Codes = append(e.Codes, ed)
			return r
		}
	}
	r.Errors = append(r.Errors, &Error{
		Name:  field,
		Codes: []ErrorDetails{ed},
	})
	return r
}

func (r *Result) AddDetails(details string, formatArgs ...interface{}) *Result {
	r.Details = fmt.Sprintf(details, formatArgs...)
	return r
}

func (r *Result) AddCode(code string) *Result {
	r.Code = code
	return r
}

func EitherPhoneOrEmail() ErrorDetails {
	return ErrorDetails{
		Message: "either phone or e-mail must be provided",
		Code:    notImplementedCode,
	}
}

func EmptyBirthDate() ErrorDetails {
	return ErrorDetails{
		Message: "empty birthday",
		Code:    notImplementedCode,
	}
}

func NotOnlyLetters() ErrorDetails {
	return ErrorDetails{
		Message: "field must contain only letters",
		Code:    notImplementedCode,
	}
}

func UserAlreadyExists() ErrorDetails {
	return ErrorDetails{
		Message: "user already exists",
		Code:    notImplementedCode,
	}
}

func UserNotFound() ErrorDetails {
	return ErrorDetails{
		Message: "user not found",
		Code:    notImplementedCode,
	}
}

func Unauthorized() ErrorDetails {
	return ErrorDetails{
		Message: "unauthorized",
		Code:    notImplementedCode,
	}
}

func EmptyPassword() ErrorDetails {
	return ErrorDetails{
		Message: "empty password provided",
		Code:    notImplementedCode,
	}
}

func InvalidPassword() ErrorDetails {
	return ErrorDetails{
		Message: "password must contain at least 1 lowercased letter, 1 capital letter, 1 digit, 1 special char and be minimum 8 chars long",
		Code:    notImplementedCode,
	}
}

func RulesNotAccepted() ErrorDetails {
	return ErrorDetails{
		Message: "rules were not accepted",
		Code:    notImplementedCode,
	}
}

func NameIsTooShort() ErrorDetails {
	return ErrorDetails{
		Message: "name cannot be less than 2 characters",
		Code:    notImplementedCode,
	}
}

func TooYoungAge() ErrorDetails {
	return ErrorDetails{
		Message: "age must be more than 18 years",
		Code:    notImplementedCode,
	}
}

func WrongPhoneFormat() ErrorDetails {
	return ErrorDetails{
		Message: "wrong phone format: it should contain only digits",
		Code:    notImplementedCode,
	}
}

func InvalidEmail() ErrorDetails {
	return ErrorDetails{
		Message: "email is empty or has invalid format",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func UnknownCountry() ErrorDetails {
	return ErrorDetails{
		Message: "such country does not exist",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func InvalidPhone() ErrorDetails {
	return ErrorDetails{
		Message: "phone is empty",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func InvalidIPAddress(ip string) ErrorDetails {
	return ErrorDetails{
		Message: fmt.Sprintf("ip_address is empty or has invalid format: %s", ip),
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func EmptyDevice() ErrorDetails {
	return ErrorDetails{
		Message: "empty device provided",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func EmptyRefreshToken() ErrorDetails {
	return ErrorDetails{
		Message: "empty refresh token provided",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func InvalidAntiPhishingCode() ErrorDetails {
	return ErrorDetails{
		Message: "anti-phishing code is invalid",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func InvalidCountryCallingCodeFormat() ErrorDetails {
	return ErrorDetails{
		Message: "country calling code format is invalid",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func WrongCountryCallingCode() ErrorDetails {
	return ErrorDetails{
		Message: "country calling code does not match selected country",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func InvalidPageLimit() ErrorDetails {
	return ErrorDetails{
		Message: "page limit is invalid",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func InvalidOrderColumn() ErrorDetails {
	return ErrorDetails{
		Message: "order column is invalid",
		Code:    notImplementedCode, // TODO error code should be implemented
	}
}

func InvalidOtpCode() ErrorDetails {
	return ErrorDetails{
		Message: "invalid otp code provided",
		Code:    notImplementedCode,
	}
}

func InvalidAction() ErrorDetails {
	return ErrorDetails{
		Message: "action is invalid",
		Code:    notImplementedCode,
	}
}

func Invalid2FAMethod() ErrorDetails {
	return ErrorDetails{
		Message: "method is invalid",
		Code:    notImplementedCode,
	}
}

func InvalidCode() ErrorDetails {
	return ErrorDetails{
		Message: "code is invalid",
		Code:    notImplementedCode,
	}
}

func InvalidKey() ErrorDetails {
	return ErrorDetails{
		Message: "key is empty",
		Code:    notImplementedCode,
	}
}

func InvalidResetToken() ErrorDetails {
	return ErrorDetails{
		Message: "reset-token is invalid",
		Code:    notImplementedCode,
	}
}

func InvalidConfirmPassword() ErrorDetails {
	return ErrorDetails{
		Message: "confirm password is invalid",
		Code:    notImplementedCode,
	}
}

func InvalidKeys() ErrorDetails {
	return ErrorDetails{
		Message: "phone or email method is missing",
		Code:    notImplementedCode,
	}
}
