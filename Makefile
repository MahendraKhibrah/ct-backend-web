wire-gen :
	cd Route
	wire gen
	cd ..

build-prod :
	docker build -f dockerfile-prod -t 10.0.0.2:81/prod/ct-core:latest .
	docker push 10.0.0.2:81/prod/ct-core:latest

build-staging :
	docker build -f dockerfile-staging -t 10.0.0.2:81/staging/ct-core:latest .
	docker push 10.0.0.2:81/staging/ct-core:latest

generate-kube-config :
	kompose -f docker-compose.yml convert