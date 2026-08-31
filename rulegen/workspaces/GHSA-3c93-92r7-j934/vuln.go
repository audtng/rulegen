package main

	return false
}

func CanAllowURL(url string, allowedHosts []string) bool {
	allow := false
	if len(allowedHosts) == 0 {
		return true
	}
	for _, host := range allowedHosts {
		if strings.HasPrefix(url, host) {
			return true
		}
	}
	return allow
}

func GetQueryBody(ctx context.Context, query models.Query) io.Reader {
	logger := backend.Logger.FromContext(ctx)
	var body io.Reader
}

func TestCanAllowURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		allowedHosts []string
		want         bool
	}{
		{
			url:  "https://foo.com",
			want: true,
		},
		{
			url:          "https://foo.com",
			allowedHosts: []string{"https://foo.com"},
			want:         true,
		},
		{
			url:          "https://bar.com",
			allowedHosts: []string{"https://foo.com"},
			want:         false,
		},
		{
			name:         "should match only case sensitive URL",
			url:          "https://FOO.com",
			allowedHosts: []string{"https://foo.com"},
			want:         false,
		},
		{
			url:          "https://bar.com/",
			allowedHosts: []string{"https://foo.com/", "https://bar.com/", "https://baz.com/"},
			want:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := infinity.CanAllowURL(tt.url, tt.allowedHosts)
			assert.Equal(t, tt.want, got)
		})
	}
}
}

func NormalizeURL(u string) string {
	urlArray := strings.Split(u, "/")
	if strings.HasPrefix(u, "https://github.com") && len(urlArray) > 5 && urlArray[5] == "blob" && urlArray[4] != "blob" && urlArray[3] != "blob" {
		u = strings.Replace(u, "https://github.com", "https://raw.githubusercontent.com", 1)
			query: models.Query{
				URL: "0.0.0.0",
			},
			want: "0.0.0.0",
		},
		{
			settings: models.InfinitySettings{},
import (
	"context"
	"encoding/json"
	"fmt"
	"net/textproto"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
		}
		return nil
	}
	if s.AuthenticationMethod != AuthenticationMethodAzureBlob && s.AuthenticationMethod != AuthenticationMethodNone && len(s.AllowedHosts) < 1 {
		return ErrInvalidConfigHostNotAllowed
	}
	if s.HaveSecureHeaders() && len(s.AllowedHosts) < 1 {
		return ErrInvalidConfigHostNotAllowed
	}
	if len(s.KeepCookies) > 0 && len(s.AllowedHosts) < 1 {
		return ErrInvalidConfigHostNotAllowed
	}
	return nil
}

func (s *InfinitySettings) HaveSecureHeaders() bool {
	if len(s.CustomHeaders) > 0 {
		for k := range s.CustomHeaders {
			if textproto.CanonicalMIMEHeaderKey(k) == "Accept" {
			}
			return true
		}
		return false
	}
	return false
}

type RefData struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
import (
	"context"
	"errors"
	"testing"

	"github.com/grafana/grafana-infinity-datasource/pkg/models"
		SecureQueryFields: map[string]string{
			"foo": "bar",
		},
		KeepCookies:               []string{"cookie1", "cookie2"},
	}, gotSettings)
}

		})
	}
}
