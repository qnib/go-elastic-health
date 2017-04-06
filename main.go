package main

import (
	"os"
	"fmt"
	"log"
	"strings"

	"github.com/urfave/cli"
	"github.com/ddliu/go-httpclient"
)



func main() {
	app := cli.NewApp()
	app.Name = "go-elastic-health"
	app.Usage = "Checks if an elasticsearch node/cluster is healthy"
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "host",
			Value: "localhost",
			Usage: "Hostname/IP of the host to query",
		},
		cli.IntFlag{
			Name: "port",
			Value: 9200,
			Usage: "Port of elasticsearch node",
		},
	}
	app.Action = func(c *cli.Context) (err error) {
		httpclient.Defaults(httpclient.Map {
			httpclient.OPT_USERAGENT: "go-elastic-health",
			"Accept-Language": "en-us",
		})
		addr := fmt.Sprintf("http://%s:%d/_cat/health?h=status", c.String("host"), c.Int("port"))
		log.Printf("Check URL: %s", addr)
		res, err := httpclient.Get(addr, nil)
		if err != nil {
			log.Printf("Query not successful: %v", err.Error())
			os.Exit(1)
		}
		bodyString,err := res.ToString()
		cStatus := strings.TrimSpace(bodyString)
		log.Printf("%d | Status;%s", res.StatusCode, cStatus)
		if res.StatusCode != 200 {
			log.Printf("Return code %d != 200", res.StatusCode)
			os.Exit(res.StatusCode)
		}
		switch cStatus {
		case "green":
			log.Println("Cluster is green...")
			os.Exit(0)
		case "red":
			log.Println("Cluster is RED. panic!")
			os.Exit(1)
		case "yellow":
			log.Println("Cluster is yellow... should be alright...")
			os.Exit(0)
		default:
			log.Printf("Cluster-status is '%s'?; Neigher 'green', 'yellow' or 'red'", cStatus)
			os.Exit(1)
		}
		return
		/*
		if [ ${STATUS} = "green" ];then
  			echo "Cluster is green, all fine..."
  			exit 0
		elif [ ${STATUS} = "red" ];then
		  echoerr "Cluster is red!... "
		  exit 1
		elif [ ! -f /tmp/es_initialised ] && [ ${STATUS} = "yellow" ];then
		  echoerr "Node has not initialised, and cluster is yellow -> considered an error..."
		  exit 1
		elif [ ${STATUS} = "yellow" ];then
		  echo "Cluster is yellow but we are already over the initialisation phase and assume someone else was updated."
		  exit 0
		else
		  echoerr "Cluster is neither green,yello or red, but '${STATUS}'"
		  exit 1
		fi
		 */
	}


	app.Run(os.Args)
}
