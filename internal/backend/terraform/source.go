package terraform

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (p terraformBackend) ociDownload(sourceUrl string) (string, error) {
	p.logger.Info("pulling task image", zap.String("image", sourceUrl))
	img, err := crane.Pull(sourceUrl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = crane.Export(img, &buf)
	if err != nil {
		return "", err
	}

	tr := tar.NewReader(&buf)
	content := map[string]string{}
	workdir := fmt.Sprintf("/tmp/cloudx/executions/%s", uuid.New().String())
	err = os.MkdirAll(workdir, os.ModePerm)
	if err != nil {
		return "", err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		b, err := io.ReadAll(tr)
		if err != nil {
			return "", err
		}

		content[hdr.Name] = string(b)
		f, err := os.Create(fmt.Sprintf("%s/%s", workdir, hdr.Name))
		if err != nil {
			return "", err
		}

		_, err = f.Write(b)
		if err != nil {
			return "", err
		}
	}

	return workdir, nil
}

func (t terraformBackend) dowloadSource(source string) (string, error) {
	s := strings.Split(source, "://")
	if len(s) <= 1 {
		return "", fmt.Errorf("invalid source. Plese use protocol://source-url.")
	}

	protocol, sourceUrl := s[0], s[1]
	switch protocol {
	case "git":
		return "", nil
	case "s3":
		return "", nil
	case "oci":
		return t.ociDownload(sourceUrl)
	default:
		return "", fmt.Errorf("Invalid protocol")
	}
}
