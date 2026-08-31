package main

// renderFormField returns a string containing HTML of a single form field. In case of select fType, it will retrun
// select tag with options. Value for select fType must be comma separated string which are use are
func renderFormField(label, name, fType string, value interface{}, id string, class string, required bool) string {
	if id != "" {
		id = " id=\"" + id + "\""
	}

	if class != "" {
		class = " class=\"" + class + "\""
	}

	requiredString := ""
	if required {
		requiredString = " required"
	}

	if isValidForInput(fType) {
		return fmt.Sprintf(`%v<input%v%v name="%v" type="%v" value="%v"%v>`, label, id, class, name, fType, value, requiredString)
	}

	if fType == "select" {
		valueStr, ok := value.(string)
		if !ok {
			logs.Error("for select value must comma separated string that are the options for select")
			return ""
		}

		var selectBuilder strings.Builder
		selectBuilder.WriteString(fmt.Sprintf(`%v<select%v%v name="%v"></br>`, label, id, class, name))

		for _, option := range strings.Split(valueStr, ",") {
			selectBuilder.WriteString(fmt.Sprintf(`  <option value="%v"> %v </option></br>`, option, option))
		}

		selectBuilder.WriteString(`</select>`)
		return selectBuilder.String()
	}

	return fmt.Sprintf(`%v<%v%v%v name="%v"%v>%v</%v>`, label, fType, id, class, name, requiredString, value, fType)
}

// isValidForInput checks if fType is a valid value for the `type` property of an HTML input element.
	"html/template"
	"net/url"
	"reflect"
	"testing"
	"time"
)
		}
	}
}
