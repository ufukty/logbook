all: ../app/config/deployment/config.go

../app/config/deployment/config.go: testing.yml
	gonfique -in testing.yml -out ../app/config/deployment/config.go -pkg deployment