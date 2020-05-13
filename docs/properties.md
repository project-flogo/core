#Properties
In Flogo, the concept of an application-level property bag is available to allow reuse of properties across different actions, triggers and activities.  For example, they can be used in trigger or activity settings. Properties are accessed via the `$property` resolver and made available to the scopes defined in the [mappings](mapping.md) documentation.

### Configuration

Properties are defined within the top-level of the flogo application json, as shown below via the `properties` element.

```json
{
  "name": "default_app",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "Sample flogo app",
  "properties": [
     {
       "name": "my_property",
       "type": "string",
       "value": "My Property Value"
     }
  ]
```

As previously stated, properties are accessible via the `$property` resolver. Consider the following mappings into a log activity:

```json
{
  "id": "log_2",
  "name": "Log",
  "description": "Logs a message",
  "activity": {
    "ref": "#log",
    "input": {
      "message": "=$property[my_property]"
    }
  }
}
```

### Grouping of properties
Even though the engine itself doesn't support property grouping, this can be accomplished by using a naming convention in your application. Since property names allow the use of `.`, a naming convention like `<group>.<sub-group>...<name>` can be used to create an artifical grouping of related properties. 

```json
{
  "name": "default_app",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "Sample flogo app",
  "properties": [
     {
       "name": "PURCHASE.SERVICE.DB.URL",
       "type": "string",
       "value": "postgres://10.10.10.10:5370/mydb"
     },
     {
       "name": "PURCHASE.SERVICE.DB.USER",
       "type": "string",
       "value": "testuser"
     },
     {
        "name": "INVENTORY.SERVICE.DB.URL",
        "type": "string",
        "value": "postgres://10.10.10.20:5370/mydb"
     },
     {
        "name": "INVENTORY.SERVICE.DB.USER",
        "type": "string",
        "value": "testuser"
     }
  ]
```


These properties can be accessed via `$property[PURCHASE.SERVICE.DB.URL]` or `$property[INVENTORY.SERVICE.DB.URL]`

### Overriding properties at runtime

In order to override properties at runtime, you have to enable external property resolvers.

This can be done by setting the `FLOGO_APP_PROP_RESOLVERS` environment variable.  Currently, there are two built-in external
property resolvers: json(JSON) and env(Environment Variable).


```terminal
FLOGO_APP_PROP_RESOLVERS=env,json ./<app_binary>
```

You can override app properties at runtime in two ways:

#### Resolver: json

When using the `json` property resolver, you can provide a comma separated list of json files that
will override the application's existing property values.
```env
FLOGO_APP_PROPS_JSON=app1.json,common.json
```

**Example**

Let's say you want to override some of your properties.  You will need to define your new value for a given property in your json file.

_props.json_

```json
{
 "MyProp1": "This is new value",
 "MyProp2": 20
}
```

Now run your application:

```terminal

export FLOGO_APP_PROPS_JSON=props.json 
FLOGO_APP_PROP_RESOLVERS=json ./MyApp
```

#### Resolver: env

In order to override properties using environment variables, you just need enable the `env` property resolver

```terminal
FLOGO_APP_PROP_RESOLVERS=env ./<app_binary>
```

**Example**

Let's say you want to override `myprop` property in your app.  You would do the following:

```terminal

export myprop=bar
FLOGO_APP_PROP_RESOLVERS=env ./MyApp
```


### Custom External Resolver

You can plug-in your own property resolver to resolve application property values from external configuration management services, such as, Consul, Spring Cloud Config etc. Just implement the following interface and register implementation with the runtime:

```go
// Resolver used to resolve property value from external configuration like env, file etc
type ExternalResolver interface {
	// Name of the resolver (e.g., consul)
	Name() string
	// Should return value and true if the given key exists in the external configuration otherwise should return nil and false.
	LookupValue(key string) (interface{}, bool)
}

```

#### Sample Resolver

```go
package sampleresolver

import "github.com/project-flogo/core/data/property"


type SamplePropertyResolver struct {
}

func init() {
  _ = property.RegisterExternalResolver(&SamplePropertyResolver{})
}

func (resolver *SamplePropertyResolver) Name() string {
   return "sampleresolver"
}

func (resolver *SamplePropertyResolver) LookupValue(propertyName string) (interface{}, bool) {
   // Resolve property value
  return "some_value"", true
}
```
*Note: In order for your resolver to be loaded in the go code, you need to add an entry to your resolver in the imports section of the engine.json*


Set the `FLOGO_APP_PROP_RESOLVERS` environment variable to `sampleresolver` while running your application. For example:

```terminal
FLOGO_APP_PROPS_RESOLVERS=sampleresolver ./<app_binary>
```