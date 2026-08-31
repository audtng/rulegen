package main

// EncodeIdentity encodes a plugin identity string for a plugin with the given
// name. If the name is invalid, it returns an empty string.
func EncodeIdentity(name string, data []byte) string {
	if !validPluginName(name) {
		return ""
	}
	s, _ := bech32.Encode("AGE-PLUGIN-"+strings.ToUpper(name)+"-", data)
	return s
}
	}
	name = strings.TrimSuffix(strings.TrimPrefix(hrp, "AGE-PLUGIN-"), "-")
	name = strings.ToLower(name)
	if !validPluginName(name) {
		return "", nil, fmt.Errorf("invalid plugin name: %q", name)
	}
	return name, data, nil
}

// EncodeRecipient encodes a plugin recipient string for a plugin with the given
// name. If the name is invalid, it returns an empty string.
func EncodeRecipient(name string, data []byte) string {
	if !validPluginName(name) {
		return ""
	}
	s, _ := bech32.Encode("age1"+strings.ToLower(name), data)
	return s
}
		return "", nil, fmt.Errorf("not a plugin recipient: %v", err)
	}
	name = strings.TrimPrefix(hrp, "age1")
	if !validPluginName(name) {
		return "", nil, fmt.Errorf("invalid plugin name: %q", name)
	}
	return name, data, nil
}

func validPluginName(name string) bool {
	if name == "" {
		return false
	}
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-._"
	for _, r := range name {
		if !strings.ContainsRune(allowed, r) {
			return false
		}
	}
	return true
}
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	exec "golang.org/x/sys/execabs"

func NewIdentityWithoutData(name string, ui *ClientUI) (*Identity, error) {
	s := EncodeIdentity(name, nil)
	if s == "" {
		return nil, fmt.Errorf("invalid plugin name: %q", name)
	}
	return &Identity{
		name: name, encoding: s, ui: ui,
	}, nil
	path := "age-plugin-" + name
	if testOnlyPluginPath != "" {
		path = filepath.Join(testOnlyPluginPath, path)
	} else if strings.ContainsRune(name, os.PathSeparator) {
		return nil, fmt.Errorf("invalid plugin name: %q", name)
	}
	cmd := exec.Command(path, "--age-plugin="+protocol)

