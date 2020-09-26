package config

var (
	ServerDefaults = map[string]string{
		"name":     "consul-adapter",
		"port":     "10002",
		"traceurl": "none",
		"version":  "v0.1.0",
	}

	MeshDefaults = map[string]string{
		"name":    "consul",
		"status":  "not installed", // TODO: status should be a type / an enum
		"version": "1.8.2",
	}

	ViperDefaults = map[string]string{
		"filepath": "~/.meshery-consul",
		"filename": "config.yml",
		"filetype": "yaml",
	}
)
