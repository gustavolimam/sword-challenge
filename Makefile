build:
	mkdir -p etc/rabbitmq/conf	
	sudo cp rabbitmq.conf ./etc/rabbitmq/conf	
	docker-compose up -d