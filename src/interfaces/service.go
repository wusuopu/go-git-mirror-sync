package interfaces

import (
	"app/models"

	"github.com/go-git/go-git/v5/plumbing/transport"
)


type IRepositoryService interface {
	Add(a, b int) int
	Minus(a, b int) int
	GetCoreStorePath(r models.Repository) string
	MakeGitAuth(r models.Repository) (transport.AuthMethod, error) 
	Clone(r models.Repository) error
	SyncOrigin(r models.Repository) error
	SyncMirror(r models.Repository) error
	CreateRemote(r models.Repository, m models.Mirror) error
	DeleteRemote(r models.Repository, m models.Mirror) error
	BuildBranchInfo(r models.Repository) error
}