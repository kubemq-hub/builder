package manager

import (
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestConnectorsCatalog_initResourceNoLocalFile(t *testing.T) {
	file, err := ioutil.TempFile("./", "")
	defer func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}()
	require.NoError(t, err)
	cc := NewCatalogManager()
	resource, err := cc.initResource(file.Name(), targetRemoteUrl, targetRemoteHash, true)
	require.NoError(t, err)
	require.NotNil(t, resource)
	m, err := common.LoadManifest(resource)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, "targets", m.Schema)
}
func TestConnectorsCatalog_initResourceUpdate(t *testing.T) {
	file, err := ioutil.TempFile("./", "")
	defer func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}()
	require.NoError(t, err)
	cc := NewCatalogManager()
	resource, err := cc.initResource(file.Name(), targetRemoteUrl, targetRemoteHash, true)
	require.NoError(t, err)
	require.NotNil(t, resource)
	m, err := common.LoadManifest(resource)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, "targets", m.Schema)
	currentVersion := m.Version
	m.Version = "update"
	_ = file.Close()
	err = m.SaveFile(file.Name())
	require.NoError(t, err)
	resource, err = cc.initResource(file.Name(), targetRemoteUrl, targetRemoteHash, true)
	require.NoError(t, err)
	require.NotNil(t, resource)
	m, err = common.LoadManifest(resource)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, "targets", m.Schema)
	require.Equal(t, currentVersion, m.Version)
}
func TestConnectorsCatalog_initResourceNoHash(t *testing.T) {
	file, err := ioutil.TempFile("./", "")
	defer func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}()
	require.NoError(t, err)
	cc := NewCatalogManager()
	resource, err := cc.initResource(file.Name(), targetRemoteUrl, "http://localhost", true)
	require.NoError(t, err)
	require.NotNil(t, resource)
	m, err := common.LoadManifest(resource)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, "targets", m.Schema)

}
func TestConnectorsCatalog_initResourceNoAccess(t *testing.T) {
	file, err := ioutil.TempFile("./", "")
	defer func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}()
	require.NoError(t, err)
	cc := NewCatalogManager()
	resource, err := cc.initResource(file.Name(), "http://localhost", "http://localhost", true)
	require.NoError(t, err)
	require.NotNil(t, resource)
	m, err := common.LoadManifest(resource)
	require.Error(t, err)
	require.Nil(t, m)

}
