dcnet:
	docker network rm traefik_web
	docker network create  -d bridge traefik_web

dcupdb:
	docker-compose -f docker-compose.mydb.yaml -f docker-compose.traefik.yaml -f docker-compose.redis.yaml up -d

downdb:
	docker-compose -f docker-compose.mydb.yaml -f docker-compose.traefik.yaml -f docker-compose.redis.yaml down

rmnet:
	docker network rm traefik_web
	docker network rm dawnshine_default
	
dcapp:
	docker-compose up -d
