resource "aws_lambda_function" "this" {
  function_name = var.name
  role          = var.role_arn
  image_uri = var.image_uri
  package_type  = "Image"
  architectures = ["x86_64"]


  environment {
    variables = {
    }
  }
}