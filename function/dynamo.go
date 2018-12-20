package function

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	cfg "github.com/alknopfler/alexa-mibebe/config"
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

func getRecordsName(key, value string) ([]RecordName, error){

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
		return nil, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(cfg.DynamoTableName),
	}

	result, err := svc.Scan(params)
	var item []RecordName
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)
	if err != nil {
		log.Println("Failed to unmarshal Record: "+ err.Error())
		return nil, err
	}
	log.Println(item)
	return item, nil

}


func getRecordsPeso(key, value, oldTime, newTime string) ([]RecordPeso, error){

	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return nil, err
	}
	svc := dynamodb.New(sess)

	filt := expression.Name("fecha").Between(expression.Value(oldTime),expression.Value(newTime)).And(expression.Name("email").Equal(expression.Value(value)))
	proj := expression.NamesList(expression.Name("email"), expression.Name("fecha"), expression.Name("nombre"), expression.Name("peso"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("Got error building expression: "+err.Error())
		return nil, err
	}
	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		ScanIndexForward: 		   aws.Bool(false),
		Limit: 					   aws.Int64(1),
		TableName:                 aws.String(cfg.DynamoTablePeso),
	}

	result, err := svc.Query(params)
	if err!= nil{
		log.Println("error scanning")
	}
	var item []RecordPeso
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)
	if err != nil {
		log.Println("Failed to unmarshal Record: "+ err.Error())
		return nil, err
	}
	return item, nil

}


func getRecordsToma(key, value, oldTime, newTime string) ([]RecordToma, error){

	sess, err := session.NewSession(&aws.Config{Region: &cfg.AWS_Region})
	if err != nil{
		log.Println("Error creating session with aws: " + err.Error())
		return nil, err
	}
	svc := dynamodb.New(sess)

	filt := expression.Name("fecha").Between(expression.Value(oldTime),expression.Value(newTime)).And(expression.Name("email").Equal(expression.Value(value)))
	proj := expression.NamesList(expression.Name("email"), expression.Name("fecha"), expression.Name("nombre"), expression.Name("toma"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Println("Got error building expression: "+err.Error())
		return nil, err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(cfg.DynamoTableToma),
	}

	result, err := svc.Scan(params)
	if err!= nil{
		log.Println("error scanning")
	}
	var item []RecordToma
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &item)
	if err != nil {
		log.Println("Failed to unmarshal Record: "+ err.Error())
		return nil, err
	}
	return item, nil

}
