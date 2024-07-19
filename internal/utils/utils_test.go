package utils

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleRetryFunc() {
	err := RetryFunc(context.Background(), func() error {
		// ... do something
		return nil
	})
	if err != nil {
		fmt.Printf("error")
	}
}

func TestGetMD5Hash(t *testing.T) {
	type args struct {
		text string
	}
	type want struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Empty string",
			args: args{text: ""},
			want: want{text: "d41d8cd98f00b204e9800998ecf8427e"},
		},
		{
			name: "String with numbers",
			args: args{text: "1234567890"},
			want: want{text: "e807f1fcf82d132f9bb018ca6738a19f"},
		},
		{
			name: "JSON string",
			args: args{text: "{\"field1\":\"value1\",\"filed2\":\"value2\"}"},
			want: want{text: "8f4723681f45117c2457645d793f3f7e"},
		},
		{
			name: "Long string",
			args: args{text: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."},
			want: want{text: "01aad0e51fcd5582b307613842e4ffe5"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GetMD5Hash(test.args.text)
			assert.Equal(t, test.want.text, got)
		})
	}
}
