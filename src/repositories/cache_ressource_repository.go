package repositories

import (
	"gonion/src/domain"
	"log/slog"
)

type CacheRessourceRepository struct{
	FsRepo 		*FsRessourceRepository
	NpmRepo 	*NpmRessourceRepository 
}

func NewCacheRessourceRepository(fsRepo *FsRessourceRepository, npmRepo *NpmRessourceRepository) *CacheRessourceRepository{
	return &CacheRessourceRepository{
		FsRepo: fsRepo,
		NpmRepo: npmRepo,
	}
}

func (c *CacheRessourceRepository) LoadRessource(ressource domain.Ressource) (string, error) {
	if(c.FsRepo.IsCached(ressource)){
		slog.Info("GET (from cache)", "ressource", ressource.RessourceName())
		return c.FsRepo.GetRessourcePath(ressource), nil;
	}

	writer, err := c.FsRepo.GetRessourceWriter(ressource)

	if(err != nil){
		return "", err
	}

	netError :=c.NpmRepo.DownloadRessource(ressource, writer)

	if(netError != nil){
		return "", netError
	}

	switch ressource.(type) {
		case *domain.NpmMeta:
			err = c.FsRepo.RewriteUrls(ressource)
			if err != nil{
				return "", err
			}
	}
	slog.Info("GET (from fetch)", "ressource", ressource.RessourceName())

	return c.FsRepo.GetRessourcePath(ressource), nil;
}
