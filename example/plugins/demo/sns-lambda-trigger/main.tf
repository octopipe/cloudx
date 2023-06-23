data "aws_lambda_function" "lambda" {
  function_name = var.lambda_name
}

resource "aws_lambda_permission" "lambda_sns_permission" {
  statement_id  = "AllowExecutionFromSNS"
  action        = "lambda:InvokeFunction"
  function_name =  var.lambda_name
  principal     = "sns.amazonaws.com"
  source_arn    = var.sns_arn
}


resource "aws_sns_topic_subscription" "lambda_source_mapping" {
  topic_arn = var.sns_arn
  protocol  = "lambda"
  endpoint  = data.aws_lambda_function.lambda.arn
}