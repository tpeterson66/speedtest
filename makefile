run:
	docker-compose --env-file .env up -d

destroy:
	docker-compose down