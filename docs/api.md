# Go API to Run Flogo Application

## Run Flow
This short tutorial walks through how to run Flogo Actions using the Go-API provided by the flogo core.

First Import the action you need to run .

```go 
import (
    ...
    "github.com/project-flogo/core/api"
    "github.com/project-flogo/flow"
    ...
)
```

Then initalize the Flogo app to use Go-API

```go
app := api.NewApp()
```

Setup the trigger and handler that invokes the action.   In this case REST trigger is selected to run on port 8080. GET "/blah/:value" (on REST trigger) is registered as a handler.
```go

trg := app.NewTrigger(&rest.Trigger{}, &rest.Settings{Port: 8080})
h, _ := trg.NewHandler(&rest.HandlerSettings{Method: "GET", Path: "/blah/:num"})

```

Now setup Flow Action that runs the flow from local file myflow.json. Then associate that action with our handler.
```go

settings :=  &flow.Settings{FlowURI:"file://myflow.json"}
a, _ h.NewAction(&flow.FlowAction{}, settings)

```

Now lets map the path parameter from the handler to input of our flow.
```go

a.SetInputMappings("in=$.pathParams.val")

```

### Full Example
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
	h, _ := trg.NewHandler(&rest.HandlerSettings{Method: "GET", Path: "/blah/:val"})

	settings :=  &flow.Settings{FlowURI:"file://myflow.json"}
	a, _ := h.NewAction(&flow.FlowAction{}, settings)
	a.SetInputMappings("in=$.pathParams.val")

	return app
}
```
## Run Activities.  

Instead of having a predefined flow in a JSON file, you might want to run selected activities. To achieve this, a new activity is initialized and registered to the Trigger Handler. 

```go
import (
...
"github.com/project-flogo/core/activity"
"github.com/project-flogo/core/api"
...
)

// A Map containing Activities
activities := make(map[string]activity.Activity)

//Initialize Activity
logAct, err := api.NewActivity(&log.Activity{})

//Register Activites in Map
activities["logAct"] = logAct
```

The Trigger handler has an activity or collection of activities which is runned by calling "EvalActivity".

```go
// Trigger Handler Signature func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)
// Tigger Handler
func RunActivities(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
    
    // Setting the inputs from trigger to type of
    // Rest Tirgger Output
    trgOut := &rest.Output{}
    trgOut.FromMap(inputs)
    
    // Get Message from Inputs
    msg, _ := coerce.ToString(trgOut.PathParams)
    
    // Run Log Activity.
    out, err := api.EvalActivity(activities["log"], &log.Input{Message: msg})
    if err != nil {
    	return nil, err
    }
    
    //Converting Activity Output to Trigger Output Type
    reply := &rest.Reply{Code: 200, Data: out["someOutput"]}
    return reply.ToMap(), nil
}
```  

Registering the Trigger Handler to Trigger.

```go
    h, err := trg.NewHandler(&rest.HandlerSettings{Method: "GET", Path: "/blah/:num"})
	h.NewAction(RunActivities)
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

// A Map containing Activities
actions := make(map[string]action.Action)

//Initialize Activity
cmlAct, err := api.NewIndependentAction(&cml.Action{}, map[string]interface{}{"catalystMLURI":"file://URI"})

//Register Activites in Map
action["cmlAction"] = cmlAct
```

The Trigger handler contains the Action to Run.

```go
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

Registering the Trigger Handler to Trigger.

```go
    h, err := trg.NewHandler(&rest.HandlerSettings{Method: "GET", Path: "/blah/:num"})
	h.NewAction(RunActivities)
```