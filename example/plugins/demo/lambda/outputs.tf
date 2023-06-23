output "arn" {
  value = aws_lambda_function.this.arn
}

output "invoke-arn" {
  value = aws_lambda_function.this.invoke_arn
}
