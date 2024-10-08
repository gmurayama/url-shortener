GOMOCK = go run go.uber.org/mock/mockgen@v0.4.0

.PHONY: mocks run
mocks:
	$(GOMOCK) -source=application/interfaces.go -destination=mocks/application.go -package=mocks

run:
	go run gateways/api/main.go
