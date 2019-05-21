package generate_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	_ "github.com/project-flogo/core/examples/action"
)

var testFlogoJSON = `{
  "name": "Test",
  "type": "flogo:app",
  "version": "1.0.0",
  "description": "This is a test application.",
	"imports": [
		"github.com/project-flogo/core/examples/action",
		"github.com/project-flogo/core/examples/trigger",
		"_ github.com/project-flogo/core/data/expression/script"
	],
  "properties": [
		{"name": "test0", "type": "string", "value": "test"},
		{"name": "test1", "type": "int", "value": 1},
		{"name": "test2", "type": "bool", "value": true}
	],
  "channels": [
    "test0:1",
		"test1:2",
		"test2:3"
  ],
  "triggers": [
    {
      "name": "flogo-test0",
      "id": "test0",
      "ref": "#trigger",
      "settings": {
        "aSetting": 123
      },
      "handlers": [
        {
          "settings": {
            "aSetting": 123
          },
          "actions": [
            {
              "id": "action:Test0"
            }
          ]
        }
      ]
    },
		{
      "name": "flogo-test1",
      "id": "test1",
      "ref": "#trigger",
      "settings": {
        "aSetting": 123
      },
      "handlers": [
        {
          "settings": {
            "aSetting": 123
          },
          "actions": [
            {
              "id": "action:Test1"
            }
          ]
        }
      ]
    },
		{
      "name": "flogo-test1",
      "id": "test1",
      "ref": "github.com/project-flogo/core/examples/trigger",
      "settings": {
        "aSetting": 123
      },
      "handlers": [
        {
          "settings": {
            "aSetting": 123
          },
          "actions": [
            {
							"if": "1 == 1",
              "id": "action:Test1",
							"input": {
								"test0": "=1",
								"test1": "=2",
								"test2": "=3"
							},
							"output": {
								"test0": "=1",
								"test1": "=2",
								"test2": "=3"
							}
            },
						{
							"ref": "github.com/project-flogo/core/examples/action",
				      "settings": {
				        "aSetting": "action:Test"
				      }
						}
          ]
        }
      ]
    }
  ],
  "resources": [
    {
      "id": "action:Test",
      "compressed": false,
      "data": {
				"message": "hello world"
			}
    }
  ],
  "actions": [
    {
      "ref": "github.com/project-flogo/core/examples/action",
      "settings": {
        "aSetting": "action:Test"
      },
      "id": "action:Test0",
      "metadata": null
    },
		{
      "ref": "github.com/project-flogo/core/examples/action",
      "settings": {
        "aSetting": "action:Test"
      },
      "id": "action:Test1",
      "metadata": null
    }
  ]
}`

func TestGenerate(t *testing.T) {
	app, err := engine.LoadAppConfig(testFlogoJSON, false)
	if err != nil {
		t.Fatal(err)
	}
	os.Mkdir("test", 0777)
	api.Generate(app, "./test/test.go")
	cmd := exec.Command("go", "build", "-o", "./test/test", "./test/")
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	err = os.RemoveAll("./test")
	if err != nil {
		t.Fatal(err)
	}
}
