package filefinder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFiles(t *testing.T) {
	type args struct {
		p          string
		endPattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"Directory as input",
			args{
				p:          "test/",
				endPattern: "**/*.tfstate",
			},
			[]string{"test/test.tfstate", "test/sub/terraform.tfstate", "test/sub/subsub/terraform.tfstate", "test/terraform.tfstate/terraform.tfstate"},
			false,
		},
		{
			"**/*.tfstate as input",
			args{
				p:          "test/**/*.tfstate",
				endPattern: "**/*.tfstate",
			},
			[]string{"test/test.tfstate", "test/sub/terraform.tfstate", "test/sub/subsub/terraform.tfstate", "test/terraform.tfstate/terraform.tfstate"},
			false,
		},
		{
			"*.tfstate as input",
			args{
				p:          "test/*.tfstate",
				endPattern: "**/*.tfstate",
			},
			[]string{"test/test.tfstate"},
			false,
		},
		{
			"test.tfstate as input",
			args{
				p:          "test/test.tfstate",
				endPattern: "**/*.tfstate",
			},
			[]string{"test/test.tfstate"},
			false,
		},
		{
			"**/*.notexist as input",
			args{
				p:          "test/**/*.notexist",
				endPattern: "**/*.tfstate",
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindFiles(tt.args.p, tt.args.endPattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
