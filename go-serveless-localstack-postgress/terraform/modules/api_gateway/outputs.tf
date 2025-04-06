output "api_id" {
  description = "ID da API REST do API Gateway"
  value       = aws_api_gateway_rest_api.ticket_api.id
}