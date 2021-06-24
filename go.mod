module github.com/layer5io/meshery-consul

go 1.15

replace (
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
)

require (
	github.com/layer5io/meshery-adapter-library v0.1.20
	github.com/layer5io/meshkit v0.2.14
	github.com/layer5io/service-mesh-performance v0.3.3
	helm.sh/helm/v3 v3.3.1
)
