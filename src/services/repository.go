package services

import (
	"app/config"
	"app/di"
	"app/models"
	"app/utils"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	ssh2 "golang.org/x/crypto/ssh"
)

const RepositoryBasePath = "data/repositories"

type RepositoryService struct {
}

func (s *RepositoryService) Add(a, b int) int {
	return a + b
}

func (s *RepositoryService) Minus(a, b int) int {
	return a - b
}

func (s *RepositoryService) GetCoreStorePath(r models.Repository) string {
	target := path.Join(RepositoryBasePath, r.Name, ".git")
	return target
}
func (s *RepositoryService) MakeGitAuth(r models.Repository) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	if r.AuthType == "password" {
		auth = &http.BasicAuth{
			Username: *r.Username,	
			Password: *r.Password,
		}
	} else {
		var user string
		var password = *r.Password
		if *r.Username == "" {
			user = "git"
		} else {
			user = *r.Username
		}
		key, err := ssh.NewPublicKeys(user, []byte(*r.SSHKey), password)
		if err != nil {
			return nil, err
		}
		auth = key
		// 忽略 SSH_KNOWN_HOSTS 检查
		auth.(*ssh.PublicKeys).HostKeyCallbackHelper = ssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh2.InsecureIgnoreHostKey(),
		}
	}
	return auth, nil
}
// 克隆新的仓库
func (s *RepositoryService) Clone(r models.Repository) error {
	if r.Name == "" {
		return fmt.Errorf("repository name is blank")
	}
	auth, err := s.MakeGitAuth(r)
	if err != nil {
		return err
	}
	
	di.Container.Logger.Info(fmt.Sprintf("Start to Clone %s", r.Url))
	target := path.Join(RepositoryBasePath, r.Name)
	// 先清空目录
	if utils.FsIsExist(target) {
		err := os.RemoveAll(target)
		if err != nil {
			return err
		}
	}

	_, err = git.PlainClone(target, false, &git.CloneOptions{
		URL: r.Url,
		RemoteName: "origin",
		Auth: auth,
		Tags: git.AllTags,
	})
	if err != nil {
		di.Container.Logger.Error(fmt.Sprintf("Clone Error %s %s", r.Url, err.Error()))
	} else {
		di.Container.Logger.Info(fmt.Sprintf("Finish Clone %s", r.Url))
	}
	
	return err
}

