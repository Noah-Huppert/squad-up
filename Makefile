.PHONY: db-create db-destroy db-start db-stop app-run
db-create:
	docker run \
		--name squad-up-postgres \
		-e POSTGRES_USER=username \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=squad-up \
		-p 5432:5432 \
		-d \
		postgres
db-destroy:
	docker rm squad-up-postgres
db-start:
	docker start squad-up-postgres
db-stop:
	docker stop squad-up-postgres

app-run:
	go run ./server/main.go