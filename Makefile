include .bingo/Variables.mk

install-tool-deps: ## Installs dependencies for code generation
install-tool-deps: $(MOCKGEN) $(OAPI_CODEGEN)
	@echo ">>GOBIN=$(GOBIN)"