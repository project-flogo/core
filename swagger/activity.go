package swagger

import(
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"github.com/project-flogo/core/api"
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

// Trigger is the swagger trigger
type Trigger struct {
	metadata 	*trigger.Metadata
	settings 	*Settings
	config   	*trigger.Config
	Server 		*http.Server
	logger 		log.Logger
	response	string
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
	fmt.Println("Calling Swagger")
	response,_ := Swagger("hostname",config)
	fmt.Println("After Swagger")


	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	trigger := &Trigger{
		metadata: f.Metadata(),
		config:   config,
		response: string(response),
		Server: server,
	}
	mux.HandleFunc("/swagger", trigger.SwaggerHandler)

	return trigger, nil
}

func (t *Trigger) SwaggerHandler(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, t.response)
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

func Swagger(hostname string, config *trigger.Config) ([]byte, error) {
	fmt.Println("Inside Swagger")
	var endpoints []Endpoint
	appConfig := *api.App.properties
	fmt.Println("Details : ",appConfig)
	if config.Ref == "github.com/project-flogo/contrib/trigger/rest" {
		for _, handler := range config.Handlers{
			fmt.Println("Inside Swagger : for")
			var endpoint Endpoint
			endpoint.Name = config.Id
			fmt.Println("ID")
			endpoint.Method = handler.Settings["method"].(string)
			fmt.Println("Method")
			endpoint.Path = handler.Settings["path"].(string)
			fmt.Println("Path")
			endpoint.Description = config.Settings["description"].(string)
			fmt.Println("Description")
			var beginDelim, endDelim rune
			switch config.Ref {
			case "github.com/project-flogo/contrib/trigger/rest":
				beginDelim = ':'
				endDelim = '/'
			default:
				beginDelim = '{'
				endDelim = '}'
			}
			endpoint.BeginDelim = beginDelim
			endpoint.EndDelim = endDelim
			endpoints = append(endpoints, endpoint)
		}
	}
	fmt.Println("Before Generate")
	/*byteArray,error := Generate(hostname, appConfig.Name, appConfig.Description, appConfig.Version, endpoints)
	if error != nil {
		fmt.Println("Inside error")
		return "",error
	}
	fmt.Println("Before return")
	return string(byteArray[:]), nil*/
	fmt.Println("App name :", appConfig.Name())
	//fmt.Println("Description :", appConfig.Description)
	fmt.Println("Version :", appConfig.Version())
	return Generate(hostname, appConfig.Name, appConfig.Description, appConfig.Version, endpoints)
}


func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	return nil
}