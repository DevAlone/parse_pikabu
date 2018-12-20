package parser

type ParserConfig struct {
	ParserId                         string
	ApiURL                           string
	ApiTimeout                       int
	ApiSessionId                     string
	ProxyProviderAPIURL              string
	ProxyProviderTimeout             int
	PikagoTimeout                    uint
	PikagoWaitBetweenProcessingPages int
	WaitAfterErrorSeconds            int
	WaitNoTaskSeconds                int
	AMQPAddress                      string
}

func NewParserConfigFromFile(filename string) (*ParserConfig, error) {
	/*file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ParserConfig)

	return err*/
	// TODO: complete
	config := &ParserConfig{}

	config.ParserId = "unique_parser_id"
	config.ApiURL = "http://localhost:8080/api/v1"
	config.ProxyProviderAPIURL = "https://eivailohciihi4uquapach7abei9iesh.d3d.info/api/v1/"
	config.ProxyProviderTimeout = 60
	config.PikagoTimeout = 5
	config.PikagoWaitBetweenProcessingPages = 1
	config.ApiTimeout = 60
	config.WaitAfterErrorSeconds = 5
	config.WaitNoTaskSeconds = 2
	config.ApiSessionId = "parser_oogoShaituNoh8iebaesiYaeh"
	config.AMQPAddress = "amqp://guest:guest@localhost:5672"

	return config, nil
}
