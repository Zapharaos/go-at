package goat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTemplate_Render tests the Render method of the Template struct
func TestTemplate_Render(t *testing.T) {
	t.Run("empty template", func(t *testing.T) {
		tmpl := Template{}

		result, err := tmpl.Render()
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("template with incorrect content", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content",
			Data:       nil,
		}

		result, err := tmpl.Render()
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("template without data", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       nil,
		}

		result, err := tmpl.Render()
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("template with incorrect data", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       map[string]interface{}{"Content": func() {}},
		}

		result, err := tmpl.Render()
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("template with wrong data", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       map[string]string{"Wrong": "wrong"},
		}

		result, err := tmpl.Render()
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
	t.Run("successful render", func(t *testing.T) {
		tmpl := Template{
			Name:       "test_template",
			ContentRaw: "{{.Content}}",
			Data:       map[string]string{"Content": "Hello World"},
		}

		result, err := tmpl.Render()
		assert.NoError(t, err)
		assert.Equal(t, "Hello World", result)
	})
}
