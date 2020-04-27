# Go API to Run Flogo Application

## Go API's
#### func  EvalActivity

```go
func EvalActivity(act activity.Activity, input interface{}) (map[string]interface{}, error)
```
EvalActivity evaluates the specified activity using the provided inputs

#### func  NewActivity

```go
func NewActivity(act activity.Activity, settings ...interface{}) (activity.Activity, error)
```
NewActivity creates an instance of the specified activity

#### func  NewEngine

```go
func NewEngine(a *App) (engine.Engine, error)
```
NewEngine creates a new flogo Engine from the specified App

#### func RunAction

```go
func RunAction(ctx context.Context, act action.Action, inputs map[string]interface{}) (results map[string]interface{}, err error)
```
RunAction runs a specified action and returns the output 
#### type Action

```go
type Action struct {
}
```

Action is the structure that defines the Action for a Handler

#### func (*Action) Condition

```go
func (a *Action) Condition() string
```
Condition returns the condition

#### func (*Action) InputMappings

```go
func (a *Action) InputMappings() []string
```
InputMappings gets the Action's input mappings

#### func (*Action) OutputMappings

```go
func (a *Action) OutputMappings() []string
```
OutputMappings gets the Action's output mappings

#### func (*Action) SetCondition

```go
func (a *Action) SetCondition(condition string)
```
SetCondition sets the conditional expression which determines if the action is
executed

#### func (*Action) SetInputMappings

```go
func (a *Action) SetInputMappings(mappings ...string)
```
SetInputMappings sets the input mappings for the Action, which maps the outputs
of the Trigger to the inputs of the Action

#### func (*Action) SetOutputMappings

```go
func (a *Action) SetOutputMappings(mappings ...string)
```
SetOutputMappings sets the output mappings for the Action, which maps the
outputs of the Action to the return of the Trigger

#### func (*Action) Settings

```go
func (a *Action) Settings() map[string]interface{}
```
Settings gets the settings of the Action

#### type App

```go
type App struct {
}
```

App is the structure that defines an application

#### func  NewApp

```go
func NewApp() *App
```
NewApp creates a new Flogo application

#### func (*App) Actions

```go
func (a *App) Actions() map[string]*Action
```
Triggers gets the Triggers of the application

#### func (*App) AddAction

```go
func (a *App) AddAction(id string, act action.Action, settings interface{}) error
```
AddAction adds an action to the application

#### func (*App) AddProperty

```go
func (a *App) AddProperty(name string, dataType data.Type, value interface{}) error
```
AddProperty adds a shared property to the application

#### func (*App) AddResource

```go
func (a *App) AddResource(id string, data json.RawMessage)
```
AddResource adds a Flogo resource to the application

#### func (*App) NewIndependentAction

```go
func (a *App) NewIndependentAction(act action.Action, settings interface{}) (action.Action, error)
```

#### func (*App) NewTrigger

```go
func (a *App) NewTrigger(trg trigger.Trigger, settings interface{}) *Trigger
```
NewTrigger adds a new trigger to the application

#### func (*App) Properties

```go
func (a *App) Properties() map[string]data.TypedValue
```
Properties gets the shared properties of the application

#### func (*App) Triggers

```go
func (a *App) Triggers() []*Trigger
```
Triggers gets the Triggers of the application

#### type Handler

```go
type Handler struct {
}
```

Handler is the structure that defines the handler for a Trigger

#### func (*Handler) Actions

```go
func (h *Handler) Actions() []*Action
```
Actions gets the Actions of the Handler

#### func (*Handler) NewAction

```go
func (h *Handler) NewAction(handlerAction interface{}, settings ...interface{}) (act *Action, err error)
```
NewAction adds a new Action to the Handler

#### func (*Handler) Settings

```go
func (h *Handler) Settings() map[string]interface{}
```
Settings gets the Handler's settings

#### type HandlerFunc

```go
type HandlerFunc func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)
```

HandlerFunc is the signature for a function to use as a handler for a Trigger

#### type Trigger

```go
type Trigger struct {
}
```

Trigger is the structure that defines a Trigger for the application

#### func (*Trigger) Handlers

```go
func (t *Trigger) Handlers() []*Handler
```
Handlers gets the Trigger's Handlers

#### func (*Trigger) NewHandler

```go
func (t *Trigger) NewHandler(settings interface{}) (*Handler, error)
```
NewHandler adds a new Handler to the Trigger

#### func (*Trigger) Settings

```go
func (t *Trigger) Settings() map[string]interface{}
```
Settings gets the Trigger's settings


