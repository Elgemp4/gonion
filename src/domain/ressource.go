package domain

type Ressource interface {
	RelativeCachedPath() string
	RessourceName() string
}