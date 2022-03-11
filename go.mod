module github.com/layer5io/meshery-consul

go 1.17

replace (
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
)

require (
	github.com/layer5io/meshery-adapter-library v0.5.2
	github.com/layer5io/meshkit v0.5.6
	github.com/layer5io/service-mesh-performance v0.3.3
	gopkg.in/yaml.v2 v2.4.0
)
