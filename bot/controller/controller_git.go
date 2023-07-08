package controller

import d "ticket-pimp/bot/domain"

func (wc *WorkflowController) CreateRepo(name string) (*d.Git, error) {
	//Create git repository with iGit interface;
	repo, err := wc.iGit.NewRepo(name)
	if err != nil {
		return nil, err
	}

	//Set 'apps' as collaborator to created repository;
	_, err = wc.iGit.AppsAsCollaboratorTo(repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
