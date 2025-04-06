resource "aws_api_gateway_rest_api" "ticket_api" {
  name        = "TicketAPI"
  description = "API para gerenciamento de tickets"
}

# Recurso /tickets
resource "aws_api_gateway_resource" "tickets" {
  rest_api_id = aws_api_gateway_rest_api.ticket_api.id
  parent_id   = aws_api_gateway_rest_api.ticket_api.root_resource_id
  path_part   = "tickets"
}

# Recurso /tickets/{id}
resource "aws_api_gateway_resource" "ticket" {
  rest_api_id = aws_api_gateway_rest_api.ticket_api.id
  parent_id   = aws_api_gateway_resource.tickets.id
  path_part   = "{id}"
}

# Método GET para listar todos os tickets
resource "aws_api_gateway_method" "tickets_get" {
  rest_api_id   = aws_api_gateway_rest_api.ticket_api.id
  resource_id   = aws_api_gateway_resource.tickets.id
  http_method   = "GET"
  authorization = "NONE"
}

# Método POST para criar ticket
resource "aws_api_gateway_method" "tickets_post" {
  rest_api_id   = aws_api_gateway_rest_api.ticket_api.id
  resource_id   = aws_api_gateway_resource.tickets.id
  http_method   = "POST"
  authorization = "NONE"
}

# Método GET para obter ticket por ID
resource "aws_api_gateway_method" "ticket_get" {
  rest_api_id   = aws_api_gateway_rest_api.ticket_api.id
  resource_id   = aws_api_gateway_resource.ticket.id
  http_method   = "GET"
  authorization = "NONE"
  request_parameters = {
    "method.request.path.id" = true
  }
}

# Método PUT para atualizar ticket
resource "aws_api_gateway_method" "ticket_put" {
  rest_api_id   = aws_api_gateway_rest_api.ticket_api.id
  resource_id   = aws_api_gateway_resource.ticket.id
  http_method   = "PUT"
  authorization = "NONE"
  request_parameters = {
    "method.request.path.id" = true
  }
}

# Método DELETE para excluir ticket
resource "aws_api_gateway_method" "ticket_delete" {
  rest_api_id   = aws_api_gateway_rest_api.ticket_api.id
  resource_id   = aws_api_gateway_resource.ticket.id
  http_method   = "DELETE"
  authorization = "NONE"
  request_parameters = {
    "method.request.path.id" = true
  }
}

# Integrações com Lambda

# Integração GET /tickets
resource "aws_api_gateway_integration" "tickets_get" {
  rest_api_id = aws_api_gateway_rest_api.ticket_api.id
  resource_id = aws_api_gateway_resource.tickets.id
  http_method = aws_api_gateway_method.tickets_get.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# Integração POST /tickets
resource "aws_api_gateway_integration" "tickets_post" {
  rest_api_id = aws_api_gateway_rest_api.ticket_api.id
  resource_id = aws_api_gateway_resource.tickets.id
  http_method = aws_api_gateway_method.tickets_post.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# Integração GET /tickets/{id}
resource "aws_api_gateway_integration" "ticket_get" {
  rest_api_id = aws_api_gateway_rest_api.ticket_api.id
  resource_id = aws_api_gateway_resource.ticket.id
  http_method = aws_api_gateway_method.ticket_get.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# Integração PUT /tickets/{id}
resource "aws_api_gateway_integration" "ticket_put" {
  rest_api_id = aws_api_gateway_rest_api.ticket_api.id
  resource_id = aws_api_gateway_resource.ticket.id
  http_method = aws_api_gateway_method.ticket_put.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# Integração DELETE /tickets/{id}
resource "aws_api_gateway_integration" "ticket_delete" {
  rest_api_id = aws_api_gateway_rest_api.ticket_api.id
  resource_id = aws_api_gateway_resource.ticket.id
  http_method = aws_api_gateway_method.ticket_delete.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda_invoke_arn
}

# Implantação da API
resource "aws_api_gateway_deployment" "api_deployment" {
  depends_on = [
    aws_api_gateway_integration.tickets_get,
    aws_api_gateway_integration.tickets_post,
    aws_api_gateway_integration.ticket_get,
    aws_api_gateway_integration.ticket_put,
    aws_api_gateway_integration.ticket_delete
  ]

  rest_api_id = aws_api_gateway_rest_api.ticket_api.id
  stage_name  = "prod"
}
