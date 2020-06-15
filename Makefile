################################################################################
ifeq (, $(shell which buildtool))
$(warning "Installing buildtool...")
$(warning "go get github.com/gearboxworks/buildtool")
$(shell go get github.com/gearboxworks/buildtool)
endif
BUILDTOOL := $(shell which buildtool)
ifeq (, $(BUILDTOOL))
$(error "No buildtool found...")
endif
################################################################################

all:
	@echo ""
	@echo "build		- Build for local testing."
	@echo "release		- Build for published release."
	@echo "push		- Push repo to GitHub."
	@echo ""
	@$(BUILDTOOL) get all

build:
	@make pkgreflect
	@$(BUILDTOOL) build

release:
	@make pkgreflect
	@$(BUILDTOOL) release

push:
	@make pkgreflect
	@$(BUILDTOOL) push

pkgreflect:
	@$(BUILDTOOL) pkgreflect .

