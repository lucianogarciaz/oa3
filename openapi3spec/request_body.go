package openapi3spec

import (
	"errors"
	"fmt"
	"strings"
)

// RequestBody for an operation
type RequestBody struct {
	Description *string               `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty"`
}

// Validate request body
func (r *RequestBody) Validate(cmps Components) error {
	if r.Description != nil && len(strings.TrimSpace(*r.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	if len(r.Content) == 0 {
		return fmt.Errorf("content must not be empty")
	}
	for k, c := range r.Content {
		if err := c.Validate(cmps); err != nil {
			return fmt.Errorf("content(%s).%w", k, err)
		}
	}

	return nil
}

// RequestBodyRef refers to a request body
type RequestBodyRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*RequestBody
}

// Validate request body ref
func (r *RequestBodyRef) Validate(c Components) error {
	// Don't validate references
	if r == nil || len(r.Ref) != 0 {
		return nil
	}

	if err := r.RequestBody.Validate(c); err != nil {
		return err
	}

	return nil
}