## Run Flow
This short tutorial walks through how to run Flogo Actions using the Go-API provided by the flogo core.

First Import the action you need to run. In this example import Flow Action. Initialize the Flogo App. Initialiaze and configure a Rest Trigger. 
Add the Flow Action to the desired Trigger Handler.

```go
package main

import (
	"fmt"

	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/flow"

	_ "github.com/project-flogo/contrib/activity/log" //our flow contains a log activity, so we need to include this
	"github.com/project-flogo/contrib/trigger/rest"
)

func main() {
	// Initialize the Flogo App
	app, err := myApp()
    if err != nil {
    		fmt.Println("Error:", err)
    		return
    }
    // Create a Flogo Engine using the Flogo App.
	e, err := api.NewEngine(app)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
    // Run the engine.
	engine.RunEngine(e)
}

func myApp() (*api.App, error) {
	// Initialize a Flogo App.
    app := api.NewApp()
    
    // Intialize a New Rest Trigger.
	trg := app.NewTrigger(&rest.Trigger{}, &rest.Settings{Port: 8080})
    
    // Create a New Handler using the Trigger Handler
	h, err := trg.NewHandler(&rest.HandlerSettings{Method: "GET", Path: "/blah/:val"})
    if err != nil {
        return nil, err
    }
	
    // Set Flow Action Settings. The flow is defined in local file system
    settings :=  &flow.Settings{FlowURI:"file://myflow.json"}
    
    // Initialize a Flow Action using settings.
    act, err := h.NewAction(&flow.FlowAction{}, settings)
    if err != nil {
    	return nil, err
    }
    
    // map the path parameter from the handler to input of the flow
    act.SetInputMappings("in=$.pathParams.val")

	return app, nil
}
```


## Run Activities.  

Instead of having a predefined flow in a JSON file, you might want to run selected activities. To achieve this, a new activity is initialized and registered to the Trigger Handler. 

Complete Example

```go
import (
...
"github.com/project-flogo/core/activity"
"github.com/project-flogo/core/api"
...
)

func main() {

	app := myApp()

	e, err := api.NewEngine(app)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	engine.RunEngine(e)
}

func myApp() *api.App {
	app := api.NewApp()

	trg := app.NewTrigger(&rest.Trigger{}, &rest.Settings{Port: 8080})
	h, err := trg.NewHandler(&rest.HandlerSettings{Method: "GET", Path: "/blah/:num"})
	if err != nil {
        //Handle Error.
    }
    // Register Trigger handler.
	h.NewAction(RunActivities)
    
    // Initialize the Log Activity
	logAct, err := api.NewActivity(&log.Activity{})
    if err != nil {
    // Handler Error
    }
    //store in map to avoid activity instance recreation
	activities = map[string]activity.Activity{"log": logAct}

	return app
}

var activities map[string]activity.Activity

// Trigger Handler
func RunActivities(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
	
    trgOut := &rest.Output{}
    
    // Get the inputs from the trigger. 
	trgOut.FromMap(inputs)
    
	msg, _ := coerce.ToString(trgOut.PathParams)
    
    // Run the Log Activity. The output of the activity is returned in map.
	out, err := api.EvalActivity(activities["log"], &log.Input{Message: msg})
	if err != nil {
		return nil, err
	}

	reply := &rest.Reply{Code: 200, Data: out["someOutput"]}

	return reply.ToMap(), nil
}
```

## Run Independent Action.

Just like running activities, you might want to directly run a certain action. To achieve this a New Independent Action is initialized and Registered to Trigger Handler

```go
import (
...
"github.com/project-flogo/core/action"
"github.com/project-flogo/core/api"
cml"github.com/project-flogo/catalystml-flogo/action"
...
)
var actions map[string]action.Action


func myApp() *api.App {
    ...
    // A Map containing Activities
    actions := make(map[string]action.Action)
    
    //Initialize Action
    cmlAct, err := api.NewIndependentAction(&cml.Action{}, map[string]interface{}{"catalystMLURI":"file://URI"})
    
    //Register Activites in Map
    action["cmlAction"] = cmlAct
    ...	
}


// Trigger Handler Signature func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)
// Tigger Handler
func RunActivities(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
	
    // Run CML Action.
    out, err := api.RunAction(ctx, action["cmlAction"], inputs)
    if err != nil {
    	return nil, err
    }
    
    //Converting Activity Output to Trigger Output Type
    reply := &rest.Reply{Code: 200, Data: out["someOutput"]}
    return reply.ToMap(), nil
}
``` 
