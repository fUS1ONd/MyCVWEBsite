package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `validate:"required,min=3,max=10"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,gt=0"`
}

func TestValidate_ValidStruct(t *testing.T) {
	valid := TestStruct{
		Name:  "John",
		Email: "john@example.com",
		Age:   25,
	}

	err := Validate(valid)
	assert.NoError(t, err)
}

func TestValidate_RequiredField(t *testing.T) {
	invalid := TestStruct{
		Email: "john@example.com",
		Age:   25,
	}

	err := Validate(invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Name")
	assert.Contains(t, err.Error(), "required")
}

func TestValidate_MinLength(t *testing.T) {
	invalid := TestStruct{
		Name:  "Jo",
		Email: "john@example.com",
		Age:   25,
	}

	err := Validate(invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Name")
	assert.Contains(t, err.Error(), "at least 3")
}

func TestValidate_MaxLength(t *testing.T) {
	invalid := TestStruct{
		Name:  "VeryLongName",
		Email: "john@example.com",
		Age:   25,
	}

	err := Validate(invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Name")
	assert.Contains(t, err.Error(), "at most 10")
}

func TestValidate_Email(t *testing.T) {
	invalid := TestStruct{
		Name:  "John",
		Email: "invalid-email",
		Age:   25,
	}

	err := Validate(invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Email")
	assert.Contains(t, err.Error(), "valid email")
}

func TestValidate_GreaterThan(t *testing.T) {
	type TestStructWithGT struct {
		Value int `validate:"gt=10"`
	}

	invalid := TestStructWithGT{
		Value: 5,
	}

	err := Validate(invalid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Value")
	assert.Contains(t, err.Error(), "greater than")
}
