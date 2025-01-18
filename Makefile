
compose_build:
	docker compose --file deployments/docker-compose.yaml build

up: compose_build
	docker compose --file deployments/docker-compose.yaml up -d

cleanup: clean up

down:
	docker compose --file deployments/docker-compose.yaml down

clean:
	rm -rf deployments/data

restart: down up

cleanrestart: down cleanup