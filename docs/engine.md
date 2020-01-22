# Engine


## Configuration DLS

Engine specific configuration can be set in an *engine.json* file.  This file should be placed 
along side the flogo.json

Sections:

* [Imports](#imports "Goto Imports") - Go package imports
* [ActionSettings](#actionsettings "Goto Action Settings") - Action Runtime Settings
* [Services](#services "Goto Services") - Engine Service Configurations
    
[Full Example](#full-example "Full Example") 

## Imports
The `imports` section allows one to define all the non-contribution packages that should be imported by the engine. 
These are references to go code implementations that aren't specific to the application logic.  They tend to be 
runtime and environment specific configuration.  For example here is where you would add a reference to a custom
property resolver.  

```json
  "imports": [
    "github.com/myuser/mypropertyresolver"
  ],
```

## ActionSettings
The `actionSettings` section allows one to specify runtime specific settings of an action that aren't related to 
the application logic.  For example, you might want to enable state recording for a flow action.

```json
  "actionSettings": {
    "github.com/project-flogo/flow": {
      "stepRecordingMode": "full",
    }
  },
```

Look at an action's documentation to see if they expose runtime specific properties.

## Services
The services section allows one to define engine services.  To configure the service, the
ref to the service implementation has to be specified.  Any necessary settings for that service
can be specified via the `settings` section.  You can enabled or disabled via the `enabled`
property.


#### Example Remote State Recorder
```json
    {
      "name": "flowStateRecorder",
      "ref": "github.com/project-flogo/services/flow-state/client/rest",
      "enabled": true,
      "settings": {
        "host": "192.168.1.50",
        "port": "9190"
      }
    }
```
#### Example Local State Recorder
```json
    {
      "name": "flowStateRecorder",
      "ref": "github.com/project-flogo/services/flow-state/client/local",
      "enabled": true
    }
```

## Full Example
Sample engine runtime configuration file. 

```json
{
  "type": "flogo:engine",
  "imports": [
    "github.com/project-flogo/services/flow-state/client/rest@master",
    "github.com/project-flogo/stream/service/telemetry@master"
  ],
  "actionSettings": {
    "github.com/project-flogo/flow": {
      "stepRecordingMode": "full"
    }
  },
  "services": [
    {
      "name": "telemetry",
      "ref": "github.com/project-flogo/stream/service/telemetry",
      "enabled": true
    },
    {
      "name": "flowStateRecorder",
      "ref": "github.com/project-flogo/services/flow-state/client/rest",
      "enabled": true,
      "settings": {
        "host": "192.168.1.50",
        "port": "9190"
      }
    }
  ]
}
```
