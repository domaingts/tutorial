package option

import "github.com/sagernet/sing/common/json/badoption"

type ExperimentalOptions struct {
	CacheFile *CacheFileOptions `json:"cache_file,omitempty"`
	Debug     *DebugOptions     `json:"debug,omitempty"`
}

type CacheFileOptions struct {
	Enabled     bool               `json:"enabled,omitempty"`
	Path        string             `json:"path,omitempty"`
	CacheID     string             `json:"cache_id,omitempty"`
	StoreFakeIP bool               `json:"store_fakeip,omitempty"`
	StoreRDRC   bool               `json:"store_rdrc,omitempty"`
	RDRCTimeout badoption.Duration `json:"rdrc_timeout,omitempty"`
}
