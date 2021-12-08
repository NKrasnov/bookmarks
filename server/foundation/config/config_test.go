package config

import (
	"reflect"
	"testing"
)

func TestParseCommandLineParameters(t *testing.T) {
	tt := []struct {
		name   string
		get    []string
		expect map[string]string
	}{
		{name: "tc1", get: []string{"", "-host=127.0.0.1", "port=8080"}, expect: map[string]string{"host": "127.0.0.1", "port": "8080"}},
		{name: "tc2", get: []string{"", "-host=127.0.0.1", "-port=8080"}, expect: map[string]string{"host": "127.0.0.1", "port": "8080"}},
		{name: "tc3", get: []string{"", "-host=", "port=8080"}, expect: map[string]string{"host": "127.0.0.1", "port": "8080"}},
		{name: "tc4", get: []string{"", "host=", "-port="}, expect: map[string]string{"host": "127.0.0.1", "port": "8080"}},
		{name: "tc5", get: []string{"", "-host=127.0.0.1"}, expect: map[string]string{"host": "127.0.0.1"}},
	}

	for _, tc := range tt {
		t.Run(
			tc.name,
			func(t *testing.T) {
				res, err := parseCommandLineParameters(tc.get)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(res, tc.expect) {
					t.Errorf("Name: %s\nexpected: %v\n, got:%v", tc.name, tc.expect, res)
				}
			})
	}
}
