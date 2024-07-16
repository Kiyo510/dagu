// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DagLogGridItem dag log grid item
//
// swagger:model dagLogGridItem
type DagLogGridItem struct {

	// name
	// Required: true
	Name *string `json:"Name"`

	// vals
	// Required: true
	Vals []int64 `json:"Vals"`
}

// Validate validates this dag log grid item
func (m *DagLogGridItem) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVals(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DagLogGridItem) validateName(formats strfmt.Registry) error {

	if err := validate.Required("Name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *DagLogGridItem) validateVals(formats strfmt.Registry) error {

	if err := validate.Required("Vals", "body", m.Vals); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this dag log grid item based on context it is used
func (m *DagLogGridItem) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DagLogGridItem) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DagLogGridItem) UnmarshalBinary(b []byte) error {
	var res DagLogGridItem
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}