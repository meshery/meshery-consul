module github.com/layer5io/meshery-consul

go 1.15

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f

require (
	github.com/layer5io/meshery-adapter-library v0.1.13
	github.com/layer5io/meshkit v0.2.6
	github.com/layer5io/service-mesh-performance v0.3.2
	helm.sh/helm/v3 v3.3.1
)
