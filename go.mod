module github.com/layer5io/meshery-consul

go 1.16

replace (
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
)

require (
	github.com/layer5io/meshery-adapter-library v0.1.25
	github.com/layer5io/meshkit v0.2.36
	github.com/layer5io/service-mesh-performance v0.3.3
	gopkg.in/yaml.v2 v2.4.0
)
