up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose up -d --build

logs:
	docker-compose logs --tail 1000 -f app

bash:
	docker-compose exec app sh
