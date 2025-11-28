package services

import (
	"gonion/src/domain"
	"gonion/src/repositories"
	"net/url"
	"path/filepath"
	"strings"
)

type NpmCacheService struct{
	RessourceRepository	repositories.RessourceRepository
}

func NewNpmCacheService(repository repositories.RessourceRepository) *NpmCacheService {
	return &NpmCacheService{
		RessourceRepository: repository,
	}
}

func (ncs* NpmCacheService) sanytizeRessource(path string) (string, error){
	decodedUrl, decodingErr := url.QueryUnescape(path)
	
	cleanedUrl := filepath.Clean(decodedUrl)

	return cleanedUrl, decodingErr
}

func (ncs* NpmCacheService) resolveRessource(ressourceName string) domain.Ressource {
	var ressource domain.Ressource

	if(strings.HasSuffix(ressourceName, ".tgz")){
		println("Package asked")
		ressource = domain.NewNPMArchive(ressourceName)
	}else{
		println("Meta asked")
		ressource = domain.NewNpmMeta(ressourceName)
	}

	return ressource
}

func (ncs *NpmCacheService) GetRessourceLocalPath(name string) (string, error) {
	ressourceName, sanytiszeError := ncs.sanytizeRessource(name)

	if(sanytiszeError != nil) {
		return "", domain.ErrBadPackageName
	}

	ressource := ncs.resolveRessource(ressourceName)

	return ncs.RessourceRepository.LoadRessource(ressource)
}