wire-gen :
	cd Route
	wire gen
	cd ..

build-prod :
	docker build -f Dockerfile-prod -t 10.0.0.2:81/prod/ct-backend:latest .
	docker push 10.0.0.2:81/prod/ct-backend:latest

build-dev :
	docker build -f Dockerfile-dev -t 10.0.0.2:81/dev/ct-backend:latest .
	docker push 10.0.0.2:81/dev/ct-backend:latest

generate-kube-config :
	kompose -f docker-compose.yml convert