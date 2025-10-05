wire-gen :
	cd Route
	wire gen
	cd ..

build-prod :
	docker build -f dockerfile-prod -t debian-server.taile49fd9.ts.net:81/prod/ct-core:latest .
	docker push debian-server.taile49fd9.ts.net:81/prod/ct-core:latest

build-staging :
	docker build -f dockerfile-staging -t debian-server.taile49fd9.ts.net:81/staging/ct-core:latest .
	docker push debian-server.taile49fd9.ts.net:81/staging/ct-core:latest

generate-kube-config :
	kompose -f docker-compose.yml convert