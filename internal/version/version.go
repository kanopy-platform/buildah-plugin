package version

const (
	PluginTypeDrone = "drone"
)

var (
	// these are set using ldflags
	version    = "0.0.0"
	gitCommit  = ""
	pluginType = PluginTypeDrone
)

type Info struct {
	Version    string `json:"version,omitempty"`
	GitCommit  string `json:"gitCommit,omitempty"`
	PluginType string `json:"pluginType,omitempty"`
}

func Get() Info {
	return Info{
		Version:    version,
		GitCommit:  gitCommit,
		PluginType: pluginType,
	}
}
