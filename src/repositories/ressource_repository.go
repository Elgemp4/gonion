package repositories

import "gonion/src/domain"


type RessourceRepository interface {
	LoadRessource(ressource domain.Ressource) (string, error)
}