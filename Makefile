.ONESHELL:
SHELL := /bin/bash


.PHONY: up
up:
	@docker compose up -d


# Core
.PHONY: restart
restart:
	@docker compose build core
	@docker compose up -d --force-recreate --no-deps core
	@$(MAKE) clear-cache

.PHONY: migrate
migrate:
	@docker exec -it core-core-1 ./tern migrate --migrations migrations


# Postgres
.PHONY: peek-db
peek-db:
	@docker exec -it core-postgres-1 psql -U causality


# Valkey
.PHONY: clear-cache
clear-cache:
	@docker exec -it core-valkey-1 redis-cli FLUSHALL


# Airflow
.PHONY: init-airflow
init-airflow:
	@docker exec -it core-airflow-scheduler-1 airflow db migrate

.PHONY: create-airflow-user
create-airflow-user:
	@docker exec -it core-airflow-scheduler-1 airflow users \
		create --role Admin --username admin --email admin@airflow.local \
		--firstname Causality --lastname Admin --password admin
