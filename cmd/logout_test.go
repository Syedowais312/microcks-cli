package cmd

import (
	"os"
	"testing"

	"github.com/microcks/microcks-cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const logoutTestConfig = `current-context: dev
contexts:
- name: dev
  server: http://localhost:8585
  user: http://localhost:8585
  instance: ""
- name: http://localhost:8585
  server: http://localhost:8585
  user: http://localhost:8585
  instance: ""
servers:
- name: microcks
  server: http://localhost:8585
  insecureTLS: true
  keycloakEnable: false
users:
- name: http://localhost:8585
  auth-token: token-value
  refresh-token: refresh-value
auths:
- server: http://localhost:8585
  clientid: ""
  clientsecret: ""
`

func TestLogoutContextNamedAlias(t *testing.T) {
	configPath := writeLogoutTestConfig(t)

	err := logoutContext("dev", configPath)
	require.NoError(t, err)

	localCfg, err := config.ReadLocalConfig(configPath)
	require.NoError(t, err)

	user, err := localCfg.GetUser("http://localhost:8585")
	require.NoError(t, err)
	assert.Empty(t, user.AuthToken)
	assert.Empty(t, user.RefreshToken)
}

func TestLogoutContextServerURL(t *testing.T) {
	configPath := writeLogoutTestConfig(t)

	err := logoutContext("http://localhost:8585", configPath)
	require.NoError(t, err)

	localCfg, err := config.ReadLocalConfig(configPath)
	require.NoError(t, err)

	user, err := localCfg.GetUser("http://localhost:8585")
	require.NoError(t, err)
	assert.Empty(t, user.AuthToken)
	assert.Empty(t, user.RefreshToken)
}

func writeLogoutTestConfig(t *testing.T) string {
	t.Helper()

	f, err := os.CreateTemp(t.TempDir(), "microcks-logout-config-*.yaml")
	require.NoError(t, err)

	_, err = f.WriteString(logoutTestConfig)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	require.NoError(t, os.Chmod(f.Name(), 0o600))

	return f.Name()
}
