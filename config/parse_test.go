package config

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/caddyserver/caddy"
)

func TestParse(t *testing.T) {
	successData, err := ioutil.ReadFile("testdata/success")
	if err != nil {
		t.Fatal(err)
	}

	dynamicConf, err := ioutil.ReadFile("testdata/dynamic-conf")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		c *caddy.Controller
	}
	tests := []struct {
		name    string
		args    args
		wantCfg config
		wantErr bool
	}{
		{
			name: "success",
			args: args{c: caddy.NewTestController("dns", string(successData))},
			wantCfg: config{
				dynamicConfigPath:    "testdata/dynamic-conf",
				dynamicConfigContent: bytes.NewReader(dynamicConf),
				key:                  "",
				names:                []string{"p1", "p2"},
				paths: map[string]string{
					"p1": "p1.so",
					"p2": "p2.so",
				},
				setups: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := parse(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("parse() gotCfg = %v, \nwant %v", gotCfg, tt.wantCfg)
			}
		})
	}
}
