# Data

## Types

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

## Coercion

Values are automatically coerced to the expected type.  Most conversion are straight forward.  In cases of 
complex types (params, map, object, array), json representations are also converted. 

