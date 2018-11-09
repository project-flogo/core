# Swagger Feature Example

## Install

To install run the following commands:
```
flogo create -f flogo.json
cd swagger
flogo build
```

## Testing

Run:
```
bin/swagger
```

Then open another terminal and run:
```
###  Format : curl http://localhost:1234/swagger/<triggerName> ###

curl http://localhost:1234/swagger/swagdocs
```

You should then see something like:
```
{
    "host": "Temporarys-MacBook-Pro.local: 1234",
    "info": {
        "description": "1.0.0",
        "title": "swagger",
        "version": "This is a simple proxy."
    },
    "paths": {
        "/swagger": {
            "get": {
                "description": "Simple swagger doc Trigger",
                "parameters": [],
                "responses": {
                    "200": {
                        "description": "Simple swagger doc Trigger"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "swagdocs"
                ]
            }
        }
    },
    "swagger": "2.0"
}
```
#####
```
curl http://localhost:1234/swagger/MyProxy
```

You should then see something like:
```
{
    "host": "Temporarys-MacBook-Pro.local: 9096",
    "info": {
        "description": "1.0.0",
        "title": "swagger",
        "version": "This is a simple proxy."
    },
    "paths": {
        "/pets": {
            "get": {
                "description": "Simple REST Trigger",
                "parameters": [],
                "responses": {
                    "200": {
                        "description": "Simple REST Trigger"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "MyProxy"
                ]
            }
        }
    },
    "swagger": "2.0"
}
```
