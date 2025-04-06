#!/bin/bash
set -e

echo "Deploy Application Localstack..."

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_ENDPOINT_URL=http://localhost:4566

if ! curl -s http://localhost:4566 > /dev/null; then
    echo "Localstack not running. Please run 'docker-compose up -d'"
    exit 1
fi

echo "Create Zip Lambda..."
zip -j function.zip bootstrap

echo "create/update Lambda function..."
FUNCTION_ARN=$(aws --endpoint-url=http://localhost:4566 lambda list-functions --query "Functions[?FunctionName=='ticket-function'].FunctionArn" --output text)

if [ -z "$FUNCTION_ARN" ]; then
    echo "Create new function'..."
    aws --endpoint-url=http://localhost:4566 lambda create-function \
        --function-name ticket-function \
        --runtime provided.al2 \
        --handler bootstrap \
        --role arn:aws:iam::000000000000:role/lambda-role \
        --zip-file fileb://function.zip \
        --environment "Variables={DB_HOST=postgres,DB_PORT=5432,DB_USER=postgres,DB_PASSWORD=postgres,DB_NAME=tickets_db}" \
        --timeout 30
else
    # Update function
    echo "Update Function'..."
    aws --endpoint-url=http://localhost:4566 lambda update-function-code \
        --function-name ticket-function \
        --zip-file fileb://function.zip
fi

echo "Configure API Gateway..."
API_ID=$(aws --endpoint-url=http://localhost:4566 apigateway get-rest-apis --query "items[?name=='TicketAPI'].id" --output text)

if [ -z "$API_ID" ]; then
    # Create API Gateway
    echo "Create new API Gateway 'TicketAPI'..."
    API_ID=$(aws --endpoint-url=http://localhost:4566 apigateway create-rest-api \
        --name TicketAPI \
        --query "id" \
        --output text)

    # Get ID
    ROOT_RESOURCE_ID=$(aws --endpoint-url=http://localhost:4566 apigateway get-resources \
        --rest-api-id $API_ID \
        --query "items[?path=='/'].id" \
        --output text)

    # Create resource /tickets
    TICKETS_RESOURCE_ID=$(aws --endpoint-url=http://localhost:4566 apigateway create-resource \
        --rest-api-id $API_ID \
        --parent-id $ROOT_RESOURCE_ID \
        --path-part "tickets" \
        --query "id" \
        --output text)

    # Create resource /tickets/{id}
    TICKET_ID_RESOURCE_ID=$(aws --endpoint-url=http://localhost:4566 apigateway create-resource \
        --rest-api-id $API_ID \
        --parent-id $TICKETS_RESOURCE_ID \
        --path-part "{id}" \
        --query "id" \
        --output text)

    # GET - List All tickets
    aws --endpoint-url=http://localhost:4566 apigateway put-method \
        --rest-api-id $API_ID \
        --resource-id $TICKETS_RESOURCE_ID \
        --http-method GET \
        --authorization-type NONE

    # POST - Create New ticket
    aws --endpoint-url=http://localhost:4566 apigateway put-method \
        --rest-api-id $API_ID \
        --resource-id $TICKETS_RESOURCE_ID \
        --http-method POST \
        --authorization-type NONE

    # GET - Get ticket by id
    aws --endpoint-url=http://localhost:4566 apigateway put-method \
        --rest-api-id $API_ID \
        --resource-id $TICKET_ID_RESOURCE_ID \
        --http-method GET \
        --authorization-type NONE \
        --request-parameters "method.request.path.id=true"

    # PUT - Update ticket
    aws --endpoint-url=http://localhost:4566 apigateway put-method \
        --rest-api-id $API_ID \
        --resource-id $TICKET_ID_RESOURCE_ID \
        --http-method PUT \
        --authorization-type NONE \
        --request-parameters "method.request.path.id=true"

    # DELETE - Delete ticket
    aws --endpoint-url=http://localhost:4566 apigateway put-method \
        --rest-api-id $API_ID \
        --resource-id $TICKET_ID_RESOURCE_ID \
        --http-method DELETE \
        --authorization-type NONE \
        --request-parameters "method.request.path.id=true"

    # Configure integrations
    for METHOD in GET POST; do
        aws --endpoint-url=http://localhost:4566 apigateway put-integration \
            --rest-api-id $API_ID \
            --resource-id $TICKETS_RESOURCE_ID \
            --http-method $METHOD \
            --type AWS_PROXY \
            --integration-http-method POST \
            --uri "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:ticket-function/invocations"
    done

    for METHOD in GET PUT DELETE; do
        aws --endpoint-url=http://localhost:4566 apigateway put-integration \
            --rest-api-id $API_ID \
            --resource-id $TICKET_ID_RESOURCE_ID \
            --http-method $METHOD \
            --type AWS_PROXY \
            --integration-http-method POST \
            --uri "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:ticket-function/invocations"
    done

    # Deploy
    aws --endpoint-url=http://localhost:4566 apigateway create-deployment \
        --rest-api-id $API_ID \
        --stage-name prod
else
    echo "API Gateway 'TicketAPI' already exists: $API_ID"
    aws --endpoint-url=http://localhost:4566 apigateway create-deployment \
        --rest-api-id $API_ID \
        --stage-name prod
fi

aws --endpoint-url=http://localhost:4566 lambda add-permission \
    --function-name ticket-function \
    --statement-id apigateway-test \
    --action lambda:InvokeFunction \
    --principal apigateway.amazonaws.com \
    --source-arn "arn:aws:execute-api:us-east-1:000000000000:$API_ID/*/*/*" \
    2>/dev/null || true

echo "Localstack Done!"
echo "API URL: http://localhost:4566/restapis/$API_ID/prod/_user_request_/"
echo "For example: curl -X GET http://localhost:4566/restapis/$API_ID/prod/_user_request_/tickets"
