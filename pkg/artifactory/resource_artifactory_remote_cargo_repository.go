package artifactory

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var cargoRemoteSchema = mergeSchema(baseRemoteSchema, map[string]*schema.Schema{
	"git_registry_url": {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		Description:  `This is the index url, expected to be a git repository. for remote artifactory use "arturl/git/repokey.git"`,
	},
	"anonymous_access": {
		Type:     schema.TypeBool,
		Optional: true,
		Description: "(On the UI: Anonymous download and search) Cargo client does not send credentials when performing download and search for crates. " +
			"Enable this to allow anonymous access to these resources (only), note that this will override the security anonymous access option.",
	},
})

type CargoRemoteRepo struct {
	RemoteRepositoryBaseParams
	RegistryUrl     string `json:"gitRegistryUrl"`
	AnonymousAccess bool   `json:"cargoAnonymousAccess"`
}

var cargoRemoteRepoReadFun = mkRepoRead(packCargoRemoteRepo, func() interface{} {
	return &CargoRemoteRepo{
		RemoteRepositoryBaseParams: RemoteRepositoryBaseParams{
			Rclass: "remote",
			PackageType: "cargo",
		},
	}
})

func resourceArtifactoryRemoteCargoRepository() *schema.Resource {
	return &schema.Resource{
		Create: mkRepoCreate(unpackCargoRemoteRepo, cargoRemoteRepoReadFun),
		Read:   cargoRemoteRepoReadFun,
		Update: mkRepoUpdate(unpackCargoRemoteRepo, cargoRemoteRepoReadFun),
		Delete: deleteRepo,
		Exists: repoExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: cargoRemoteSchema,
	}
}

func unpackCargoRemoteRepo(s *schema.ResourceData) (interface{}, string, error) {
	d := &ResourceData{s}
	repo := CargoRemoteRepo{
		RemoteRepositoryBaseParams: unpackBaseRemoteRepo(s),
		RegistryUrl:                d.getString("git_registry_url", false),
		AnonymousAccess:            d.getBool("anonymous_access", false),
	}
	repo.PackageType = "cargo"
	return repo, repo.Key, nil
}

func packCargoRemoteRepo(r interface{}, d *schema.ResourceData) error {
	repo := r.(*CargoRemoteRepo)
	setValue := packBaseRemoteRepo(d, repo.RemoteRepositoryBaseParams)
	setValue("git_registry_url", repo.RegistryUrl)
	errors := setValue("anonymous_access", repo.AnonymousAccess)

	if len(errors) > 0 {
		return fmt.Errorf("%q", errors)
	}

	return nil
}
