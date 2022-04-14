.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: start
start: ## start full stack inside docker environment
	docker-compose up --build --abort-on-container-exit --timeout 0

.PHONY: stop
stop: ## stop docker containers without cleaning up
	docker-compose stop

.PHONY: setup
setup: ## install development dependencies
	npm --prefix ./frontend install
	poetry install
	docker-compose build

.PHONY: frontend
frontend: ## start fronted server
	npm --prefix ./frontend start

.PHONY: backend_py
backend_py: ## start backend sever with hot reload
	cd ./backend_py && poetry run uvicorn main:app --port 3000 --reload

.PHONY: backend_go
backend_go: ## start backend_go sever with hot reload
	cd ./backend_go && fresh

.PHONY: keycloak
keycloak: ## star keycloak service with demo realm in docker
	docker-compose up keycloak

.PHONY: backup
backup: ## Export 'myrealm' realm to ke/realm-export.json (WARNING: requires keycloak to be running; the script don't exit after export you should kill them manually)
	docker-compose exec keycloak /opt/jboss/keycloak/bin/standalone.sh \
-Djboss.socket.binding.port-offset=100 -Dkeycloak.migration.action=export \
-Dkeycloak.migration.provider=singleFile \
-Dkeycloak.migration.realmName=myrealm \
-Dkeycloak.migration.usersExportStrategy=REALM_FILE \
-Dkeycloak.migration.file=/tmp/realm-export.json

.PHONY: clean
clean: ## Cleanup docker and volumes
	docker-compose down --remove-orphans --rmi all --volumes --timeout 0
