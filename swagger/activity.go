package swagger

import(
	"fmt"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
)
var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{})
const DefaultPort = "9096"

func init() {
	trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}


// New implements trigger.Factory.New
func (f *Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}
	fmt.Println("config is:",config)
	port := strconv.Itoa(config.Settings["port"].(int))
	if len(port) == 0 {
		port = DefaultPort
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	trigger := &Trigger{
		metadata: f.Metadata(),
		config:   config,
		response: string("Hello this is test for swagger"),
		Server: server,
	}
	mux.HandleFunc("/swagger", trigger.SwaggerHandler)
	//response := Swagger("hostname",config)

	return trigger, nil
}

func (t *Trigger) SwaggerHandler(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "{\"response\":\"Ping successful\"}\n")
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	go func() {
		if err := t.Server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Errorf("Ping service err:", err)
		}
	}()
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	if err := t.Server.Shutdown(nil); err != nil {
		fmt.Errorf("[mashling-ping-service] Ping service error when stopping:", err)
		return err
	}
	return nil
}

/*func Swagger(hostname string, config *trigger.Config) string {
	if config.Ref == "github.com/project-flogo/contrib/trigger/rest" {
		fmt.Print
	}

}*/


/*if trigger.Ref == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest"  {
for _, handler := range trigger.Handlers {
var endpoint swagger.Endpoint
endpoint.Name = trigger.Name
endpoint.Method = handler.Settings["method"].(string)
endpoint.Path = handler.Settings["path"].(string)
endpoint.Description = trigger.Description
var beginDelim, endDelim rune
switch trigger.Type {
case "github.com/TIBCOSoftware/flogo-contrib/trigger/rest":
beginDelim = ':'
endDelim = '/'
case "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger":
beginDelim = '{'
endDelim = '}'
default:
beginDelim = '{'
endDelim = '}'
}
endpoint.BeginDelim = beginDelim
endpoint.EndDelim = endDelim
endpoints = append(endpoints, endpoint)
}
}*/


func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	return nil
}