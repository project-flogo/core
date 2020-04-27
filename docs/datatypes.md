# Data Types

#### Standard Data Types:
|Name|Go Type |Description|
|--- |--- |--- |
| any     | interface{} |  Can be any value |
| string  | string      | A string|
| int     | int   | Integer, the size is system dependent |
| int32   | int32 | 32 bit integer |
| int64   | int64 | 64 bit integer |
| float32 | float32 | 32 bit float |
| float64 | float64 | 64 bit float |
| bool    | bool    | Boolean |
| bytes   | []byte  | Byte array |
| array   | []interface{} | Array of "any value" |
| object  | map[string]interface{} | Object - typically translated JSON object |
| params  | map[string]string      | Parameter map with string keys and values |
| map     | map[interface{}]interface{} | A map with any type of key or value |

#### Special Data Types:
|Name|Description|
|--- |--- |
| connection | Special type to denote a connection |

## Connection
The **connection** datatype is used to indicate that the value is "connection".  A connection is
typically used to connect to a networked resource.  These connections can be shared or defined
in-line.  When a connection is defined in-line, it is only used by that contribution.  Shared 
connection is defined once in the flogo.json and can shared between multiple resources.  

#### In-line

```json
"activity" : {
  "ref": "#sql",
  "settings": {
    "connection": {
      "ref" : "github.com/project-flogo/contrib/connection/sql",
      "settings" : {
        "dbType": "mysql",
        "driver": "mysql",
        "dataSource": "username:password@tcp(host:port)/dbName",
      }
    }
  }
}
```

#### Shared

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

```json
"activity" : {
  "id": "a1",
  "ref": "#sql",
  "settings": {
    "connection": "conn://myConn"
  }
},
"activity" : {
  "id": "a2",
  "ref": "#sql",
  "settings": {
    "connection": "conn://myConn"
  }
}
```

## Coercion

Values are automatically coerced to the expected type.  Most conversion are straight forward.  In cases of 
complex types (params, map, object, array), json representations are also converted. 

