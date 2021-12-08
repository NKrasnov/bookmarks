package main

import (
	"fmt"
	"os"

	"github.com/nkrasnov/bookmarks/server/foundation/config"
)

func main() {

	_ = run(os.Args)
}

func run(args []string) error {

	// App configuration struct
	// tags description
	// cmd     - command line parameter name
	// env     - environment variable name
	// default - default value for the parameter.
	//           command line parameters take precedence over environment variables
	// usage   - short description
	cfg := struct {
		APIHost        string `param:"cmd=host,env=BM_API_HOST,default=127.0.0.1,usage=IP or DNS Name"`
		APIPort        int    `param:"cmd=port,env=BM_API_PORT,default=8081,usage=API server port"`
		DBHost         string `param:"cmd=dbhost,env=BM_API_DBHOST,default=127.0.0.1,usage=IP or DNS Name"`
		DBPort         int    `param:"cmd=dbport,env=BM_API_DBPORT,default=5432,usage=Port number a database server is listens to"`
		DBUser         string `param:"cmd=dbuser,env=BM_API_DBUSER,default=postgres,usage=database user name"`
		DBPwd          string `param:"cmd=dbpwd,env=BM_API_DBPWD,usage=database user password"`
		ReadTimeout    int    `param:"cmd=srto,env=BM_API_READ_TIMEOUT,default=10,usage=API server read timeout"`
		WriteTimeout   int    `param:"cmd=swto,env=BM_API_WRITE_TIMEOUT,default=10,usage=API server write timeout"`
		RequestTimeout int    `param:"cmd=rto,env=BM_API_REQUEST_TIMEOUT,default=5,usage=API request timeout"`
	}{}
	err := config.Parse(&cfg, args)

	//if errors.Is(err, config.ErrHelpNeeded) {
	config.PrintUsage()
	//	return nil
	//}
	if err != nil {
		fmt.Println(err)
	}
	config.PrintUsage()
	return nil
}
