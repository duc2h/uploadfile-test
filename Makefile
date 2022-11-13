init: 
	docker compose down
	docker compose up

unit-test:
	go test -count=3 ./internals/...

call-api-with-normal-payload:
	curl -X POST http://localhost:8080/user/batch --data "@files/payload.json" -H "Content-Type: application/json"

call-api-with-heavy-payload:
	curl -X POST http://localhost:8080/user/batch --data "@files/payload-heavy.json" -H "Content-Type: application/json"