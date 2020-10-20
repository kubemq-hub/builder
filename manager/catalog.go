package manager

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/kubemq-hub/builder/pkg/utils"
	"github.com/kubemq-hub/builder/survey"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	sourceLocalFile  = "./sources-manifest.json"
	sourceRemoteUrl  = "https://raw.githubusercontent.com/kubemq-hub/kubemq-sources/master/sources-manifest.json"
	sourceRemoteHash = "https://raw.githubusercontent.com/kubemq-hub/kubemq-sources/master/sources-manifest-hash.txt"

	targetLocalFile  = "./target-manifest.json"
	targetRemoteUrl  = "https://raw.githubusercontent.com/kubemq-hub/kubemq-targets/master/targets-manifest.json"
	targetRemoteHash = "https://raw.githubusercontent.com/kubemq-hub/kubemq-targets/master/targets-manifest-hash.txt"
)

type ConnectorsCatalog struct {
	SourcesManifest []byte
	TargetsManifest []byte
}

func NewConnectorCatalog() *ConnectorsCatalog {
	return &ConnectorsCatalog{}
}

func (cc *ConnectorsCatalog) loadFromFile(filename string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return data
}
func (cc *ConnectorsCatalog) saveToFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0600)
}

func (cc *ConnectorsCatalog) loadFromUrl(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var data []byte
	buff := bytes.NewBuffer(data)
	_, err = io.Copy(buff, resp.Body)
	if err != nil {
		return nil
	}
	return buff.Bytes()
}
func (cc *ConnectorsCatalog) initResource(localFile, remoteUrl, remoteHash string) ([]byte, error) {
	localData := cc.loadFromFile(localFile)
	remoteHashStringData := cc.loadFromUrl(remoteHash)
	var remoteHashData []byte
	if remoteHashStringData != nil {
		remoteHashData, _ = hex.DecodeString(string(remoteHashStringData))
	}
	if !hashIsEqual(localData, remoteHashData) {
		if remoteUrlData := cc.loadFromUrl(remoteUrl); remoteUrlData != nil {
			err := cc.saveToFile(localFile, remoteUrlData)
			if err != nil {
				return nil, err
			}
			return remoteUrlData, nil
		}
	}
	if localData == nil {
		return nil, fmt.Errorf("no resource found")
	}
	return localData, nil
}
func (cc *ConnectorsCatalog) init() error {

	var err error
	if cc.SourcesManifest, err = cc.initResource(sourceLocalFile, sourceRemoteUrl, sourceRemoteHash); err != nil {
		return fmt.Errorf("error loading sources connector catalog,%s", err.Error())
	}
	if cc.SourcesManifest, err = cc.initResource(targetLocalFile, targetRemoteUrl, targetRemoteHash); err != nil {
		return fmt.Errorf("error loading targets connector catalog,%s", err.Error())
	}

	return nil
}

func (cc *ConnectorsCatalog) browseTargets() error {
	return nil
}
func (cc *ConnectorsCatalog) browseSources() error {
	return nil
}
func (cc *ConnectorsCatalog) updateCatalog() error {
	utils.Println(promptCatalogLoadingStarted)
	err := cc.init()
	if err != nil {
		utils.Println(promptCatalogLoadingError, err.Error())
		return err
	}
	utils.Println(promptCatalogLoadingCompleted)
	return nil
}
func (cc *ConnectorsCatalog) Render() error {
	if cc.SourcesManifest == nil || cc.TargetsManifest == nil {
		if err := cc.init(); err != nil {
			return err
		}
	}
	if err := survey.NewMenu("Connectors Catalog Management: Select operation").
		AddItem("Browse Targets Catalog", cc.browseTargets).
		AddItem("Browse Sources Catalog", cc.browseSources).
		AddItem("Update Catalogs", cc.updateCatalog).
		AddItem("<-back", nil).
		Render(); err != nil {
		return err
	}
	return nil
}
func hashIsEqual(a, b []byte) bool {
	return a != nil && b != nil && hash(a) == hash(b)
}

func hash(data []byte) string {
	if data == nil {
		return ""
	}
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
