resource "aws_sns_topic" "this" {
  name = var.name
  tags = {
    owner = "stackspot-1"
  }
}