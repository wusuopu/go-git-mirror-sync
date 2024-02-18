package services

import (
	"app/config"
	"app/di"
	"app/models"
	"app/utils"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
		return err
	}

	di.Container.Logger.Info(fmt.Sprintf("Finish Clone %s", r.Url))
	// 仅保留 .git 目录，其他文件都删除，节约磁盘空间
	files, err := os.ReadDir(target)
	if err == nil {
		for _, entry := range files {
			if entry.Name() == ".git" {
				continue
			}
			os.RemoveAll(path.Join(target, entry.Name()))
		}
	}
	return nil
}

// 获取对应 hash 的 commit 信息
func (s *RepositoryService) getCommintLogForHash(repo *git.Repository,target string, refName string, hash string) (*object.Commit, error) {
	commit, err := repo.CommitObject(plumbing.NewHash(hash))
	if err == nil {
		return commit, nil
	}

	// 有时 tag 的 hash 指向的 commit 不正确
	cmd := exec.Command("git", "log", "-n", "1", "--pretty=format:%H%n%ct%n%B", refName)
	cmd.Dir = target
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	data := strings.SplitN(string(out), "\n", 3)
	hash = data[0]		// 正确的 hash 值
	sec, _ := strconv.ParseInt(data[1], 10, 64)
	date := time.Unix(sec, 0)
	message := data[2]
	commit = &object.Commit{
		Hash: plumbing.NewHash(hash),
		Message: message,
		Author: object.Signature{
			When: date,
		},
	}
	return commit, nil
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
// 获取仓库的所有分支名或Tag名
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
// 列出远程仓库的所有 tag、分支信息
func (s *RepositoryService) listRemoteAllRefs(repo *git.Repository, remoteName string, auth transport.AuthMethod) (*[]map[string]string, error) {
	var results []map[string]string
	remote, err := repo.Remote(remoteName)
	if err != nil {
		return &results, err
	}
	refs, err := remote.List(&git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		return &results, err
	}
	for _, ref := range refs {
		item := make(map[string]string)
		item["ref"] = ref.Name().String()
		if ref.Name().IsTag() {
			item["tag"] = "true"
			item["name"] = strings.TrimPrefix(item["ref"], "refs/tags/")
		} else {
			item["tag"] = ""
			item["name"] = strings.TrimPrefix(item["ref"], "refs/heads/")
		}
		results = append(results, item)
	}

	return &results, nil
}
// 删除远程仓库中所有在本地不存在的分支及 tag
func (s *RepositoryService) pruneRemoteRefs(repo *git.Repository, remoteName string, auth transport.AuthMethod, target string) error {
	allBranches, err := s.getAllBranchNames(target, false, true)
	if err != nil {
		return err
	}
	allTags, err := s.getAllBranchNames(target, true, false)
	if err != nil {
		return err
	}
	allRemoteRefs, err := s.listRemoteAllRefs(repo, remoteName, auth) 
	if err != nil {
		return err
	}
	remote, err := repo.Remote(remoteName)
	if err != nil {
		return err
	}
	cfg := remote.Config()


	di.Container.Logger.Info(fmt.Sprintf("Start to Prune %s", cfg.URLs[0]))

	var willRemoveSpecs []gitconfig.RefSpec
	for _, item := range *allRemoteRefs {
		if item["tag"] == "" {
			if (*allBranches)[item["name"]] == "" {
				// 该分支在本地不存在
				di.Container.Logger.Debug(fmt.Sprintf("Remove remote %s branch %s", remoteName, item["name"]))
				willRemoveSpecs = append(willRemoveSpecs, gitconfig.RefSpec(":" + item["ref"]))
			}
		} else {
			if (*allTags)[item["name"]] == "" {
				// 该 tag 在本地不存在
				di.Container.Logger.Debug(fmt.Sprintf("Remove remote %s tag %s", remoteName, item["name"]))
				willRemoveSpecs = append(willRemoveSpecs, gitconfig.RefSpec(":" + item["ref"]))
			}
		}
	}

	if len(willRemoveSpecs) > 0 {
		err = repo.Push(&git.PushOptions{
			RemoteName: cfg.Name,
			RemoteURL: cfg.URLs[0],
			Auth: auth,
			Force: true,
			RefSpecs: willRemoveSpecs,
		})
		if err != nil && errors.Is(err, git.NoErrAlreadyUpToDate) {
			// 忽略 already up-to-date 的错误
			err = nil
		}
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
		commit, err := s.getCommintLogForHash(repo, target, name, hash)
		if err != nil {
			di.Container.Logger.Error(fmt.Sprintf("%s get tag log for %s(%s) error %s", r.Name, name, hash, err.Error()))
		} else {
			obj.Hash = commit.Hash.String()
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
func (s *RepositoryService) SyncMirror(r models.Repository, m models.Mirror) error {
	if r.Name == "" {
		return fmt.Errorf("repository name is blank")
	}
	target := path.Join(RepositoryBasePath, r.Name)
	repo, err := git.PlainOpen(target)
	if err != nil {
		return err
	}
	remote, err := repo.Remote(m.Name)
	if err != nil {
		return err
	}
	cfg := remote.Config()
	auth, err := s.MakeGitAuth(models.Repository{
		AuthType: m.AuthType,
		Username: m.Username,
		Password: m.Password,
		SSHKey: m.SSHKey,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	di.Container.Logger.Info(fmt.Sprintf("Start to Push %s", cfg.URLs[0]))
	// 强制推送所有的分支、tag
	err = repo.Push(&git.PushOptions{
		RemoteName: cfg.Name,
		RemoteURL: cfg.URLs[0],
		Auth: auth,
		Force: true,
		RefSpecs: []gitconfig.RefSpec{gitconfig.RefSpec("refs/tags/*:refs/tags/*"), gitconfig.RefSpec("refs/heads/*:refs/heads/*")},
	})
	if err != nil && errors.Is(err, git.NoErrAlreadyUpToDate) {
		// 忽略 already up-to-date 的错误
		err = nil
	}
	if err != nil {
		di.Container.Logger.Error(fmt.Sprintf("Push Error %s %s", cfg.URLs[0], err.Error()))
		return err
	}
	di.Container.Logger.Info(fmt.Sprintf("Finish Push %s", cfg.URLs[0]))

	err = s.pruneRemoteRefs(repo, m.Name, auth, target)
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
		// remote url不一致，先删除再重新创建
		s.DeleteRemote(r, m)
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

	os.RemoveAll(path.Join(target, ".git/refs/remotes", m.Name))
	return repo.DeleteRemote(m.Name)
}