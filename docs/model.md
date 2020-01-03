# Application Model

Sections:

* [Imports](#imports "Goto Imports") - Go package and contribution imports
* [Properties](#properties "Goto Properties") - Shared Properties
* [Channels](#channels "Goto Channels") - Internal Communication Channels
* [Triggers](#imports "Goto Triggers") - Triggers
* [Actions](#actions "Goto Actions") - Shared Actions
* [Resources](#resources "Goto Resources") - Shared Resources
* [Schemas](#schemas "Goto Schemas") - Shared Schemas
* [Connections](#connections "Goto Connections") - Shared Connections
    
[Full Example](#full-example "Full Example") 

## Imports
The imports section allows one to define all the contributions that should be imported by the engine. In some instances
 this may include go code references, for example a specific database driver if your application logic is dependent on 
 it.

```json
  "imports": [
    "github.com/project-flogo/flow",
    "github.com/project-flogo/contrib/trigger/rest",
    "github.com/project-flogo/contrib/activity/log"
  ],
```

A contribution can be referenced directly using the full go package or indirectly to its imported package.

Direct: `"ref" : "github.com/project-flogo/flow"`
 
Indirect: `"ref" : "#flow"` 

## Properties
The properties section allows one to define properties that can be shared by the application.

```json
  "properties": [
    {"name":"myProp", "type":"string", "value":"myValue" }
  ]
```
The properties section allows one to define properties that can be shared by the application.

A property can be resolved in an expression using the property resolver  `$property[propertyName]`

Example:

```json
  "settings": {
    "mySetting": "=$property[myProp]"
  } 
```

## Channels
The channels section allows one to define internal communications channels for the engine.

```json
  "channels": [
    "myChannel:5"
  ]
```
A channel is used for internal communications in an engine.  It is defined by a channel name and buffer size. 

The `github.com/project-flogo/contrib/activity/channel` activity can be used to put a message on the channel.  The `github.com/project-flogo/contrib/trigger/channel` trigger can be used to listen on a channel and trigger actions from messages received on that channel.

## Triggers
The triggers section is used to define the triggers that will be used by the application.
 
```json
 "triggers": [ 
    { 
      "id": "my_rest_trigger",
      "ref": "#rest",
      "settings": { 
        "port": "9233"
      },
      "handlers": [ 
        { 
          "settings": { 
            "method": "GET",
            "path": "/test"
          },
          "action": { 
            "id": "sharedAction"
          },
        }
      ]
    }
 ]
```
Handlers can define the actions they use in two ways.
One can refer to a shared action via 'id' like above, where the handler refers to the shared action "sharedAction"

One can also define the action inline:

```json
"handlers": [ 
  { 
    "action": {
      "ref": "#flow",
      "settings": {
        "flowURI": "res://flow:myflow"
      }
    }
  }
]
```
There is also a special case where multiple actions can be defined in a single handler.  The handler will have
an `actions` section instead of `action`.  The actions can have an "if" property to determine which action will 
be invoked.  The first action whose 'if' expression evaluates to true will be executed.  An action that doesn't 
have an 'if' property it will be considered the same as `"if":"true"`

Example:

```json
"actions": 
  [
    {
      "if": "$.headers.Foo == \"bar\"",
      "id": "sharedAction"
    },
    {
      "id": "sharedActionDefault"
    }
  ]
```
In this example the action "sharedAction" is executed if the header Foo = "bar", otherwise the "sharedActionDefault"
is executed.

## Actions
The actions section is used to define shared actions that can be referenced by id.

```json
  "actions": [
    {
      "id" : "sharedAction",
      "ref": "#flow",
      "settings": {
        "flowURI": "res://flow:myflow"
      }
    }
  ]
```

## Resources
The resources section contains the resources used by actions.

```json
  "resources": [
    {
      "id": "flow:myflow",
      "data": {
        "name": "My Flow",
        "description": "Example description",
        "tasks": [
          ...
        ]
      }
    }
  ]  
```

## Schemas
The schemas section contains schemas that are shared in the application.

```json
  "schemas": {
    "mySchema": { 
      "type": "json", 
      "value": "{\"$schema\": \"http://json-schema.org/draft-07/schema#\", ... }"
    }
  }  
```

Schemas can be referenced in metadata or in schema sections in contributions.

```json
  "activity": {
   "ref": "#myActivity",
    "input": {
      "val": "=$flow.Val"
    },
   "schemas": {
       "input": {
         "val": "schema://mySchema"
       }
   }
}
```
Schemas can also be defined inline:

```json
  "activity": {
   "ref": "#myActivity",
    "input": {
      "val": "=$flow.Val"
    },
   "schemas": {
       "input": {
         "val": { "type": "json","value": "{\"$schema\": \"http://json-schema.org/draft-07/schema#\", ... }" }
       }
   }
}
```

## Connections
The connections section contains connections that can be shared by triggers, actions or
activities.  It allows you to define a connection once and use it in multiple places.

```json
"connections": {
  "myConn": {
    "ref" : "github.com/project-flogo/contrib/connection/sql",
    "settings" : {
      "dbType": "mysql",
      "driver": "mysql",
      "dataSource": "username:password@tcp(host:port)/dbName",
    }
  }
}
```

These connections can be reference by values of type `connection`, 

```json
"activity" : {
  "ref": "#sql",
  "settings": {
    "connection": "conn://myConn"
  }
}
```

## Full Example
Sample flogo application configuration file. 

```json
{
  "name": "simpleApp",
  "type": "flogo:app",
  "version": "0.0.1",
  "appModel": "1.0.0",
  "description": "My flogo application description",
  "imports": [
    "github.com/project-flogo/flow",
    "github.com/project-flogo/contrib/trigger/rest",
    "github.com/project-flogo/contrib/activity/log"
  ],
  "triggers": [
    {
      "id": "my_rest_trigger",
      "ref": "#rest",
      "settings": {
        "port": "9233"
      },
      "handlers": [
        {
          "settings": {
            "method": "GET",
            "path": "/test"
          },
          "action": {
            "ref": "#flow",
            "settings": {
              "flowURI": "res://flow:myflow"
            },
            "input": {
              "orderType": "standard"
            },
            "output": {
              "data": "=$.value"
            }
          }
        },
        {
          "settings": {
            "method": "GET",
            "path": "/test2"
          },
          "actions": [
            {
              "if": "$.headers.Foo == \"bar\"",
              "id": "sharedAction",
              "input": {
                "orderType": "bar"
              },
              "output": {
                "data": "=$.value"
              }
            },
            {
              "id": "sharedActionDefault",
              "input": {
                "orderType": "foo"
              },
              "output": {
                "data": "fixed"
              }
            }
          ]
        }
      ]
    }
  ],
  "actions":[
    {
      "id" : "sharedAction",
      "ref": "#flow",
      "settings": {
        "flowURI": "res://flow:myflow"
      }
    },
    {
      "id" : "sharedActionDefault",
      "ref": "#flow",
      "settings": {
        "flowURI": "res://flow:myflow"
      }
    }
  ],
  "resources": [
    {
      "id": "flow:myflow",
      "data": {
        "name": "My Flow",
        "description": "Example description",
        "metadata": {
          "input": [
            { "name":"customerId", "type":"string" },
            { "name":"orderId", "type":"string" },
            { "name":"orderType", "type":"string" }
          ],
          "output":[
            { "name":"value", "type":"string" }
          ]
        },
        "tasks": [
          {
            "id": "FirstLog",
            "name": "FirstLog",
            "type": "iterator",
            "settings": {
              "iterate": 10
            },
            "activity": {
              "ref": "#log",
              "input": {
                "message": "=$iteration[key]"
              }
            }
          },
          {
            "id": "SecondLog",
            "name": "SecondLog",
            "activity" : {
              "ref": "#log",
              "input": {
                "message": "test message"
              }
            }
          }
        ],
        "links": [
          {
            "from": "FirstLog",
            "to": "SecondLog",
            "type": "expression",
            "value": "$flow.orderId > 1000"
          }
        ],
        "errorHandler": {
          "tasks": [
            {
              "id": "ErrorLog",
              "activity": {
                "ref": "#log",
                "input": {
                  "message": "log in error handler"
                }
              }
            }
          ]
        }
      }
    }
  ]
}
```
