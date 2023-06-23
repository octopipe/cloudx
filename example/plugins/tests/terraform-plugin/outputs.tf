output "uid" {
  value = kubectl_manifest.test.uid
  sensitive = true
}