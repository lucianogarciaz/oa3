// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"fmt"
	"strings"

	"github.com/aarondl/oa3/support"
)

// Recursive definition of an array
type ArrayRecursive [][][]string

// VVValidateSchemaArrayRecursive validates the object and returns
// errors that can be returned to the user.
func (o ArrayRecursive) VVValidateSchema() support.Errors {
	var ctx []string
	var ers []error
	var errs support.Errors
	_, _, _ = ctx, ers, errs

	if err := support.ValidateMaxItems(o, 10); err != nil {
		ers = append(ers, err)
	}
	if err := support.ValidateMinItems(o, 2); err != nil {
		ers = append(ers, err)
	}

	for i, o := range o {
		var ers []error
		ctx = append(ctx, fmt.Sprintf("[%d]", i))
		if err := support.ValidateMaxItems(o, 8); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMinItems(o, 5); err != nil {
			ers = append(ers, err)
		}

		for i, o := range o {
			var ers []error
			ctx = append(ctx, fmt.Sprintf("[%d]", i))
			if err := support.ValidateMaxItems(o, 15); err != nil {
				ers = append(ers, err)
			}
			if err := support.ValidateMinItems(o, 12); err != nil {
				ers = append(ers, err)
			}

			errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
			ctx = ctx[:len(ctx)-1]
		}
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}

	return errs
}
