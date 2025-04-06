resource "aws_iam_role" "lambda_role" {
  name = "lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "ticket_function" {
  filename         = "${path.module}/../../../function.zip"
  function_name    = "ticket-function"
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  runtime          = "provided.al2"
  source_code_hash = filebase64sha256("${path.module}/../../../function.zip")
  timeout          = 30

  environment {
    variables = {
      DB_HOST     = "postgres"
      DB_PORT     = "5432"
      DB_USER     = "postgres"
      DB_PASSWORD = "postgres"
      DB_NAME     = "tickets_db"
    }
  }
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ticket_function.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:us-east-1:000000000000:*/*/*"
}

output "function_name" {
  value = aws_lambda_function.ticket_function.function_name
}

output "invoke_arn" {
  value = aws_lambda_function.ticket_function.invoke_arn
}