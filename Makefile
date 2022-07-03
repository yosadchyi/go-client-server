.PHONY: up down

up:
	(cd deployments && docker-compose up --build)

down:
	(cd deployments && docker-compose down)
