# Mapping

## What are mappings?

A mapping in Flogo is used to assign the value of one parameter (a trigger output, for example) to that an input 
parameter (a flow's input, for example).  This expression can also be used to manipulate that value during the 
assignment.

#### Expression
The simplest type of mapping is an expression, which is denoted by a string that starts with `=`.

These mappings are quite straightforward, for example:
```json
{
  "myInput":"=$.pathParams.myParam"
}
```

The above mapping indicates that the value of `pathParams.myParam` from a trigger input should be mapped to the action input 
named `myInput`. 

#### Object
Another type of mapping is an object mapping.  These are denoted by having the value of the parameter be an object 
named `mapping`.  The `mapping` object is used to define how the object should be constructed and how the various fields 
within the object are mapped. Mapping expressions can be used to assign values to the fields of the object.

Example:
```json
{
  "bookDetails": {
    "mapping": {
      "Author": "=$flow.author",
      "ISBN": "=$flow.isbn",
      "Price": "$20"
    }
  }
}
```

In the `mapping` object you can also iterate through values of a parameter to construct the object.  The `@foreach` directive is used to iterate
over an array of values.


Example: Iterate over the array `$fow.store.books` to construct `books`
```json
{
  "books": {
    "mapping": {
      "@foreach($flow.store.books)": {
        "author": "=$loop.author",
        "title": "=$loop.title",
        "price": "=$loop.price"
      }
    }
  }
}
```

The above example shows how `foreach` works. It will iterate over `$flow.store.books` and extract the author, title and
 price to target for each entry and assign it to the array `books`, The final target books:
```json
{
  "books": [
   {
    "author":"an author",
    "title":"title 1",
    "price": 33.33
    },
   {
    "author":"another author",
    "title":"title 2",
    "price": 22.22
    }
  ] 
}
```

## Mapping Resolvers

Mapping resolvers are used in mapping expression to lookup a value.

The following table is a list of the standard resolvers. These are available in most mapping expressions.

|Resolver|Description|
|--- |--- |
|$env|Used to resolve an environment variable|
|$property|Used to resolve properties from the global application property bag|
|$loop|Used to resolve an the current value in the foreach loop of an object mapping|
|$.|Used to resolve in the current scope|


Individual actions can also have their own set of resolvers. For example, in addition to the standard resolvers 
the flow has additional resolvers that are available in its mapping expressions.  The following table is a list of resolvers
used by the `flow` action.  

|Resolver|Description|
|--- |--- |
|$flow|Used to resolve params from within the current flow. If a flow has a single trigger and no input params defined, then the output of the trigger is made available via $flow|
|$activity|Used to resolve activity params. Activities are referenced by id, for example, $activity[activity_id].activity_property.|
|$iteration[key] |Used to resolve data scoped to an iterator, key is the key of the current iteration|
|$iteration[valye] |Used to resolve data scoped to an iterator, value is the value of the current iteration|

#### Scopes

Resolvers are used to look up data in different scopes.  For example, the `.` indicates that the value is available within the current scope. 

```json
{
  "isbn": "=$.event.isbn"
}
```
    
The above mapping is from the Trigger/Handler, which we know, based on the indication of the `.`, we can only access trigger scoped (output) variables, 
thus `event.isbn` is within the trigger scope, as indicated by the preceding `.`.


What if you’re accessing a value outside of the immediate scope? The mapping expression should use the corresponding resolver to access the value. For example, 
consider the following.
```json
{
  "flowName": "=$flow.name"
}
```
    
This mapping is associated with Flow Action. Let's say we have an activity that takes the flow name.  We know that the value isn't in our immediate scope, 
hence the `$flow` resolver should be used. In the above snippet, we’re grabbing the value of the flow variable named `name`, hence `$flow.name` is used. If 
we wanted to grab the value of an environment variable we could use `$env.VarName`.


## Functions

Mapping expressions also support functions.  These functions are used to manipulate the data that is being assigned.
You may want to add some custom logic to the mapping, such as concat/substring/length of a string or generate a random 
number base on a range and so on. Refer to the [functions repository](https://github.com/project-flogo/contrib/tree/master/function) 
for available functions. Also note, you can install functions using the CLI’s `flogo install` command.

Example:
```json
{
  "description": "=string.concat(\"The pet category name is: \", $flow.pet.category.name)"
}
```   

## Additional Details

### Accessing object properties

Most of the time you wont want to perform a direct assigning from one complex object to another, rather you’ll want to grab a simple type property from one complex object and perform a direct 
assigning to another property. This can be done accessing children using a simple dot notation. For example, consider the following mapping.

    {
      "someObject": {
        "mapping": {
          "Title": "=$activity[rest_3].result.items[0].volumeInfo.title",
          "PublishedDate": "=$activity[rest_3].result.items[0].volumeInfo.publishedDate",
          "Description": "=$activity[rest_3].result.items[0].volumeInfo.description"
        }
      }
    }

someObject is of type `object` and has the properties `Titie`, `PublishedDate`, `Description` which are being mapped from the response of an activity, this is fetched using the `$activity` 
resolver. Consider one of the examples:

`$activity[rest_3].result.items[0].volumeInfo.title`

We’re referencing the result property from the activity named rest_3. We’re then accessing an items array (the first entry of the array) to another complex object, where finally 
we’re at a simple string property named title.


### Handling arrays in mappings

There are lots of use cases for array mapping, map entire array to another or iterator partial array to another with functions The array mapping value comes from a JSON format

Case 1: iterate on array `$flow.store.books` and assign value to `books`

    {
      "books": {
        "mapping": {
          "@foreach($flow.store.books)": {
            "author": "=$loop.author",
            "title": "=$loop.title",
            "price": "=$loop.price"
          }
        }
      }
    }

Case 2: Copy original array `$fow.store.books` to target array `books`

    {
      "books": {
        "mapping": {
          "@foreach($flow.store.books)": {
            "=": "$loop"
          }
        }
      }
    }

Case 3: Iterate on array `$fow.store.books` and assign to primitive array `titles`

    {
      "titles": {
        "mapping": {
          "@foreach($flow.store.books)": {
            "=": "$loop.title"
          }
        }
      }
    }

Case 4: Accessing parent loop data.

    {
      "books": {
        "mapping": {
          "@foreach($flow.store.books, bookLoop)": {
            "title": "=$loop.title",
            "price": "=$loop.price",
            "author": {
             "@foreach($loop.author, authorLoop)": {
                "firstName": "=$loop.firstName",
                "lastName": "=$loop[authorLoop].lastName",
                "bookTitle": "=$loop[bookLoop].title"
             }
            }
          }
        }
      }
    }

Case 5: Using fixed array

    {
      "store": {
        "mapping": {
          "store": {
            "books": [
              {
                "author": "=string.concat($activity[rest].result.firstName, $activity[rest].result.lastName)",
                "title": "Five little ducks",
                "price": 19.99
              },
              {
                "author": "=string.concat($activity[rest2].result.firstName, $activity[rest2].result.lastName)",
                "title": "I love trucks",
                "price": 11.99
              }
            ]
          }
        }
      }
    }

1.  Adding `@foreach(source, <optional loopName>)` to indicate iterating on a value
2.  Using `$loop.xxx` to access the current loop data `xxx` is the object field name
3.  Using `$loop[loopName].xxx` to access specific loop data

**Note** You can use any literal, functions, expression in object mappings.

    {
      "books": {
        "mapping": {
          "@foreach($flow.store.books)": {
            "author": "=string.concat($activity[rest].result.firstName, $activity[rest].result.lastName)",
            "title": "Five little ducks",
            "price": 19.99
          }
        }
      }
    }

