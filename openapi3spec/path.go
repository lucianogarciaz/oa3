package openapi3spec

import (
	"errors"
	"fmt"
	"strings"
)

// Paths holds the relative paths to the individual endpoints and their
// operations. The path is appended to the URL from the Server Object in order
// to construct the full URL. The Paths MAY be empty, due to ACL constraints.
//
// Technically Paths can have extensions as per the spec, but we make a choice
// not to conform in order to be able to avoid an object graph that is also
// against the spec: OpenAPI3.Paths.Paths["/url"]
type Paths map[string]*PathRef

// Path describes the operations available on a single path. A Path Item MAY
// be empty, due to ACL constraints. The path itself is still exposed to the
// documentation viewer but they will not know which operations and parameters
// are available.
type Path struct {
	Summary     *string        `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description *string        `json:"description,omitempty" yaml:"description,omitempty"`
	Get         *Operation     `json:"get,omitempty" yaml:"get,omitempty"`
	Put         *Operation     `json:"put,omitempty" yaml:"put,omitempty"`
	Post        *Operation     `json:"post,omitempty" yaml:"post,omitempty"`
	Delete      *Operation     `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options     *Operation     `json:"options,omitempty" yaml:"options,omitempty"`
	Head        *Operation     `json:"head,omitempty" yaml:"head,omitempty"`
	Patch       *Operation     `json:"patch,omitempty" yaml:"patch,omitempty"`
	Trace       *Operation     `json:"trace,omitempty" yaml:"trace,omitempty"`
	Servers     []Server       `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters  []ParameterRef `json:"parameters,omitempty" yaml:"parameters,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// Validate path
func (p *Path) Validate(cmps Components, pathTemplates []string, opIDs map[string]struct{}) error {
	if p == nil {
		return nil
	}

	if p.Summary != nil && len(strings.TrimSpace(*p.Summary)) == 0 {
		return errors.New("summary if present must not be blank")
	}
	if p.Description != nil && len(strings.TrimSpace(*p.Description)) == 0 {
		return errors.New("description if present must not be blank")
	}

	if err := p.Get.Validate(cmps, pathTemplates, opIDs); err != nil {
		return fmt.Errorf("get.%w", err)
	}
	if err := p.Put.Validate(cmps, pathTemplates, opIDs); err != nil {
		return fmt.Errorf("put.%w", err)
	}
	if err := p.Post.Validate(cmps, pathTemplates, opIDs); err != nil {
		return fmt.Errorf("post.%w", err)
	}
	if err := p.Delete.Validate(cmps, pathTemplates, opIDs); err != nil {
		return fmt.Errorf("delete.%w", err)
	}
	if err := p.Options.Validate(cmps, pathTemplates, opIDs); err != nil {
		return fmt.Errorf("options.%w", err)
	}
	if err := p.Head.Validate(cmps, pathTemplates, opIDs); err != nil {
		return fmt.Errorf("head.%w", err)
	}
	if err := p.Patch.Validate(cmps, pathTemplates, opIDs); err != nil {
		return fmt.Errorf("patch.%w", err)
	}
	if err := p.Trace.Validate(cmps, pathTemplates, opIDs); err != nil {
		return fmt.Errorf("trace.%w", err)
	}

	for i, s := range p.Servers {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("servers[%d].%w", i, err)
		}
	}

	if err := paramDuplicateKeyCheck(p.Parameters); err != nil {
		return fmt.Errorf("parameters.%w", err)
	}

	for i, p := range p.Parameters {
		if err := p.Validate(cmps, pathTemplates); err != nil {
			return fmt.Errorf("parameters[%d].%w", i, err)
		}
	}

	return nil
}

// PathRef refers to a path item
type PathRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*Path
}

// Validate path ref
func (p *PathRef) Validate(cmps Components, pathTemplates []string, opIDs map[string]struct{}) error {
	// Don't validate references
	if p == nil || len(p.Ref) != 0 {
		return nil
	}

	if err := p.Path.Validate(cmps, pathTemplates, opIDs); err != nil {
		return err
	}

	return nil
}
