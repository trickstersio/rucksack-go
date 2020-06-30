package apierr

// Validator contains information about some payload validation
type Validator struct {
	err *Err
}

// AddError is adding new validation error
func (v *Validator) AddError(err *Err) {
	if v.err == nil {
		v.err = New()
	}

	v.err = v.err.Add(err)
}

// Error returns validation error
func (v *Validator) Error() *Err {
	return v.err
}
