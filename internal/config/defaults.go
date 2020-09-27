package config

var (
	ServerDefaults = map[string]string{
		"name":     "consul-adapter",
		"port":     "10002",
		"traceurl": "none",
		"version":  "v0.1.0",
	}

	MeshSpecDefaults = map[string]string{
		"name":    "Consul",
		"status":  "not installed", // TODO: status should be a type / an enum
		"version": "1.8.2",
	}

	MeshInstanceDefaults = map[string]string{
		"name":    "Consul",
		"status":  "not installed", // TODO: status should be a type / an enum
		"version": "1.8.2",
	}

	ViperDefaults = map[string]string{
		"filepath": "~/.meshery-consul",
		"filename": "config.yml",
		"filetype": "yaml",
	}
)
