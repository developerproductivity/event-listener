package database

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/dynamodb/stubs"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func enterTest() (*testtools.AwsmStubber, *Table) {
	stubber := testtools.NewStubber()
	basics := &Table{TableName: "test-tekton", DynamoDbClient: dynamodb.NewFromConfig(*stubber.SdkConfig)}
	return stubber, basics
}

func TestTableBasics_ListTables(t *testing.T) {
	t.Run("NoErrors", func(t *testing.T) { ListTables(nil, t) })
	t.Run("TestError", func(t *testing.T) { ListTables(&testtools.StubError{Err: errors.New("TestError")}, t) })
}

func ListTables(raiseErr *testtools.StubError, t *testing.T) {
	stubber, basics := enterTest()

	tableNames := []string{"Table 1", "Table 2", "Table 3"}

	stubber.Add(stubs.StubListTables(tableNames, raiseErr))

	tables, err := basics.ListTables()

	testtools.VerifyError(err, raiseErr, t)
	if err == nil {
		if !reflect.DeepEqual(tables, tableNames) {
			t.Errorf("got %v, expected %v", tables, tableNames)
		}
	}

	testtools.ExitTest(stubber, t)
}

func TestTableBasics_InsertRecordInDatabase(t *testing.T) {
	//t.Run("NoErrors", func(t *testing.T) { InsertRecordInDatabase(nil, t) })
	t.Run("TestError", func(t *testing.T) { ListTables(&testtools.StubError{Err: errors.New("TestError")}, t) })
}

func InsertRecordInDatabase(raiseErr *testtools.StubError, t *testing.T) {
	stubber, basics := enterTest()
	stubber.Add(stubs.StubCreateTable(basics.TableName, raiseErr))
	stubber.Add(stubs.StubDescribeTable(basics.TableName, raiseErr))
	var object v1.PipelineRun
	var dummyTime *metav1.Time
	object.Name = "example"
	object.Status.StartTime = dummyTime
	object.Status.CompletionTime = dummyTime
	object.Status.Provenance.RefSource.URI = "https://example.com"
	err := basics.InsertRecordInDatabase(object)
	testtools.VerifyError(err, raiseErr, t)
	testtools.ExitTest(stubber, t)
}

func TestTableBasics_DeleteTable(t *testing.T) {
	t.Run("NoErrors", func(t *testing.T) { DeleteTable(nil, t) })
	t.Run("TestError", func(t *testing.T) { DeleteTable(&testtools.StubError{Err: errors.New("TestError")}, t) })
}

func DeleteTable(raiseErr *testtools.StubError, t *testing.T) {
	stubber, basics := enterTest()

	stubber.Add(stubs.StubDeleteTable(basics.TableName, raiseErr))

	err := basics.DeleteTable()

	testtools.VerifyError(err, raiseErr, t)
	testtools.ExitTest(stubber, t)
}
