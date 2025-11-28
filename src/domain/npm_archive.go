package domain

import (
	"fmt"
)

type NpmArchive struct {
	Name string
}

func NewNPMArchive(name string) *NpmArchive {
	return &NpmArchive{
		Name: name,
	}
}


func (a *NpmArchive) RelativeCachedPath() string{
	return fmt.Sprintf("/packages/%s",  a.Name)
}

func (a *NpmArchive) RessourceName() string{
	return a.Name
}