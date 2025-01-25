compose_build:
	docker compose --file deployments/docker-compose.yaml build

up: compose_build
	docker compose --file deployments/docker-compose.yaml up -d

cleanup: clean up

down:
	docker compose --file deployments/docker-compose.yaml down

clean:
	sudo rm -rf deployments/data

restart: down up

cleanrestart: down cleanup

test_add_user: TEST = add_user
test_add_user: USR = 200
test_add_user: DUR = 10s
test_add_user: test

test_make_url: TEST = make_url
test_make_url: USR = 200
test_make_url: DUR = 60s
test_make_url: test

test_make_url_same: TEST = make_url_same
test_make_url_same: USR = 1
test_make_url_same: DUR = 1s
test_make_url_same: test

test:
	mkdir -p k6/result/ k6/export && \
	if [ -f "k6/result/${TEST}.json" ] ; then rm -f k6/result/${TEST}.json ; fi && \
	if [ -f "k6/export/${TEST}.json" ] ; then rm -f k6/export/${TEST}.json ; fi && \
	k6 run -u ${USR} -d ${DUR} --summary-export=k6/export/${TEST}.json --out json=k6/result/${TEST}.json k6/${TEST}.js
