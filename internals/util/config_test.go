package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: check value
func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name     string
		resource string
		path     string
		check    func(*testing.T, *Config, error)
	}{
		{
			name:     "Fail: wrong path",
			resource: "test",
			path:     "wrong_path",
			check: func(t *testing.T, conf *Config, err error) {
				assert.Error(t, err)
				assert.Equal(t, "Config File \"app.test\" Not Found in \"[/home/duchh/Desktop/edarha/uploadfile-test/internals/util/wrong_path]\"", err.Error())
				assert.Nil(t, conf)
			},
		},
		{
			name:     "Fail: wrong env",
			resource: "wrong_env",
			path:     "../../configs",
			check: func(t *testing.T, conf *Config, err error) {
				assert.Error(t, err)
				assert.Equal(t, "Config File \"app.wrong_env\" Not Found in \"[/home/duchh/Desktop/edarha/uploadfile-test/configs]\"", err.Error())
				assert.Nil(t, conf)
			},
		},
		{
			name:     "Success: load config success",
			resource: "test",
			path:     "../../configs",
			check: func(t *testing.T, conf *Config, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, conf)
				assert.Equal(t, "google_project_id", conf.GoogleProjectID)
				assert.Equal(t, "aws_secret_key", conf.AwsSecretKey)
				assert.Equal(t, time.Second*30, conf.Timeout)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// os.Setenv("ENVIRONMENT", tc.env)
			conf := &Config{}
			err := LoadConfig(tc.path, tc.resource, conf)
			tc.check(t, conf, err)
		})
	}
}
