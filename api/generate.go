package api

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"regexp"
	"strings"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data/expression/function"
	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/core/data/schema"
	"github.com/project-flogo/core/support"
)

// Import is a package import
type Import struct {
	Alias  string
	Import string
	Used   bool
}

// GetAlias gets the import alias and marks it as used
func (i *Import) GetAlias() string {
	i.Used = true
	return i.Alias
}

// Imports are the package imports
type Imports struct {
	Imports []*Import
}

// Ensure looks up an import and adds it if it is missing
func (i *Imports) Ensure(path string, name ...string) *Import {
	if strings.HasPrefix(path, "#") {
		alias := strings.TrimPrefix(path, "#")
		for _, port := range i.Imports {
			if port.Alias == alias {
				return port
			}
		}
		panic(fmt.Errorf("ref %s not found", path))
	}
	for _, port := range i.Imports {
		if port.Import == path {
			return port
		}
	}
	parts := strings.Split(path, "/")
	alias := parts[len(parts)-1]
	if len(name) == 1 {
		alias = name[0]
	}
	if alias != "_" {
		for _, port := range i.Imports {
			if port.Alias == alias {
				alias = fmt.Sprintf("port%d", len(i.Imports))
				break
			}
		}
	}
	port := &Import{
		Alias:  alias,
		Import: path,
	}
	i.Imports = append(i.Imports, port)
	return port
}

type generator struct {
	resManager *resource.Manager
	imports    Imports
}

func (g *generator) ResourceManager() *resource.Manager {
	return g.resManager
}

var flogoImportPattern = regexp.MustCompile(`^(([^ ]*)[ ]+)?([^@:]*)@?([^:]*)?:?(.*)?$`) // extract import path even if there is an alias and/or a version

// Generator generates code for an action
type Generator interface {
	Generate(settingsName string, imports *Imports, config *action.Config) (code string, err error)
}

// GenerateResource is used to determine if a resource is generated, defaults to true
type GenerateResource interface {
	Generate() bool
}

