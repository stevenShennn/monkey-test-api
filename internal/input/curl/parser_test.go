package curl

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, reqs []input.Request)
	}{
		{
			name: "基本 GET 请求",
			input: `curl 'https://api.example.com/users' \
					-H 'accept: application/json'`,
			wantErr: false,
			check: func(t *testing.T, reqs []input.Request) {
				assert.Len(t, reqs, 1)
				assert.Equal(t, "GET", reqs[0].Method)
				assert.Equal(t, "https://api.example.com/users", reqs[0].URL)
				assert.Equal(t, "application/json", reqs[0].Headers["accept"])
			},
		},
		{
			name: "带请求体的 POST 请求",
			input: `curl -X POST 'https://api.example.com/users' \
					-H 'content-type: application/json' \
					-d '{"name":"test"}'`,
			wantErr: false,
			check: func(t *testing.T, reqs []input.Request) {
				assert.Len(t, reqs, 1)
				assert.Equal(t, "POST", reqs[0].Method)
				assert.Equal(t, `{"name":"test"}`, reqs[0].Body)
			},
		},
		{
			name: "多个请求",
			input: `curl 'https://api1.example.com' \
					-H 'accept: application/json'
					curl 'https://api2.example.com' \
					-H 'accept: application/json'`,
			wantErr: false,
			check: func(t *testing.T, reqs []input.Request) {
				assert.Len(t, reqs, 2)
				assert.Equal(t, "https://api1.example.com", reqs[0].URL)
				assert.Equal(t, "https://api2.example.com", reqs[1].URL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqs, err := parser.Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			tt.check(t, reqs)
		})
	}
} 