output "monthly_planning_url" {
  value       = aws_lambda_function_url.monthly_planning.function_url
  description = "Public Lambda function URL"
}
