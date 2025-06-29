package ontology

import "errors"

// ValidationAgent performs minimal validation on TTL output.
type ValidationAgent struct{}

func NewValidationAgent() *ValidationAgent { return &ValidationAgent{} }

func (v *ValidationAgent) Validate(content string) error {
    if content == "" {
        return errors.New("empty ontology")
    }
    return nil
}
