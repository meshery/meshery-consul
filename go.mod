module github.com/layer5io/meshery-consul

go 1.13

replace github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f

require (
	github.com/layer5io/meshery-adapter-library v0.1.3
	github.com/layer5io/meshkit v0.1.20
	google.golang.org/grpc v1.32.0 // indirect
)
