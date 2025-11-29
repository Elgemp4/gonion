package repositories

import (
	"bytes"
	"fmt"
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

func  (fsRepo *FsRessourceRepository) fileExist(path string) bool{
	_, err := os.Stat(path) 
	return err == nil
}

func (fsRepo *FsRessourceRepository) IsCached(ressource domain.Ressource) bool{
	return fsRepo.fileExist(fsRepo.GetRessourcePath(ressource))
}

func (fsRepo *FsRessourceRepository) IsBeeingWritten(ressource domain.Ressource) bool{
	return fsRepo.fileExist(fsRepo.GetPartRessourcePath(ressource))
}

func (fsRepo *FsRessourceRepository) GetRessourcePath(ressource domain.Ressource) string{
	return filepath.Join(fsRepo.FsRoot, ressource.RelativeCachedPath())
}

func (fsRepo *FsRessourceRepository) GetPartRessourcePath(ressource domain.Ressource) string{
	return filepath.Join(fsRepo.FsRoot, fmt.Sprintf("%s.part", ressource.RelativeCachedPath()))
}

func (fsRepo *FsRessourceRepository) GetRessourceWriter(ressource domain.Ressource) (io.WriteCloser, error) {
	path := fsRepo.GetPartRessourcePath(ressource)

	err := os.MkdirAll(filepath.Dir(path), 0755);

	if(err != nil){
		return nil, err
	}

	return os.Create(path)
}

func (fsRepo *FsRessourceRepository) ValidateCaching(ressource domain.Ressource) error {
	partPath := fsRepo.GetPartRessourcePath(ressource)
	definitivePath := fsRepo.GetRessourcePath(ressource)

	return os.Rename(partPath, definitivePath)
}

func (fsRepo *FsRessourceRepository) CleanRessource(ressource domain.Ressource) {
	partPath := fsRepo.GetPartRessourcePath(ressource)
	definitivePath := fsRepo.GetRessourcePath(ressource)

	if(fsRepo.IsBeeingWritten(ressource)){
		os.Remove(partPath)
	}

	if(fsRepo.IsCached(ressource)){
		os.Remove(definitivePath)
	}
}

func (fsRepo *FsRessourceRepository) RewriteUrls(ressource domain.Ressource) error{
	path := fsRepo.GetPartRessourcePath(ressource)

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
