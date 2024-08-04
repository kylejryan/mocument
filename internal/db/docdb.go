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
	UpdateMany(collection string, filter, update interface{}) error
	FindDocument(collection string, filter interface{}) (interface{}, error)
	DeleteDocument(collection string, filter interface{}) error
}

type DocDB struct {
	client *docdb.DocDB
}

// CreateCluster implements DocumentDB.
func (*DocDB) CreateCluster(input *docdb.CreateDBClusterInput) (*docdb.CreateDBClusterOutput, error) {
	panic("unimplemented")
}

// CreateInstance implements DocumentDB.
func (*DocDB) CreateInstance(input *docdb.CreateDBInstanceInput) (*docdb.CreateDBInstanceOutput, error) {
	panic("unimplemented")
}

// DeleteCluster implements DocumentDB.
func (*DocDB) DeleteCluster(input *docdb.DeleteDBClusterInput) (*docdb.DeleteDBClusterOutput, error) {
	panic("unimplemented")
}

// DeleteDocument implements DocumentDB.
func (*DocDB) DeleteDocument(collection string, filter interface{}) error {
	panic("unimplemented")
}

// DeleteInstance implements DocumentDB.
func (*DocDB) DeleteInstance(input *docdb.DeleteDBInstanceInput) (*docdb.DeleteDBInstanceOutput, error) {
	panic("unimplemented")
}

// FindDocument implements DocumentDB.
func (*DocDB) FindDocument(collection string, filter interface{}) (interface{}, error) {
	panic("unimplemented")
}

// InsertDocument implements DocumentDB.
func (*DocDB) InsertDocument(collection string, document interface{}) error {
	panic("unimplemented")
}

func (db *DocDB) UpdateMany(collection string, filter, update interface{}) error {
	// Implementation for real DocumentDB
	return nil
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
