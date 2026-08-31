package main

// renderFormField returns a string containing HTML of a single form field. In case of select fType, it will retrun
// select tag with options. Value for select fType must be comma separated string which are use are
func renderFormField(label, name, fType string, value interface{}, id string, class string, required bool) string {
	// Format attributes with spaces first
	idAttr := ""
	if id != "" {
		idAttr = " id=\"" + template.HTMLEscapeString(id) + "\""
	}

	classAttr := ""
	if class != "" {
		classAttr = " class=\"" + template.HTMLEscapeString(class) + "\""
	}

	requiredAttr := ""
	if required {
		requiredAttr = " required"
	}

	// Escape all string values
	escapedName := template.HTMLEscapeString(name)
	escapedLabel := template.HTMLEscapeString(label)
	escapedType := template.HTMLEscapeString(fType)

	// Handle value specially as it's an interface{}
	escapedValue := ""
	if value != nil {
		escapedValue = template.HTMLEscapeString(fmt.Sprintf("%v", value))
	}

	if isValidForInput(fType) {
		return fmt.Sprintf(`%v<input%v%v name="%v" type="%v" value="%v"%v>`,
			escapedLabel, idAttr, classAttr, escapedName, escapedType, escapedValue, requiredAttr)
	}

	if fType == "select" {
		rawValueStr, ok := value.(string)
		if !ok {
			logs.Error("for select value must comma separated string that are the options for select")
			return ""
		}

		var selectBuilder strings.Builder
		selectBuilder.WriteString(fmt.Sprintf(`%v<select%v%v name="%v"></br>`,
			escapedLabel, idAttr, classAttr, escapedName))

		for _, option := range strings.Split(rawValueStr, ",") {
			escapedOption := template.HTMLEscapeString(option)
			selectBuilder.WriteString(fmt.Sprintf(`  <option value="%v"> %v </option></br>`,
				escapedOption, escapedOption))
		}

		selectBuilder.WriteString(`</select>`)
		return selectBuilder.String()
	}

	return fmt.Sprintf(`%v<%v%v%v name="%v"%v>%v</%v>`,
		escapedLabel, escapedType, idAttr, classAttr, escapedName, requiredAttr, escapedValue, escapedType)
}

// isValidForInput checks if fType is a valid value for the `type` property of an HTML input element.
	"html/template"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)
		}
	}
}

func TestRenderFormSecurity(t *testing.T) {
	type UserProfile struct {
		DisplayName string `form:"displayName,text,Name:"`
		Bio         string `form:",textarea"`
	}

	// Test case 1: Test proper escaping of special characters in attributes
	specialCharsProfile := UserProfile{
		DisplayName: `Special " ' < > & Characters`,
		Bio:         "Normal text content",
	}

	output := string(RenderForm(&specialCharsProfile))

	// Verify the output has all special characters properly escaped
	if strings.Contains(output, `"Special "`) {
		t.Errorf("Quotation mark not properly escaped in attribute")
	}

	if !strings.Contains(output, `value="Special`) {
		t.Errorf("Expected escaped attribute value not found")
	}

	// Test case 2: Test proper escaping of HTML-like content
	htmlContentProfile := UserProfile{
		DisplayName: "Normal Name",
		Bio:         `<div>Sample HTML content</div>`,
	}

	output = string(RenderForm(&htmlContentProfile))

	// Verify the output has HTML content properly escaped
	if strings.Contains(output, `<div>`) {
		t.Errorf("HTML tags not properly escaped in content")
	}

	if !strings.Contains(output, `&lt;div&gt;`) {
		t.Errorf("Expected escaped HTML not found")
	}
}
