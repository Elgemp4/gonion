package repositories

import (
	"bytes"
	"gonion/src/domain"
	"io"
	"os"
	"path/filepath"
)

type FsRessourceRepository struct {
	FsRoot		string
	NpmUrl		string
	LocalUrl 	string
}

func NewFsRessourceRepository(fsRoot string, npmUrl string, localUrl string) *FsRessourceRepository {
	return &FsRessourceRepository{
		FsRoot: fsRoot,
		NpmUrl: npmUrl,
		LocalUrl: localUrl,
	}
}

func (fsRepo *FsRessourceRepository) IsCached(ressource domain.Ressource) bool{
	_, err := os.Stat(fsRepo.GetRessourcePath(ressource)) 
	return err == nil
}

func (fsRepo *FsRessourceRepository) GetRessourcePath(ressource domain.Ressource) string{
	return filepath.Join(fsRepo.FsRoot, ressource.RelativeCachedPath())
}

func (fsRepo *FsRessourceRepository) GetRessourceWriter(ressource domain.Ressource) (io.WriteCloser, error) {
	path := fsRepo.GetRessourcePath(ressource)

	err := os.MkdirAll(filepath.Dir(path), 0755);

	if(err != nil){
		return nil, err
	}

	return os.Create(path)
}

func (fsRepo *FsRessourceRepository) RewriteUrls(ressource domain.Ressource) error{
	path := fsRepo.GetRessourcePath(ressource)

	data, readErr := os.ReadFile(path)

	if(readErr != nil){
		return readErr
	}

	newData := bytes.ReplaceAll(data, []byte(fsRepo.NpmUrl), []byte(fsRepo.LocalUrl))

	writeErr := os.WriteFile(path, newData, 0655)

	if(writeErr != nil){
		return writeErr
	}

	return nil
}
