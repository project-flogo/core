package api

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"regexp"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/app"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data"
)

// Import is a package import
type Import struct {
	Alias  string
	Import string
}

// Imports are the package imports
type Imports struct {
	Imports []Import
}

// Add adds an import
func (i *Imports) Add(path string) Import {
	for _, port := range i.Imports {
		if port.Import == path {
			return port
		}
	}
	port := Import{
		Import: path,
	}
	i.Imports = append(i.Imports, port)
	return port
}

// AddWithAlias adds an import with an alias
func (i *Imports) AddWithAlias(alias, path string) Import {
	for _, port := range i.Imports {
		if port.Import == path {
			return port
		}
	}
	if alias == "" {
		alias = fmt.Sprintf("port%d", len(i.Imports))
	}
	port := Import{
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

var dataTypes = map[data.Type]string{
	data.TypeUnknown: "TypeUnknown",
	data.TypeAny:     "TypeAny",
	data.TypeString:  "TypeString",
	data.TypeInt:     "TypeInt",
	data.TypeInt32:   "TypeInt32",
	data.TypeInt64:   "TypeInt64",
	data.TypeFloat32: "TypeFloat32",
	data.TypeFloat64: "TypeFloat64",
	data.TypeBool:    "TypeBool",
	data.TypeObject:  "TypeObject",
	data.TypeBytes:   "TypeBytes",
	data.TypeParams:  "TypeParams",
	data.TypeArray:   "TypeArray",
	data.TypeMap:     "TypeMap",
}

// Generate generates flogo go API code
func Generate(config *app.Config, file string) {
	if config.Type != "flogo:app" {
		panic("invalid app type")
	}

	app := generator{}

	for _, anImport := range config.Imports {
		matches := flogoImportPattern.FindStringSubmatch(anImport)
		app.imports.AddWithAlias(matches[1], matches[3])
	}

	resources := make(map[string]*resource.Resource, len(config.Resources))
	app.resManager = resource.NewManager(resources)

	for _, actionFactory := range action.Factories() {
		err := actionFactory.Initialize(&app)
		if err != nil {
			panic(err)
		}
	}

	for _, resConfig := range config.Resources {
		resType, err := resource.GetTypeFromID(resConfig.ID)
		if err != nil {
			panic(err)
		}

		loader := resource.GetLoader(resType)
		res, err := loader.LoadResource(resConfig)
		if err != nil {
			panic(err)
		}

		resources[resConfig.ID] = res
	}

	output := "/*\n"
	output += fmt.Sprintf("* Name: %s\n", config.Name)
	output += fmt.Sprintf("* Type: %s\n", config.Type)
	output += fmt.Sprintf("* Version: %s\n", config.Version)
	output += fmt.Sprintf("* Description: %s\n", config.Description)
	output += fmt.Sprintf("* AppModel: %s\n", config.AppModel)
	output += "*/\n\n"

	output += "func main() {\n"
	output += "var err error\n"
	app.imports.Add("github.com/project-flogo/core/api")
	output += "app := api.NewApp()\n"
	if len(config.Properties) > 0 {
		app.imports.Add("github.com/project-flogo/core/data")
		for _, property := range config.Properties {
			output += fmt.Sprintf("app.AddProperty(%s, data.%s, %#v)\n", property.Name(), dataTypes[property.Type()], property.Value())
		}
	}
	for i, act := range config.Actions {
		port := app.imports.AddWithAlias("", act.Ref)
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
		output += fmt.Sprintf("err = app.AddAction(\"%s\", &%s.Action{}, %s)\n", act.Id, port.Alias, settingsName)
		output += "if err != nil {\n"
		output += "panic(err)\n"
		output += "}\n"
	}
	for i, trigger := range config.Triggers {
		port := app.imports.AddWithAlias("", trigger.Ref)
		output += fmt.Sprintf("trg%d := app.NewTrigger(&%s.Trigger{}, %#v)\n", i, port.Alias, trigger.Settings)
		for j, handler := range trigger.Handlers {
			output += fmt.Sprintf("handler%d_%d, err := trg%d.NewHandler(%#v)\n", i, j, i, handler.Settings)
			output += "if err != nil {\n"
			output += "panic(err)\n"
			output += "}\n"
			for k, act := range handler.Actions {
				if act.Id != "" {
					output += fmt.Sprintf("action%d_%d_%d, err := handler%d_%d.NewAction(\"%s\")\n", i, j, k, i, j, act.Id)
				} else {
					port := app.imports.AddWithAlias("", act.Ref)
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
					output += fmt.Sprintf("action%d_%d_%d, err := handler%d_%d.NewAction(&%s.Action{}, %s)\n", i, j, k, i, j, port.Alias, settingsName)
				}
				output += "if err != nil {\n"
				output += "panic(err)\n"
				output += "}\n"
				if act.If != "" {
					output += fmt.Sprintf("action%d_%d_%d.SetCondition(%s)\n", i, j, k, act.If)
				}
				if length := len(act.Input); length > 0 {
					mappings := make([]string, length)
					for key, value := range act.Input {
						mappings = append(mappings, fmt.Sprintf("%s=%v", key, value))
					}
					output += fmt.Sprintf("action%d_%d_%d.SetInputMappings(%#v)\n", i, j, k, mappings)
				}
				if length := len(act.Output); length > 0 {
					mappings := make([]string, length)
					for key, value := range act.Output {
						mappings = append(mappings, fmt.Sprintf("%s=%v", key, value))
					}
					output += fmt.Sprintf("action%d_%d_%d.SetOutputMappings(%#v)\n", i, j, k, mappings)
				}
				output += "if err != nil {\n"
				output += "panic(err)\n"
				output += "}\n"
				output += fmt.Sprintf("_ = action%d_%d_%d\n", i, j, k)
			}
			output += fmt.Sprintf("_ = handler%d_%d\n", i, j)
		}
		output += fmt.Sprintf("_ = trg%d\n", i)
	}
	output += "e, err := api.NewEngine(app)\n"
	output += "if err != nil {\n"
	output += "panic(err)\n"
	output += "}\n"
	app.imports.Add("github.com/project-flogo/core/engine")
	output += "engine.RunEngine(e)\n"
	output += "}\n"

	header := "package main\n\n"
	header += "import (\n"
	for _, port := range app.imports.Imports {
		if port.Alias != "" {
			header += fmt.Sprintf("%s \"%s\"\n", port.Alias, port.Import)
			continue
		}
		header += fmt.Sprintf("\"%s\"\n", port.Import)
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
