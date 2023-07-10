package controller

import (
	"ticket-pimp/internal/domain"
	"ticket-pimp/internal/extapi"
)

type Git struct {
	api extapi.IGit
}

func NewGitController(url, token string) *Git {
	return &Git{
		extapi.NewGitClient(url, token),
	}
}

type RepoCreator interface {
	CreateRepo(name string) (*domain.Git, error)
}

func (g *Git) CreateRepo(name string) (*domain.Git, error) {
	//Create git repository with iGit interface;
	repo, err := g.api.NewRepo(name)
	if err != nil {
		return nil, err
	}

	//Set 'apps' as collaborator to created repository;
	_, err = g.api.AppsAsCollaboratorTo(repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
