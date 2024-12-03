package curl

import (
	"reflect"
	"testing"
	"time"

	"monkey-test-api/internal/types"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []types.Request
		wantErr bool
	}{
		{
			name:    "空输入",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "非 curl 命令",
			input:   "wget http://example.com",
			want:    nil,
			wantErr: true,
		},
		{
			name: "基本 GET 请求",
			input: `curl 'http://example.com/api/test'`,
			want: []types.Request{
				{
					Method:  "GET",
					URL:     "http://example.com/api/test",
					Headers: map[string]string{},
					Body:    map[string]interface{}{},
					Params:  map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "带请求头的 POST 请求",
			input: `curl -X POST 'http://example.com/api/users' \
                    -H 'Content-Type: application/json' \
                    -H 'Authorization: Bearer token123' \
                    -d '{"name":"test","age":25}'`,
			want: []types.Request{
				{
					Method: "POST",
					URL:    "http://example.com/api/users",
					Headers: map[string]string{
						"Content-Type":  "application/json",
						"Authorization": "Bearer token123",
					},
					Body: map[string]interface{}{
						"name": "test",
						"age":  float64(25),
					},
					Params: map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "带查询参数的请求",
			input: `curl 'http://example.com/api/search?q=test&page=1'`,
			want: []types.Request{
				{
					Method:  "GET",
					URL:     "http://example.com/api/search",
					Headers: map[string]string{},
					Body:    map[string]interface{}{},
					Params: map[string]interface{}{
						"q":    "test",
						"page": "1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "多个 curl 命令",
			input: `curl 'http://example.com/api/users'
                    curl 'http://example.com/api/posts'`,
			want: []types.Request{
				{
					Method:  "GET",
					URL:     "http://example.com/api/users",
					Headers: map[string]string{},
					Body:    map[string]interface{}{},
					Params:  map[string]interface{}{},
				},
				{
					Method:  "GET",
					URL:     "http://example.com/api/posts",
					Headers: map[string]string{},
					Body:    map[string]interface{}{},
					Params:  map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "带续行符的命令",
			input: `curl 'http://example.com/api/users' \
                    -H 'Content-Type: application/json' \
                    -d '{"name":"test"}'`,
			want: []types.Request{
				{
					Method: "GET",
					URL:    "http://example.com/api/users",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"name": "test",
					},
					Params: map[string]interface{}{},
				},
			},
			wantErr: false,
		},
	}

	parser := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// 由于 RequestID 和 Timestamp 是动态生成的，我们需要特殊处理比较
			for i := range got {
				if got[i].RequestID == "" {
					t.Errorf("Parser.Parse() RequestID is empty")
				}
				if got[i].Timestamp.IsZero() {
					t.Errorf("Parser.Parse() Timestamp is zero")
				}

				// 清除动态字段以进行比较
				got[i].RequestID = ""
				got[i].Timestamp = time.Time{}
				tt.want[i].RequestID = ""
				tt.want[i].Timestamp = time.Time{}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
} 