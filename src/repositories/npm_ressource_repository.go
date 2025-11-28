package repositories

import (
	"gonion/src/domain"
	"io"

	"github.com/go-resty/resty/v2"
)

type NpmRessourceRepository struct{
	Client 		*resty.Client
}

func NewNpmRessourceRepository(npm_url string) *NpmRessourceRepository{
	client := resty.New()
	client.BaseURL = npm_url
	return &NpmRessourceRepository{
		Client: client,
	}
}

func (npr *NpmRessourceRepository) DownloadRessource(ressource domain.Ressource, output io.WriteCloser) error {
	resp, err := npr.Client.R().
		SetDoNotParseResponse(true).
		Get(ressource.RessourceName())

	if(err != nil){
		return err
	}

	io.Copy(output, resp.RawBody())

	resp.RawBody().Close()

	return nil
}
