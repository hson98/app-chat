dockerdb:
	docker-compose --env-file ./app.env  up -d db cache
mock:
	mockgen -package mockauth -destination internal/auth/mock/pg_repository_mock.go  hson98/app-chat/internal/auth/repository Repository