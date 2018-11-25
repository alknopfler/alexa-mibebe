package function

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	cfg "github.com/alknopfler/alexa-mibebe/config"
	"errors"
)

func createRecord(data interface{}, dbTable string) error {
	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return err
	}
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(&data)
	if err != nil {
		log.Println("Got error marshalling map: "+err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dbTable),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Println("Error inserting element: "+err.Error())
		return err
	}

	return nil
}


func existRecord(key, value, dbTable string) (bool,error) {
	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return false, err
	}
	svc := dynamodb.New(sess)

	filt := expression.Name(key).Equal(expression.Value(value))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		log.Println("Got error building expression: "+err.Error())
		return false, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(dbTable),
	}

	result, err := svc.Scan(params)

	if err != nil {
		log.Println("Query API call failed: "+err.Error())
		return false, err
	}

	if len(result.Items) > 0 {
		return true, nil
	}
	return false, nil
}

func getRecord(key, value, dbTable string) (interface{}, error){

	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return nil, err
	}
	svc := dynamodb.New(sess)

	filt := expression.Name(key).Equal(expression.Value(value))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		log.Println("Got error building expression: "+err.Error())
		return false, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(dbTable),
	}

	result, err := svc.Scan(params)

	switch dbTable {
	case cfg.DynamoTableName:
		var item RecordName
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
		if err != nil {
			log.Println("Failed to unmarshal Record: "+ err.Error())
			return nil, err
		}
		log.Println("EL ITEM ES: ")
		log.Println(item)
		return item, nil
	case cfg.DynamoTablePeso:
		var item RecordPeso
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
		if err != nil {
			log.Println("Failed to unmarshal Record: "+ err.Error())
			return nil, err
		}
		log.Println(item)
		return item, nil
	case cfg.DynamoTableToma:
		var item RecordToma
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
		if err != nil {
			log.Println("Failed to unmarshal Record: "+ err.Error())
			return nil, err
		}
		return item, nil
	default:
		return nil, errors.New("Error item not found")
	}
}