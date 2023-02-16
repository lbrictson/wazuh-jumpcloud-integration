package pkg

import (
	"reflect"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    ConfigurationData
		wantErr bool
	}{
		{
			name: "TestReadConfigFileHappyPath",
			args: args{path: "../test_data/example_config.json"},
			want: ConfigurationData{
				APIKey:  "this-is-not-a-real-key",
				BaseURL: "https://api.jumpcloud.com",
				Last:    nil,
				path:    "../test_data/example_config.json",
			},
		},
		{
			name: "TestReadConfigFileBadPath",
			args: args{
				path: "../test_data/not_a_real_file.json",
			},
			want:    ConfigurationData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfigFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(*got, tt.want) {
					t.Errorf("ReadConfigFile() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
