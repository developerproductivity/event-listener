package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	eventType "github.com/event-listener/types"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// DynoObject represents an object in dynamoDB.
// Used to represent key value data such as keys, table items...
type DynoNotation map[string]types.AttributeValue

// newclient constructs a new dynamodb client using a default configuration
// and a provided profile name (created via aws configure cmd).
func Newclient() (*dynamodb.Client, error) {
	region := os.Getenv("REGION")
	url := os.Getenv("URL")
	accsKeyID := os.Getenv("ACCESSKEYID")
	secretAccessKey := os.Getenv("SECRETACCESSKEY")
	fmt.Println(url, "URL")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: url}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: accsKeyID, SecretAccessKey: secretAccessKey, SessionToken: "",
				Source: "Mock credentials used above for local instance",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	c := dynamodb.NewFromConfig(cfg)
	return c, nil
}

// createTable creates a table in the client's dynamodb instance
// using the table parameters specified in input.
func createTable(c *dynamodb.Client,
	tableName string, input *dynamodb.CreateTableInput,
) error {
	var tableDesc *types.TableDescription
	table, err := c.CreateTable(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to create table `%v` with error: %v\n", tableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(c)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Failed to wait on create table `%v` with error: %v\n", tableName, err)
		}
		tableDesc = table.TableDescription
	}
	fmt.Printf("Created table `%s` with details: %v\n\n", tableName, tableDesc)

	return err
}

// listTables returns a list of table names in the client's dynamodb instance.
func listTables(c *dynamodb.Client, input *dynamodb.ListTablesInput) ([]string, error) {
	tables, err := c.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{},
	)
	if err != nil {
		return nil, err
	}

	return tables.TableNames, nil
}

// putItem inserts an item (key + attributes) in to a dynamodb table.
func putItem(c *dynamodb.Client, tableName string, item DynoNotation) (err error) {
	_, err = c.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetCiBuildPayload(client *dynamodb.Client) []eventType.CiBuildPayload {
	var payload []eventType.CiBuildPayload
	originAttr, _ := attributevalue.Marshal("Tekton")
	keyExpr := expression.Key("origin").Equal(expression.Value(originAttr))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpr).Build()
	if err != nil {
		log.Fatal(err)
	}
	query, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String("TektonCI"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		log.Fatal(err)
	}
	// unmarshal list of items
	err = attributevalue.UnmarshalListOfMaps(query.Items, &payload)
	if err != nil {
		log.Fatal(err)
	}
	return payload
}

func InsertRecordInDatabase(object v1.PipelineRun, client *dynamodb.Client) {

	item := PrepareCiBuildData(object)
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		fmt.Println("failed to marshal Record, %w", err)
		return
	}
	fmt.Println("Response from put api ", putItem(client, "TektonCI", av))
}

func PrepareCiBuildData(obj v1.PipelineRun) eventType.CiBuildPayload {
	payload := eventType.CiBuildPayload{
		Origin:          "Tekton",
		OriginalID:      string(obj.UID),
		Name:            obj.Name,
		URL:             obj.Status.Provenance.RefSource.URI,
		CreatedAt:       obj.Status.StartTime.Time.Unix(),
		StartedAt:       obj.Status.StartTime.Time.Unix(),
		CompletedAt:     obj.Status.CompletionTime.Time.Unix(),
		Status:          string(obj.Status.Conditions[0].Type),
		Conclusion:      string(obj.Status.Conditions[0].Status),
		RepoURL:         obj.Status.Provenance.RefSource.URI,
		Commit:          "",
		PullRequestUrls: make([]string, 0),
		IsDeployment:    true,
	}
	triggeredBy := eventType.TriggeredBy{
		Name:         "Pipelines Operator",
		Email:        "dummy@redhat.com",
		AccountId:    "dummy@redhat.com",
		LastActivity: obj.Status.Conditions[0].LastTransitionTime.Inner.Unix(),
	}
	payload.TriggeredBy = triggeredBy
	var dynamicClientSet *dynamic.DynamicClient
	var err error
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Errorf("Fail to build the k8s config. Error - %s", err)
		return eventType.CiBuildPayload{}
	}
	// inorder to create the dynamic Client set
	dynamicClientSet, err = dynamic.NewForConfig(config)
	if err != nil {
		fmt.Errorf("Fail to create the dynamic client set. Errorf - %s", err)
		return eventType.CiBuildPayload{}
	}
	genericSchema := schema.GroupVersionResource{
		Group:    "tekton.dev",
		Version:  "v1",
		Resource: "taskruns",
	}
	dinterface := dynamicClientSet.Resource(genericSchema).Namespace(obj.Namespace)

	var tasks []eventType.Job
	for _, val := range obj.Status.ChildReferences {
		if val.Kind == "TaskRun" {
			var tr *unstructured.Unstructured
			tr, err = dinterface.Get(context.TODO(), val.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Printf("Error retreiving task run %v %v", val.Name, err.Error())
			}
			unstructured := tr.UnstructuredContent()
			var task v1.TaskRun
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructured, &task)
			if err != nil {
				fmt.Printf("Error converting to task run %v", val.Name)
			}
			job := eventType.Job{
				StartedAt:   task.Status.StartTime.Time.Unix(),
				CompletedAt: task.Status.CompletionTime.Time.Unix(),
				Name:        task.Name,
				Status:      string(task.Status.Conditions[0].Status),
				Conclusion:  task.Status.Conditions[0].Reason,
			}
			tasks = append(tasks, job)
		}
	}

	var stg []eventType.Stage
	stage := eventType.Stage{
		ID:          string(obj.UID),
		Name:        obj.Name,
		StartedAt:   obj.Status.StartTime.Time.Unix(),
		CompletedAt: obj.Status.CompletionTime.Time.Unix(),
		Status:      string(obj.Status.Conditions[0].Status),
		Conclusion:  obj.Status.Conditions[0].Reason,
		URL:         obj.Status.Provenance.RefSource.URI,
		Jobs:        tasks,
	}
	stg = append(stg, stage)
	payload.Stages = stg
	return payload
}
