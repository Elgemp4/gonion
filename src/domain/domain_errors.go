package domain

import "errors"

var(
	ErrBadRessourceName = errors.New("the provided ressource name is malformed")
	ErrNotFound = errors.New("the ressource could not be found upstream")
	ErrUnreachableUpstream = errors.New("the upstream server could not be contacted")
	ErrRessourceCachingFailure = errors.New("the ressource could not be stored properly")
)