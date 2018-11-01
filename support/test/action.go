package test

import (
	"github.com/project-flogo/core/app/resource"
)

func NewActionInitCtx() *ActionInitCtx {
	resources := make(map[string]*resource.Resource)
	manager := resource.NewManager(resources)
	return &ActionInitCtx{manager: manager, resources: resources}
}

type ActionInitCtx struct {
	resources map[string]*resource.Resource
	manager   *resource.Manager
}

func (ctx *ActionInitCtx) ResourceManager() *resource.Manager {
	return ctx.manager
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
