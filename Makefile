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
test_make_url: DUR = 20s
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



test_many_get_rare_make: TEST = many_get_rare_make
test_many_get_rare_make: test2

test_many_get_many_make: TEST = many_get_many_make
test_many_get_many_make: test2

test2: getshorts
	mkdir -p k6/result/ k6/export && \
	if [ -f "k6/result/${TEST}.json" ] ; then rm -f k6/result/${TEST}.json ; fi && \
	if [ -f "k6/export/${TEST}.json" ] ; then rm -f k6/export/${TEST}.json ; fi && \
	k6 run --summary-export=k6/export/${TEST}.json --out json=k6/result/${TEST}.json k6/${TEST}.js

urlstntconsole:
	tt connect 127.0.0.1:3303 -u admin -p dev

getshorts:
	curl -o k6/shorts.json http://localhost:8000/v1/shorts?limit=10000000&offset=0

docker-cmds:
	docker logs urler-2 -f 2>&1 | grep ERR -A 2
	docker logs urler-1 -n 1000
	docker logs nginx-balancer >& nginx.log
	docker logs qrer-1 -f
