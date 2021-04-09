module github.com/layer5io/meshery-consul

go 1.15

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f

require (
	github.com/layer5io/meshery-adapter-library v0.1.14
	github.com/layer5io/meshkit v0.2.7
	github.com/layer5io/service-mesh-performance v0.3.3
	helm.sh/helm/v3 v3.3.1
)
