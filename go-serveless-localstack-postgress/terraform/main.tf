provider "aws" {
  access_key                  = "test"
  secret_key                  = "test"
  region                      = "us-east-1"
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway     = "http://localhost:4566"
    lambda         = "http://localhost:4566"
    iam            = "http://localhost:4566"
    s3             = "http://localhost:4566"
  }
}

module "lambda" {
  source = "./modules/lambda"
}

module "api_gateway" {
  source          = "./modules/api_gateway"
  lambda_function = module.lambda.function_name
  lambda_invoke_arn = module.lambda.invoke_arn

}

output "api_endpoint" {
  value = "http://localhost:4566/restapis/${module.api_gateway.api_id}/prod/_user_request_/"
}