package swagger

import(
	"fmt"
	"io"
	"net/http"
	"strconv"
	"os"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"github.com/gorilla/mux"
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
	port := strconv.Itoa(config.Settings["port"].(int))
	if len(port) == 0 {
		port = DefaultPort
	}

	//mux := http.NewServeMux()
	newMux := mux.NewRouter()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: newMux,
	}
	trigger := &Trigger{
		metadata: f.Metadata(),
		config:   config,
		response: "",
		Server: server,
	}
	newMux.HandleFunc("/{triggerName}/swagger/", trigger.SwaggerHandler)

	return trigger, nil
}

func (t *Trigger) SwaggerHandler(w http.ResponseWriter, req *http.Request) {

	//mux := t.Server["Handler"]
	//vars := mux.Vars(req)
	fmt.Println("%+v\n",t.Server)
	hostName, err := os.Hostname()
	if err != nil {
		fmt.Errorf("Error in getting hostname:", err)
	}
	//triggerName := vars["triggerName"]
	response,_ := Swagger(hostName,t.config,"hello")
	io.WriteString(w, string(response))
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

func Swagger(hostname string, config *trigger.Config, triggerName string) ([]byte, error) {
	var endpoints []Endpoint
	for _, tConfig := range config.AppConfig["Trigger"].([]*trigger.Config) {
		if tConfig.Id == "" || tConfig.Id == triggerName {
			if tConfig.Ref == "github.com/project-flogo/contrib/trigger/rest" || tConfig.Ref == "github.com/project-flogo/core/swagger" {
				for _, handler := range tConfig.Handlers {
					var endpoint Endpoint
					endpoint.Name = tConfig.Id
					endpoint.Method = handler.Settings["method"].(string)
					endpoint.Path = handler.Settings["path"].(string)
					endpoint.Description = tConfig.Settings["description"].(string)
					var beginDelim, endDelim rune
					switch tConfig.Ref {
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
		}
	}
	return Generate(hostname, config.AppConfig["Name"].(string), config.AppConfig["Version"].(string), config.AppConfig["Description"].(string), endpoints)
}


func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	return nil
}