package repositories

import (
	"errors"
	"gonion/src/domain"
	"log/slog"
	"sync"
)

type CacheRessourceRepository struct{
	FsRepo 		*FsRessourceRepository
	NpmRepo 	*NpmRessourceRepository 
	Locks		sync.Map
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

	mu := &sync.Mutex{}
	actual, _ := c.Locks.LoadOrStore(ressource.RessourceName(), mu)
	lock := actual.(*sync.Mutex)

	lock.Lock()
	defer lock.Unlock()

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
		if(errors.Is(netError, domain.ErrRessourceCachingFailure)){
			c.FsRepo.CleanRessource(ressource)
		}

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
	renameErr := c.FsRepo.ValidateCaching(ressource)
	if(renameErr != nil){
		return "", renameErr
	}
	return c.FsRepo.GetRessourcePath(ressource), nil;
}
