package terraform

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
)

func (t terraformBackend) install(tfVersion string) (string, error) {
	if tfVersion != "" {
		installDirPath := filepath.Join("/tmp/cloudx/terraform-versions", tfVersion)
		if _, err := os.Stat(filepath.Join(installDirPath, "terraform")); os.IsNotExist(err) {
			err := os.MkdirAll(installDirPath, os.ModePerm)
			if err != nil {
				return "", err
			}

			installer := &releases.ExactVersion{
				Product:    product.Terraform,
				Version:    version.Must(version.NewVersion(tfVersion)),
				InstallDir: installDirPath,
			}

			t.logger.Info("install terraform by specific version")
			res, err := installer.Install(context.Background())
			return res, err
		}

		t.logger.Info("using specific terraform version")
		return filepath.Join(installDirPath, "terraform"), nil
	}

	installDirPath := filepath.Join("/tmp/cloudx/terraform-versions", "latest")
	// TODO: fix error text file busy
	if _, err := os.Stat(filepath.Join(installDirPath, "terraform")); os.IsNotExist(err) {
		err := os.MkdirAll(installDirPath, os.ModePerm)
		if err != nil {
			return "", err
		}

		installer := &releases.LatestVersion{
			Product:    product.Terraform,
			InstallDir: installDirPath,
		}

		t.logger.Info("install terraform by latest version")
		res, err := installer.Install(context.Background())
		return res, err
	}

	t.logger.Info("using latest terraform version")
	return filepath.Join(installDirPath, "terraform"), nil
}
