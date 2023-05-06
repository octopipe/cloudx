terraform {
  required_version = "> 0.8.0"
}

resource "null_resource" "health_check" {

 provisioner "local-exec" {
    
    command = "/bin/bash healthcheck.sh"
  }
}