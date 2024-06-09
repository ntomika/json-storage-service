run:
	docker build -t json-storage-service:latest .
	kubectl apply -f internal/configurations/deployment.yaml
	kubectl apply -f internal/configurations/service.yaml
	kubectl apply -f internal/configurations/roles.yaml

restart:
	docker exec $$(docker ps -q) killall main
	docker restart $$(docker ps -q)

clean:
	docker system prune -af
	kubectl delete all --all