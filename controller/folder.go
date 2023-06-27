package controller

import d "ticket-pimp/domain"

func (wc *WorkflowController) CreateFolder(name string) (*d.Folder, error) {

	//Create ownCloud folder w/ iCloud interface;
	cloud, err := wc.iCloud.CreateFolder(name)
	if cloud == nil {
		return nil, err
	}

	/* Experimental call:
	wc.iCloud.ShareToExternals(cloud)
	*/

	return cloud, err
}
