package verifier

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewVerify(t *testing.T) {
	v := NewVerify()
	assert.NotNil(t, v)
}

func TestVerify_Struct(t *testing.T) {
	cases := map[string]struct {
		input    interface{}
		expected error
	}{
		"nil struct": {
			input:    nil,
			expected: ErrStructIsNil,
		},
		"valid struct": {
			input: struct {
				Name string `validate:"required"`
			}{
				Name: "test",
			},
			expected: nil,
		},
		"invalid struct": {
			input: struct {
				Name string `validate:"required"`
			}{
				Name: "",
			},
			expected: fmt.Errorf("name is not valid\n"),
		},
		"invalid struct with multiple errors": {
			input: struct {
				Name string `validate:"required"`
				Age  int    `validate:"required"`
			}{
				Name: "",
				Age:  0,
			},
			expected: fmt.Errorf("name is not valid\nage is not valid\n"),
		},
	}

	for c, tt := range cases {
		t.Run(c, func(t *testing.T) {
			v := NewVerify()
			err := v.Struct(tt.input)
			if tt.expected == nil {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.Error(), err.Error())
			}
		})
	}
}

func TestVerify_Slice(t *testing.T) {
	type s1 struct {
		Name string `validate:"required"`
	}

	cases := map[string]struct {
		input    []interface{}
		expected error
	}{
		"nil slice": {
			input:    nil,
			expected: ErrStructIsNil,
		},
		"valid slice": {
			input: []interface{}{
				s1{Name: "test"},
				s1{Name: "test"},
			},
			expected: nil,
		},
		"invalid slice": {
			input: []interface{}{
				s1{Name: ""},
			},
			expected: fmt.Errorf("name is not valid\n"),
		},
		"invalid slice with multiple errors": {
			input: []interface{}{
				s1{Name: ""},
				s1{Name: ""},
			},
			expected: fmt.Errorf("name is not valid\n"),
		},
		"mixed slice": {
			input: []interface{}{
				s1{Name: "test"},
				struct {
					Age int `validate:"required"`
				}{Age: 0},
			},
			expected: fmt.Errorf("age is not valid\n"),
		},
	}

	for c, tt := range cases {
		t.Run(c, func(t *testing.T) {
			v := NewVerify()
			err := v.Slice(tt.input)
			if tt.expected == nil {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.Error(), err.Error())
			}
		})
	}
}
