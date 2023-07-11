package terraform

import (
	"context"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
)

func (t terraformBackend) install(tfVersion string) (string, error) {
	installDirPath := "/tmp/cloudx/terraform-versions"
	err := os.MkdirAll(installDirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	if tfVersion != "" {
		installer := &releases.ExactVersion{
			Product:    product.Terraform,
			Version:    version.Must(version.NewVersion(tfVersion)),
			InstallDir: installDirPath,
		}
		return installer.Install(context.Background())
	}

	installer := &releases.LatestVersion{
		Product:    product.Terraform,
		InstallDir: installDirPath,
	}

	return installer.Install(context.Background())
}