// Generate generates flogo go API code
func Generate(config *app.Config, file string) {
	if config.Type != "flogo:app" {
		panic("invalid app type")
	}

	app := generator{}

	for _, anImport := range config.Imports {
		matches := flogoImportPattern.FindStringSubmatch(anImport)
		alias, ref := matches[1], matches[3]
		var port *Import
		if alias == "" {
			port = app.imports.Ensure(ref)
		} else {
			port = app.imports.Ensure(ref, alias)
		}

		for _, typ := range [...]string{"activity", "action", "trigger", "function", "other"} {
			err := support.RegisterAlias(typ, port.Alias, port.Import)
			if err != nil {
				panic(err)
			}
		}

		function.SetPackageAlias(port.Import, port.Alias)
	}

	function.ResolveAliases()

	for id, def := range config.Schemas {
		_, err := schema.Register(id, def)
		if err != nil {
			panic(err)
		}
	}

	schema.ResolveSchemas()

	properties := make(map[string]interface{}, len(config.Properties))
	for _, attr := range config.Properties {
		properties[attr.Name()] = attr.Value()
	}

	propertyManager := property.NewManager(properties)
	property.SetDefaultManager(propertyManager)

	resources := make(map[string]*resource.Resource, len(config.Resources))
	app.resManager = resource.NewManager(resources)

	for _, actionFactory := range action.Factories() {
		err := actionFactory.Initialize(&app)
		if err != nil {
			panic(err)
		}
	}

	output := "/*\n"
	output += fmt.Sprintf("* Name: %s\n", config.Name)
	output += fmt.Sprintf("* Type: %s\n", config.Type)
	output += fmt.Sprintf("* Version: %s\n", config.Version)
	output += fmt.Sprintf("* Description: %s\n", config.Description)
	output += fmt.Sprintf("* AppModel: %s\n", config.AppModel)
	output += "*/\n\n"

	errorCheck := func() {
		output += "if err != nil {\n"
		output += "panic(err)\n"
		output += "}\n"
	}

	output += "func main() {\n"
	output += "var err error\n"
	port := app.imports.Ensure("github.com/project-flogo/core/api")
	output += fmt.Sprintf("app := %s.NewApp()\n", port.GetAlias())

	for i, resConfig := range config.Resources {
		resType, err := resource.GetTypeFromID(resConfig.ID)
		if err != nil {
			panic(err)
		}

		loader := resource.GetLoader(resType)
		res, err := loader.LoadResource(resConfig)
		if err != nil {
			panic(err)
		}

		generate := true
		if g, ok := loader.(GenerateResource); ok {
			generate = g.Generate()
		}
		if generate {
			port := app.imports.Ensure("encoding/json")
			output += fmt.Sprintf("resource%d := %s.RawMessage(`%s`)\n", i, port.GetAlias(), string(resConfig.Data))
			output += fmt.Sprintf("app.AddResource(\"%s\", resource%d)\n", resConfig.ID, i)
		}

		resources[resConfig.ID] = res
	}

	if len(config.Properties) > 0 {
		port := app.imports.Ensure("github.com/project-flogo/core/data")
		for _, property := range config.Properties {
			output += fmt.Sprintf("app.AddProperty(\"%s\", %s.%s, %#v)\n", property.Name(), port.GetAlias(),
				property.Type().Name(), property.Value())
		}
	}
	if len(config.Channels) > 0 {
		port := app.imports.Ensure("github.com/project-flogo/core/engine/channels")
		for i, channel := range config.Channels {
			if i == 0 {
				output += fmt.Sprintf("name, buffSize := %s.Decode(\"%s\")\n", port.GetAlias(), channel)
			} else {
				output += fmt.Sprintf("name, buffSize = %s.Decode(\"%s\")\n", port.GetAlias(), channel)
			}
			output += fmt.Sprintf("_, err = %s.New(name, buffSize)\n", port.GetAlias())
			errorCheck()
		}
	}
	for i, act := range config.Actions {
		port := app.imports.Ensure(act.Ref)
		factory, settingsName := action.GetFactory(act.Ref), fmt.Sprintf("actionSettings%d", i)
		if generator, ok := factory.(Generator); ok {
			code, err := generator.Generate(settingsName, &app.imports, act)
			if err != nil {
				panic(err)
			}
			output += "\n"
			output += code
			output += "\n"
		} else {
			output += fmt.Sprintf("%s := %#v\n", settingsName, act.Settings)
		}
		output += fmt.Sprintf("err = app.AddAction(\"%s\", &%s.Action{}, %s)\n", act.Id, port.GetAlias(), settingsName)
		errorCheck()
	}
	for i, trigger := range config.Triggers {
		port := app.imports.Ensure(trigger.Ref)
		output += fmt.Sprintf("trg%d := app.NewTrigger(&%s.Trigger{}, %#v)\n", i, port.GetAlias(), trigger.Settings)
		for j, handler := range trigger.Handlers {
			output += fmt.Sprintf("handler%d_%d, err := trg%d.NewHandler(%#v)\n", i, j, i, handler.Settings)
			errorCheck()
			for k, act := range handler.Actions {
				if act.Id != "" {
					output += fmt.Sprintf("action%d_%d_%d, err := handler%d_%d.NewAction(\"%s\")\n", i, j, k, i, j, act.Id)
				} else {
					port := app.imports.Ensure(act.Ref)
					factory, settingsName := action.GetFactory(act.Ref), fmt.Sprintf("settings%d_%d_%d", i, j, k)
					if generator, ok := factory.(Generator); ok {
						code, err := generator.Generate(settingsName, &app.imports, act.Config)
						if err != nil {
							panic(err)
						}
						output += "\n"
						output += code
						output += "\n"
					} else {
						output += fmt.Sprintf("%s := %#v\n", settingsName, act.Settings)
					}
					output += fmt.Sprintf("action%d_%d_%d, err := handler%d_%d.NewAction(&%s.Action{}, %s)\n", i, j, k, i, j, port.GetAlias(), settingsName)
				}
				errorCheck()
				if act.If != "" {
					output += fmt.Sprintf("action%d_%d_%d.SetCondition(\"%s\")\n", i, j, k, act.If)
				}
				if length := len(act.Input); length > 0 {
					mappings := make([]string, 0, length)
					for key, value := range act.Input {
						mappings = append(mappings, fmt.Sprintf("%s%v", key, value))
					}
					output += fmt.Sprintf("action%d_%d_%d.SetInputMappings(%#v...)\n", i, j, k, mappings)
				}
				if length := len(act.Output); length > 0 {
					mappings := make([]string, 0, length)
					for key, value := range act.Output {
						mappings = append(mappings, fmt.Sprintf("%s%v", key, value))
					}
					output += fmt.Sprintf("action%d_%d_%d.SetOutputMappings(%#v...)\n", i, j, k, mappings)
				}
				output += fmt.Sprintf("_ = action%d_%d_%d\n", i, j, k)
			}
			output += fmt.Sprintf("_ = handler%d_%d\n", i, j)
		}
		output += fmt.Sprintf("_ = trg%d\n", i)
	}
	port = app.imports.Ensure("github.com/project-flogo/core/api")
	output += fmt.Sprintf("e, err := %s.NewEngine(app)\n", port.GetAlias())
	errorCheck()
	port = app.imports.Ensure("github.com/project-flogo/core/engine")
	output += fmt.Sprintf("%s.RunEngine(e)\n", port.GetAlias())
	output += "}\n"

	header := "package main\n\n"
	header += "import (\n"
	for _, port := range app.imports.Imports {
		if port.Used {
			header += fmt.Sprintf("%s \"%s\"\n", port.Alias, port.Import)
			continue
		}
		header += fmt.Sprintf("_ \"%s\"\n", port.Import)
	}
	header += ")\n"

	output = header + output

	out, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	buffer := bytes.NewBufferString(output)
	fileSet := token.NewFileSet()
	code, err := parser.ParseFile(fileSet, file, buffer, parser.ParseComments)
	if err != nil {
		buffer.WriteTo(out)
		panic(fmt.Errorf("%v: %v", file, err))
	}

	formatter := printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}
	err = formatter.Fprint(out, fileSet, code)
	if err != nil {
		buffer.WriteTo(out)
		panic(fmt.Errorf("%v: %v", file, err))
	}
}
