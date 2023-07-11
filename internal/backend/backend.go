package backend

import "github.com/octopipe/cloudx/internal/backend/terraform"

const (
	TerraformBackend = "terraform"
)

type Backend struct {
	Terraform terraform.TerraformBackend
}

func NewBackend(terraform terraform.TerraformBackend) Backend {
	return Backend{
		Terraform: terraform,
	}
}
