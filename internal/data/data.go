// Package data loads embedded detection pattern tables.
package data

import (
	_ "embed"
	"encoding/json"
	"sync"
)

//go:embed patterns.json
var patternsJSON []byte

// BotPattern matches known bots and similar clients.
type BotPattern struct {
	Pattern     string `json:"pattern"`
	Name        string `json:"name"`
	DeviceClass string `json:"deviceClass"`
	AgentClass  string `json:"agentClass"`
	Heuristic   bool   `json:"heuristic"`
	Mobile      *bool  `json:"mobile"`
}

// AppPattern matches in-app browsers and applications.
type AppPattern struct {
	Pattern     string `json:"pattern"`
	Name        string `json:"name"`
	AgentClass  string `json:"agentClass"`
	DeviceClass string `json:"deviceClass"`
	Product     string `json:"product"`
}

// BrandPattern matches device brand/model hints in the UA.
type BrandPattern struct {
	Prefix  string `json:"prefix"`
	Pattern string `json:"pattern"`
	Brand   string `json:"brand"`
	Name    string `json:"name"`
	Class   string `json:"class"`
}

// Catalog is the loaded pattern set.
type Catalog struct {
	Bots   []BotPattern   `json:"bots"`
	Apps   []AppPattern   `json:"apps"`
	Brands []BrandPattern `json:"brands"`
}

var (
	loadOnce sync.Once
	catalog  Catalog
	loadErr  error
)

// Load returns the embedded pattern catalog.
func Load() (Catalog, error) {
	loadOnce.Do(func() {
		loadErr = json.Unmarshal(patternsJSON, &catalog)
	})
	return catalog, loadErr
}

// MustLoad returns the catalog or panics.
func MustLoad() Catalog {
	c, err := Load()
	if err != nil {
		panic("useragent data: " + err.Error())
	}
	return c
}
