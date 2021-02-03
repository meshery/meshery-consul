module github.com/layer5io/meshery-consul

go 1.15

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f

require (
	github.com/kumarabd/gokit v0.2.0 // indirect
	github.com/layer5io/meshery-adapter-library v0.1.12
	github.com/layer5io/meshkit v0.2.1-0.20210127211805-88e99ca45457
	google.golang.org/grpc v1.32.0 // indirect
	helm.sh/helm/v3 v3.3.1
)