// 清除本地的所有分支
func (s *RepositoryService) cleanAllLocalBranches(repo *git.Repository, target string) error {
	cfg, err := repo.Config()
	if err != nil {
		return err
	}

	for idx, _ := range cfg.Branches {
		err := repo.DeleteBranch(idx)
		if err != nil {
			return err
		}
	}
	dir := path.Join(target, ".git/refs/heads")
	if utils.FsIsExist(dir) {
		return os.RemoveAll(dir)
	}
	return nil
}
// 获取远程仓库的所有分支名或Tag名
func (s *RepositoryService) getAllBranchNames(target string, isTag bool, isRemote bool) (*map[string]string, error) {
	// {<name>: <hash>}
	var all_branches = make(map[string]string)
	var dir string
	if isTag {
		dir = path.Join(target, ".git/refs/tags")
	} else if isRemote {
		dir = path.Join(target, ".git/refs/remotes/origin")
	} else {
		dir = path.Join(target, ".git/refs/heads")
	}
	hashReg := regexp.MustCompile(`^[0-9a-z]+\s$`)
	e := filepath.Walk(dir, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			branchName := p[len(dir)+1:]
			content, _ := os.ReadFile(p)
			hash := strings.TrimSpace(string(content))
			if hashReg.Match(content) {
				all_branches[branchName] = hash
			}
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return &all_branches, nil
}
// 根据远程仓库创建同名的分支
func (s *RepositoryService) createBranchesFromRemote(repo *git.Repository, target string, allBranches *map[string]string) error {
	for name, hash := range *allBranches {
		branch := &gitconfig.Branch{
			Name: name,
			Remote: "origin",
			Merge: plumbing.NewBranchReferenceName(name),
		}
		repo.CreateBranch(branch)
		filename := path.Join(target, ".git", branch.Merge.String())
		dir := filepath.Dir(filename)
		if !utils.FsIsExist(dir) {
			os.MkdirAll(dir, os.ModePerm)
		}
		os.WriteFile(filename, []byte(hash + "\n"), os.ModePerm)
	}
	return nil
}
/*
同步拉取源仓库的记录
*/
func (s *RepositoryService) SyncOrigin(r models.Repository) error {
	if r.Name == "" {
		return fmt.Errorf("repository name is blank")
	}
	target := path.Join(RepositoryBasePath, r.Name)
	repo, err := git.PlainOpen(target)
	if err != nil {
		return err
	}

	auth, err := s.MakeGitAuth(r)
	if err != nil {
		return err
	}
	di.Container.Logger.Info(fmt.Sprintf("Start to Fetch %s", r.Url))

	// 先删除仓库在本地的所有tag、分支
	os.RemoveAll(path.Join(target, ".git/refs/tags"))
	os.RemoveAll(path.Join(target, ".git/refs/remotes/origin"))

	// 拉取所有分支、tag
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth: auth,
		Tags: git.AllTags,
		Force: true,
		InsecureSkipTLS: config.Config["GIT_INSECURE_SKIP_TLS"].(bool),
	})
	if err != nil {
		di.Container.Logger.Error(fmt.Sprintf("Fetch Error %s %s", r.Url, err.Error()))
		return err
	}

	// 远程仓库的所有分支
	allBranches, err := s.getAllBranchNames(target, false, true)
	if err != nil {
		return err
	}
	if len(*allBranches) == 0 {
		// 清的仓库
		return nil
	}
	// 将所有本地分支与远程分支同步
	err = s.cleanAllLocalBranches(repo, target)
	if err != nil {
		return err
	}
	err = s.createBranchesFromRemote(repo, target, allBranches)

	di.Container.Logger.Info(fmt.Sprintf("Finish Fetch %s", r.Url))

	return err
}
/*
获取所有分支的信息并保存到数据库中
*/
func (s *RepositoryService) BuildBranchInfo(r models.Repository) error {
	if r.Name == "" {
		return fmt.Errorf("repository name is blank")
	}
	target := path.Join(RepositoryBasePath, r.Name)
	repo, err := git.PlainOpen(target)
	if err != nil {
		return err
	}

	di.Container.Logger.Info(fmt.Sprintf("Start to BuildBranchInfo for %s", r.Name))
	// 先清空所有的分支记录
	di.Container.DB.Unscoped().Where("repository_id = ?", r.ID).Delete(&models.Branch{})

	allBranches, err := s.getAllBranchNames(target, false, false)
	di.Container.Logger.Debug(fmt.Sprintf("%s has %d branches", r.Name, len(*allBranches)))
	if err != nil {
		return err
	}
	for name, hash := range *allBranches {
		obj := models.Branch{
			Name: name,
			Hash: hash,
			IsTag: false,
			RepositoryId: r.ID,
		}
		commit, err := repo.CommitObject(plumbing.NewHash(hash))
		if err != nil {
			di.Container.Logger.Error(fmt.Sprintf("%s get branch log for %s(%s) error %s", r.Name, name, hash, err.Error()))
		} else {
			obj.CommittedAt = &commit.Author.When
			obj.CommitMsg = &commit.Message
		}
		di.Container.Logger.Debug(fmt.Sprintf("Create Branch for %s %s", r.Name, name))
		di.Container.DB.Create(&obj)
	}
	allTags, err := s.getAllBranchNames(target, true, false)
	di.Container.Logger.Debug(fmt.Sprintf("%s has %d tags", r.Name, len(*allTags)))
	if err != nil {
		return err
	}
	for name, hash := range *allTags {
		obj := models.Branch{
			Name: name,
			Hash: hash,
			IsTag: true,
			RepositoryId: r.ID,
		}
		commit, err := repo.CommitObject(plumbing.NewHash(hash))
		if err != nil {
			di.Container.Logger.Error(fmt.Sprintf("%s get tag log for %s(%s) error %s", r.Name, name, hash, err.Error()))
		} else {
			obj.CommittedAt = &commit.Author.When
			obj.CommitMsg = &commit.Message
		}
		di.Container.Logger.Debug(fmt.Sprintf("Create Tag for %s %s", r.Name, name))
		di.Container.DB.Create(&obj)
	}
	return nil
}
/*
同步上传镜像仓库的记录
*/
func (s *RepositoryService) SyncMirror(r models.Repository) error {
	if r.Name == "" {
		return fmt.Errorf("repository name is blank")
	}
	target := path.Join(RepositoryBasePath, r.Name)
	repo, err := git.PlainOpen(target)
	if err != nil {
		return err
	}
	// s.CreateRemote(r)

	// TODO
	err = repo.Push(&git.PushOptions{

	})
	return err
}
func (s *RepositoryService) CreateRemote(r models.Repository, m models.Mirror) error {
	if r.Name == "" {
		return fmt.Errorf("repository name is blank")
	}
	target := path.Join(RepositoryBasePath, r.Name)
	repo, err := git.PlainOpen(target)
	if err != nil {
		return err
	}

	remote, err := repo.Remote(m.Name)
	if remote != nil {
		cfg := remote.Config()
		// 该仓库已经存在
		if cfg.URLs[0] == m.Url {
			return nil
		}
	}
	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: m.Name,
		URLs: []string{m.Url},
		Mirror: false,
	})
	return err
}
func (s *RepositoryService) DeleteRemote(r models.Repository, m models.Mirror) error {
	if r.Name == "" {
		return fmt.Errorf("repository name is blank")
	}
	target := path.Join(RepositoryBasePath, r.Name)
	repo, err := git.PlainOpen(target)
	if err != nil {
		return err
	}

	remote, err := repo.Remote(m.Name)
	if remote == nil {
		return nil
	}

	return repo.DeleteRemote(m.Name)
}