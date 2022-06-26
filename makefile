all: run 

run:
	go run main.go

psql:
	docker run -it --name psql -d -p 5432:5432 \
	-e POSTGRES_USER=ui_test \
	-e POSTGRES_PASSWORD=ui_test \
	postgres 