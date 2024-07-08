package gvar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_prepareWelcome(t *testing.T) {
	type args struct {
		buildVersion string
		buildDate    string
		buildCommit  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All variable empty",
			args: args{
				buildVersion: "",
				buildDate:    "",
				buildCommit:  "",
			},
			want: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A",
		},
		{
			name: "buildVersion variable specified",
			args: args{
				buildVersion: "0.0.1",
				buildDate:    "",
				buildCommit:  "",
			},
			want: "Build version: 0.0.1\nBuild date: N/A\nBuild commit: N/A",
		},
		{
			name: "buildVersion and buildDate variables specified",
			args: args{
				buildVersion: "0.0.1",
				buildDate:    "01.01.1970 00:00:00",
				buildCommit:  "",
			},
			want: "Build version: 0.0.1\nBuild date: 01.01.1970 00:00:00\nBuild commit: N/A",
		},
		{
			name: "all variables specified",
			args: args{
				buildVersion: "0.0.1",
				buildDate:    "01.01.1970 00:00:00",
				buildCommit:  "super build",
			},
			want: "Build version: 0.0.1\nBuild date: 01.01.1970 00:00:00\nBuild commit: super build",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args.buildVersion != "" {
				BuildVersion = test.args.buildVersion
			}
			if test.args.buildDate != "" {
				BuildDate = test.args.buildDate
			}
			if test.args.buildCommit != "" {
				BuildCommit = test.args.buildCommit
			}

			got := prepareWelcome()
			assert.Equal(t, test.want, got)
		})
	}
}
