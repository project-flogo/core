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

// Generate generates flogo go API code
func Generate(config *app.Config, file string) {
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

	output := "func main() {\n"
	app.imports.Add("github.com/project-flogo/core/api")
	output += "app := api.NewApp()\n"
	for i, trigger := range config.Triggers {
		port := app.imports.AddWithAlias("", trigger.Ref)
		output += fmt.Sprintf("trg%d := app.NewTrigger(&%s.Trigger{}, %#v)\n", i, port.Alias, trigger.Settings)
		for j, handler := range trigger.Handlers {
			output += fmt.Sprintf("handler%d_%d, err := trg%d.NewHandler(%#v)\n", i, j, i, handler.Settings)
			output += "if err != nil {\n"
			output += "panic(err)\n"
			output += "}\n"
			for k, act := range handler.Actions {
				actionConfig := act.Config
				if actionConfig.Id != "" {
					for _, act := range config.Actions {
						if actionConfig.Id == act.Id {
							actionConfig = act
							break
						}
					}
				}
				port := app.imports.AddWithAlias("", actionConfig.Ref)
				factory, settingsName := action.GetFactory(actionConfig.Ref), fmt.Sprintf("settings%d_%d_%d", i, j, k)
				if generator, ok := factory.(Generator); ok {
					code, err := generator.Generate(settingsName, &app.imports, actionConfig)
					if err != nil {
						panic(err)
					}
					output += "\n"
					output += code
					output += "\n"
				} else {
					output += fmt.Sprintf("var %s map[string]interface{}\n", settingsName)
				}
				output += fmt.Sprintf("action%d_%d_%d, err := handler%d_%d.NewAction(&%s.Action{}, %s)\n", i, j, k, i, j, port.Alias, settingsName)
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
