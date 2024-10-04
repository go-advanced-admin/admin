package fields

import (
	"github.com/go-advanced-admin/admin/internal/form"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextFieldHTML(t *testing.T) {
	textField := &TextField{}
	err := textField.RegisterName("username")
	assert.Nil(t, err)
	placeholder := "Enter username"
	textField.Placeholder = &placeholder
	textField.Required = true
	maxLength := uint(20)
	textField.MaxLength = &maxLength

	html, err := textField.HTML()
	assert.Nil(t, err)
	assert.Contains(t, html, `name="username"`)
	assert.Contains(t, html, `placeholder="Enter username"`)
	assert.Contains(t, html, `maxlength="20"`)
	assert.Contains(t, html, `required`)
}

func TestTextFieldConversions(t *testing.T) {
	textField := &TextField{}

	htmlType, err := textField.GoTypeToHTMLType("test")
	assert.Nil(t, err)
	assert.Equal(t, form.HTMLType("test"), htmlType)

	goType, err := textField.HTMLTypeToGoType("test")
	assert.Nil(t, err)
	assert.Equal(t, "test", goType)
}
