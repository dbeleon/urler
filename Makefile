
compose_build:
	docker compose --file deployments/docker-compose.yaml build

up: compose_build # clean
	docker compose --file deployments/docker-compose.yaml up -d

down:
	docker compose --file deployments/docker-compose.yaml down

clean:
	rm -rf deployments/data

restart: down up