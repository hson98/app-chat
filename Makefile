dockerdb:
	docker-compose --env-file ./app.env  up -d db cache