package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: check value
func TestLoadConfig(t *testing.T) {
	t.Run("Fail: due to wrong resource", func(t *testing.T) {
		conf := &Config{}
		err := LoadConfig("../../configs", "wrong_resource", conf)
		assert.Error(t, err)
	})

	t.Run("Fail: due to wrong path", func(t *testing.T) {
		conf := &Config{}
		err := LoadConfig("../../wrong_path", "wrong_resource", conf)
		assert.Error(t, err)
	})

	t.Run("Success: get config from service.env", func(t *testing.T) {
		conf := &Config{}
		err := LoadConfig("../../configs", "service", conf)
		assert.NoError(t, err)
		assert.Equal(t, "google_credential_path", conf.GoogleCredentialPath)
		assert.Equal(t, time.Second*30, conf.Timeout)
	})

	t.Run("Success: get config from nats.env", func(t *testing.T) {
		conf := &NatsConf{}
		err := LoadConfig("../../configs", "nats", conf)
		fmt.Println(conf)
		assert.NoError(t, err)
		assert.Equal(t, 10, conf.MaxReconnect)
		assert.Equal(t, time.Second*10, conf.ReconnectWait)
	})
}
