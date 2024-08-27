package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	api "github.com/event-listener/api"
	db "github.com/event-listener/database"
	"github.com/kelseyhightower/envconfig"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type Data struct {
	Pipelinerun v1.PipelineRun `json:"pipelineRun"`
}

func eventReceiver(ctx context.Context, event cloudevents.Event) error {
	var dat Data
	if err := json.Unmarshal(event.DataEncoded, &dat); err != nil {
		fmt.Println("Ignore")
	}
	var table *db.Table
	var err error

	if table, err = db.Newclient("TektonCI"); err != nil {
		log.Fatalf("failed to create dynamoclient: %s", err.Error())
	}
	fmt.Println("Pipleine run", dat.Pipelinerun)
	table.InsertRecordInDatabase(dat.Pipelinerun)
	return nil
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	log.Print("Starting Event Listener")
	ctx := context.Background()

	p, err := cloudevents.NewHTTP(cloudevents.WithPort(env.Port), cloudevents.WithPath(env.Path))
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}
	c, err := cloudevents.NewClient(p,
		cloudevents.WithUUIDs(),
		cloudevents.WithTimeNow(),
	)
	if err != nil {
		log.Fatalf("failed to create client: %s", err.Error())
	}

	var table *db.Table

	if table, err = db.Newclient("TektonCI"); err != nil {
		log.Fatalf("failed to create dynamoclient: %s", err.Error())
	}

	go func() {
		t := time.Tick(60 * time.Minute)
		for {
			select {
			case <-t:
				fmt.Println("Logilica Upload")
				LogilicaUpload(table)
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Printf("listening on :%d%s\n", env.Port, env.Path)
	if err := c.StartReceiver(ctx, eventReceiver); err != nil {
		log.Fatalf("failed to start receiver: %s", err.Error())
	}

	<-ctx.Done()
}

func LogilicaUpload(table *db.Table) {
	payload := db.GetCiBuildPayload(table.TableName, table.DynamoDbClient)
	if err := api.UploadPlanningData("872a7985dd8a58328dea96015b738c317039fb5a", payload); err == nil {
		db.UpdateRecords(payload, table.TableName, table.DynamoDbClient)
	}
}
