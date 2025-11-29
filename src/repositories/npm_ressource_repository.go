package repositories

import (
	"gonion/src/domain"
	"io"
	"log/slog"

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
		return domain.ErrUnreachableUpstream
	}

	if resp.StatusCode() != 200 {
        output.Close()
        slog.Error("404 error from upstream", "ressource", ressource.RessourceName())
		return domain.ErrNotFound
    }

	_, errCopy := io.Copy(output, resp.RawBody())

	errOutputClose := output.Close()
	errInputClose := resp.RawBody().Close()

	if(errCopy != nil || errOutputClose != nil || errInputClose != nil){
		return domain.ErrRessourceCachingFailure
	}

	return nil
}
