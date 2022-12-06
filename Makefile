.DEFAULT_GOAL := help

deploy: ## Deploy all
	@make -s app-deploy
	@make -s nginx-deploy

app-deploy: ## Deploy Go App
	@sudo "Restart Go App"
	@sudo systemctl daemon-reload
	@sudo systemctl restart GoEarth.service
	@echo "Done: app-deploy"

nginx-deploy: ## Deploy nginx
	@sudo "Replace conf files"
	@sudo cp nginx/nginx.conf /etc/nginx/nginx.conf
	@sudo cp nginx/sites-available/earth /etc/nginx/sites-available/earth
	@echo "Validate nginx.conf"
	@sudo nginx -t
	@sudo "Restart nginx"
	@sudo systemctl daemon-reload
	@sudo systemctl restart nginx
	@echo "Done: nginx-deploy"

nginx-log: ## Tail nginx access.log
	@sudo tail -f /var/log/nginx/access.log

nginx-error-log: ## Tail nginx error.log
	@sudo tail -f /var/log/nginx/error.log

log: ## Tail journalctl
	@sudo journalctl -f

.PHONY: help
help:
	@grep -E '^[a-z0-9A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'