devup:
	docker-compose -p txstat -f./deploy/dev/docker-compose.yml up -d --build

devdown:
	docker-compose -p txstat -f./deploy/dev/docker-compose.yml down

test:
	go test ./internal/...
