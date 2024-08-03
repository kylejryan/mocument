// internal/db/docdb.go
package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/docdb"
	"github.com/kylejryan/mocument/internal/config"
)

type DocumentDB interface {
	CreateCluster(input *docdb.CreateDBClusterInput) (*docdb.CreateDBClusterOutput, error)
	DeleteCluster(input *docdb.DeleteDBClusterInput) (*docdb.DeleteDBClusterOutput, error)
	CreateInstance(input *docdb.CreateDBInstanceInput) (*docdb.CreateDBInstanceOutput, error)
	DeleteInstance(input *docdb.DeleteDBInstanceInput) (*docdb.DeleteDBInstanceOutput, error)
	InsertDocument(collection string, document interface{}) error
	FindDocument(collection string, filter interface{}) (interface{}, error)
	DeleteDocument(collection string, filter interface{}) error
}

type DocDB struct {
	client *docdb.DocDB
}

func NewDocDB(cfg *config.Config) (*DocDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		return nil, err
	}

	return &DocDB{
		client: docdb.New(sess),
	}, nil
}
