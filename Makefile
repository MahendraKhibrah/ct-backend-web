wire-gen :
	cd Route
	wire gen
	cd ..

docker-build :
	docker compose build

docker-up :
	docker compose up -d --build --pull never

docker-down :
	docker compose down --remove-orphans


generate-kube-config :
	kompose -f docker-compose.yml convert