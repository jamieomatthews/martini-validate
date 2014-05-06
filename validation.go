package validation

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Validation struct {
	Errors  Errors
	Request *http.Request
}

//represents one 'set' of validation errors.
type Set struct {
	Field          interface{} //a pointer to the passed in feild
	Key            string      //string key pulled from the field
	isValid        bool
	classification string
	message        string
	Validation     *Validation //for now, keep a reference to the validation to map errors back
}

func New(errors Errors, request *http.Request) *Validation {
	return &Validation{Errors: errors, Request: request}
}

func NewDefault() *Validation {
	v := &Validation{Errors: &errors{}, Request: &http.Request{}}
	return v
}

func (v *Validation) Validate(field interface{}, key string) *Set {
	s := &Set{Field: field, Key: key, isValid: true, Validation: v}

	return s
}

// func (v *Validation) ToJson() []byte {
// 	m := make(map[string][]string)
// 	for err := range v.Errors {
// 		m[err.FieldNames[0]] =
// 	}
// }

// returns true if the validator has 1 or more errors
func (s *Set) HasErrors() bool {
	return s.Validation.Errors.HasErrors()
}

func (s *Set) MaxLength(maxLength int) *Set {
	return s.validate(MaxLength{MaxLength: maxLength}, s.Len())
}

func (s *Set) MinLength(minLength int) *Set {
	fmt.Println("Min Length: '", s.Field, "'")
	return s.validate(MinLength{MinLength: minLength}, s.Len())
}

func (s *Set) Match(strMatch string, regex *regexp.Regexp) *Set {
	return s.validate(Matches{Regex: regex}, strMatch)
}

func (s *Set) Email() *Set {
	var emailPattern = regexp.MustCompile("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$")
	return s.validate(Email{Matches{emailPattern}}, s.toString())
}

func (s *Set) Classify(str string) *Set {
	s.classification = str
	return s
}

func (s *Set) Message(message string) *Set {
	s.message = message
	return s
}
func (s *Set) TrimSpace() *Set {
	if str, ok := s.Field.(string); ok {
		s.Field = strings.TrimSpace(str)
	}
	return s
}

func (s *Set) getMessage(val Validator) string {
	if s.message == "" {
		return val.DefaultMessage()
	}
	return s.message
}

func (s *Set) Len() int {
	if str, ok := s.Field.(string); ok {
		return len(str)
	}
	if array, ok := s.Field.([]interface{}); ok {
		return len(array)
	}
	if m, ok := s.Field.(map[interface{}]interface{}); ok {
		return len(m)
	}
	return 0
}

func (s *Set) toString() string {
	if str, ok := s.Field.(string); ok {
		return str
	}
	panic("This method requires a string value")
}

//runs the validation rule, returns true if the rule passed
func (s *Set) validate(validator Validator, obj interface{}) *Set {
	//check if the rule is valid
	if validator.IsValid(obj) {
		fmt.Println("Validated")
		return s
	}

	//else, add a new validation error
	s.Validation.Errors.Add([]string{s.Key}, s.classification, s.getMessage(validator))
	//s.Validation.AddError(binding.Error{FieldNames: []string{s.Key}, Classification: s.classification, Message: s.getMessage(validator)})
	s.isValid = false
	return s
}
