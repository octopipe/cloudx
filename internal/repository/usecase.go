package repository

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
	"github.com/octopipe/cloudx/internal/secret"
	"go.uber.org/zap"
)

type useCase struct {
	logger        *zap.Logger
	repository    RepositoryType
	secretUseCase secret.UseCase
}

func NewUseCase(logger *zap.Logger, repository RepositoryType, secretUseCase secret.UseCase) UseCase {
	return useCase{logger: logger, repository: repository, secretUseCase: secretUseCase}
}

// Create implements UseCase.
func (u useCase) Create(ctx context.Context, r Repository) (Repository, error) {
	_, err := u.secretUseCase.Apply(ctx, secret.Secret{
		Name:      r.Name,
		Namespace: r.Namespace,
		Data: map[string][]byte{
			"type":       []byte(r.Auth.Type),
			"username":   []byte(r.Auth.Username),
			"password":   []byte(r.Auth.Password),
			"publickKey": []byte(r.Auth.PublicKey),
		},
	})
	if err != nil {
		return Repository{}, err
	}

	r.RepositorySpec.AuthRef = commonv1alpha1.Ref{
		Name:      r.Name,
		Namespace: r.Namespace,
	}

	newRepository := commonv1alpha1.Repository{
		Spec: r.RepositorySpec,
	}

	newRepository.SetName(r.Name)
	newRepository.SetNamespace(r.Namespace)

	s, err := u.repository.Apply(ctx, newRepository)
	if err != nil {
		return Repository{}, err
	}

	return Repository{
		Name:           s.GetName(),
		Namespace:      s.GetNamespace(),
		RepositorySpec: s.Spec,
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, name string, namespace string) (Repository, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return Repository{}, err
	}

	return Repository{
		Name:           s.GetName(),
		Namespace:      s.GetNamespace(),
		RepositorySpec: s.Spec,
	}, nil
}

func (u useCase) Sync(ctx context.Context, name string, namespace string) error {
	currRepository, err := u.Get(ctx, name, namespace)
	if err != nil {
		return err
	}

	err = u.syncRemote(ctx, currRepository)
	if err != nil {
		return err
	}

	return nil
}

func (u useCase) syncRemote(ctx context.Context, currRepository Repository) error {
	tmpDir := os.Getenv("TMP_DIR")
	repoDir := fmt.Sprintf("%s/%s", tmpDir, currRepository.Url)

	secret, err := u.secretUseCase.Get(ctx, currRepository.AuthRef.Name, currRepository.AuthRef.Namespace)
	if err != nil {
		u.logger.Error("Failed to get secret from repository", zap.Error(err))
		return err
	}

	gitCloneAuth := &http.BasicAuth{
		Username: string(secret.Data["username"]),
		Password: string(secret.Data["password"]),
	}

	// if string(secret.Data["type"]) == "ssh" {
	// 	sshPublicKey, err := ssh.ParsePublicKey(secret.Data["publicKey"])

	// }

	// TODO: SSH AUTHENTICATION
	gitRepository, err := git.PlainClone(repoDir, false, &git.CloneOptions{
		Auth: gitCloneAuth,
	})
	if err != nil && !errors.Is(err, git.ErrRepositoryAlreadyExists) {
		u.logger.Error("Failed to plain clone repository", zap.Error(err))
		return err
	}

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		gitRepository, err = git.PlainOpen(repoDir)
		if err != nil {
			u.logger.Error("Failed to open repository", zap.Error(err))
			return err
		}
	}

	worktree, err := gitRepository.Worktree()
	if err != nil {
		u.logger.Error("Failed to get worktree", zap.Error(err))
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName(currRepository.Branch)})
	if err != nil {
		u.logger.Error("Failed to checkout branch", zap.Error(err))
		return err
	}

	err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		u.logger.Error("Failed to pull branch", zap.Error(err))
		return err
	}

	return nil
}

// List implements UseCase.
func (u useCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Repository], error) {
	l, err := u.repository.List(ctx, namespace, chunkPagination)
	if err != nil {
		return pagination.ChunkingPaginationResponse[Repository]{}, err
	}

	repositorys := []Repository{}
	for _, i := range l.Items {
		repositorys = append(repositorys, Repository{
			Name:           i.GetName(),
			Namespace:      i.GetNamespace(),
			RepositorySpec: i.Spec,
		})
	}

	return pagination.ChunkingPaginationResponse[Repository]{
		Items: repositorys,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, r Repository) (Repository, error) {
	_, err := u.secretUseCase.Apply(ctx, secret.Secret{
		Name:      r.Name,
		Namespace: r.Namespace,
		Data: map[string][]byte{
			"type":       []byte(r.Auth.Type),
			"username":   []byte(r.Auth.Username),
			"password":   []byte(r.Auth.Password),
			"publickKey": []byte(r.Auth.PublicKey),
		},
	})
	if err != nil {
		return Repository{}, err
	}

	r.RepositorySpec.AuthRef = commonv1alpha1.Ref{
		Name:      r.Name,
		Namespace: r.Namespace,
	}
	newRepository := commonv1alpha1.Repository{
		Spec: r.RepositorySpec,
	}

	newRepository.SetName(r.Name)
	newRepository.SetNamespace(r.Namespace)

	s, err := u.repository.Apply(ctx, newRepository)
	if err != nil {
		return Repository{}, err
	}

	return Repository{
		Name:           s.GetName(),
		Namespace:      s.GetNamespace(),
		RepositorySpec: s.Spec,
	}, nil
}
