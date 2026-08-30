package main


	ru := *u
	if _, has := ru.User.Password(); has {
		ru.User = url.UserPassword(ru.User.Username(), "redacted")
	}
	q := ru.Query()
	if q.Get("sshkey") != "" {
		q.Set("sshkey", "redacted")
		ru.RawQuery = q.Encode()
	}
	return ru.String()
}
				Path:   "this:that",
				User:   url.UserPassword("user", "password"),
			},
			want: "http://user:redacted@host.tld/this:that",
		},
		{
			name: "blank Password",
				Path:   "this:that",
				User:   url.UserPassword("", "password"),
			},
			want: "http://:redacted@host.tld/this:that",
		},
		{
			name: "blank Username, blank Password",
			url:  nil,
			want: "",
		},
		{
			name: "non-blank SSH key in URL query parameter",
			url: &url.URL{
				Scheme:   "ssh",
				User:     url.User("git"),
				Host:     "github.com",
				Path:     "hashicorp/go-getter-test-private.git",
				RawQuery: "sshkey=LS0tLS1CRUdJTiBPUE",
			},
			want: "ssh://git@github.com/hashicorp/go-getter-test-private.git?sshkey=redacted",
		},
		{
			name: "blank SSH key in URL query parameter",
			url: &url.URL{
				Scheme:   "ssh",
				User:     url.User("git"),
				Host:     "github.com",
				Path:     "hashicorp/go-getter-test-private.git",
				RawQuery: "sshkey=",
			},
			want: "ssh://git@github.com/hashicorp/go-getter-test-private.git?sshkey=",
		},
	}

	for _, tt := range cases {
