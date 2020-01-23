package test

import (
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/support/service"
)

func NewActionInitCtx() *ActionInitCtx {
	resources := make(map[string]*resource.Resource)
	resManager := resource.NewManager(resources)
	svcManager := service.NewServiceManager()
	return &ActionInitCtx{resManager: resManager, svcManager:svcManager, resources: resources}
}

type ActionInitCtx struct {
	resources  map[string]*resource.Resource
	resManager *resource.Manager
	svcManager *service.Manager
}

func (ctx *ActionInitCtx) ServiceManager() *service.Manager {
	return ctx.svcManager
}

func (ctx *ActionInitCtx) RuntimeSettings() map[string]interface{} {
	return make(map[string]interface{})
}

func (ctx *ActionInitCtx) ResourceManager() *resource.Manager {
	return ctx.resManager
}

func (ctx *ActionInitCtx) AddResource(resourceType string, config *resource.Config) error {
	loader := resource.GetLoader(resourceType)
	res, err := loader.LoadResource(config)
	if err != nil {
		return err
	}
	ctx.resources[config.ID] = res

	return nil
}
