package swagger


// Endpoint represents an endpoint in a Swagger 2.0 document.
type Endpoint struct {
	Name        string `md:"name"`
	Description string `md:"description"`
	Path        string `md:"path"`
	Method      string `md:"method"`
	BeginDelim  rune   `md:"begin_delim"`
	EndDelim    rune   `md:"end_delim"`
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


type Settings struct {
	Port int `md:"port,required"`
}

type HandlerSettings struct {
	Method string `md:"method,required,allowed(GET,POST,PUT,PATCH,DELETE)"`
	Path   string `md:"path,required"`
}
