package cmd

import (
	"fmt"
	"log"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cliConfig "github.com/yottab/cli/config"
	ybApi "github.com/yottab/proto-api/proto"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	ybRemoteName    = "io.YOTTAb.git"
	gitYbPath       = "https://git.yottab.io/%s/%s.git"
	pushLogIDFormat = "%s:%s" // imageName:imageTa
)

func pushRepository(cmd *cobra.Command, args []string) {
	baseFolder := cmd.Flag("path").Value.String() // Repository path
	imageTag := cmd.Flag("tag").Value.String()    // Repository tag
	imageName := cmd.Flag("name").Value.String()  // Repository name

	//Set default image name if needed
	if imageName == "" {
		imageName = defaultRepositoryName(baseFolder)
	}

	checkExistDockerfile(baseFolder)

	//Let's select and inform build tag
	user := viper.GetString(cliConfig.KEY_USER)
	repo := getYbRepo(baseFolder, imageName, user)
	latestTag := latestTag(repo)
	if repositoryYbIsClean(repo) == false {
		log.Fatal("Repository is dirty! please commit your changes and try again!")
	}
	commitHash := latestCommit(repo)
	//Select tag based on user choice
	imageTag = selectTag(imageTag, flagGitCommitHash, flagGitTag, latestTag, commitHash)
	confirmF(cmd, "build & push docker image \033[1m%s:%s\033[0m into yottab hub", imageName, imageTag)

	//Push and request build
	push(repo, user, imageName)
	client := grpcConnect()
	defer client.Close()
	req := new(ybApi.ImgBuildReq)
	req.RepositoryName = imageName
	req.RepositoryTag = imageTag

	_, err := client.V2().ImgBuild(client.Context(), req)
	uiCheckErr("Could not Build the Repository", err)
	log.Print("Build started!\r\nWaiting for builder log to get ready...")
	time.Sleep(20 * time.Second)
	streamBuildLog(imageName, imageTag, true)
	log.Printf("Enter this command to see more:\r\n$: yb push log --name=%s --tag=%s\r\n", imageName, imageTag)
}

func pushLog(cmd *cobra.Command, args []string) {
	imageTag := cmd.Flag("tag").Value.String()   // Repository tag
	imageName := cmd.Flag("name").Value.String() // Repository name

	streamBuildLog(imageName, imageTag, false)
}

func selectTag(imageTag string, flagGitCommitHash, flagGitTag bool, latestTag, commitHash string) string {
	if flagGitCommitHash && flagGitTag {
		log.Fatal("Please set either one of (tag-git, tag-commit) boolean switches!")
	}
	if flagGitTag {
		if len(latestTag) == 0 {
			log.Fatal("latest tag can't be found!")
		}
		return latestTag
	}
	if flagGitCommitHash {
		if len(commitHash) == 0 {
			log.Fatal("no commit hash detected!")
		}
		return commitHash
	}
	//NOOP short circuit
	return imageTag
}
func checkExistDockerfile(basePath string) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatalf("ioutil.ReadDir Err: %v", err)
	}

	for _, file := range files {
		if file.Name() == "Dockerfile" {
			return
		}
	}
	log.Fatalf("Err: can find Dockerfile at [%s]", basePath)
}

func commit(repo *git.Repository) (commitHash string) {
	if repositoryYbIsClean(repo) == false {
		commitHash = repositoryYbCommit(repo)
	}
	return
}
func push(repo *git.Repository, user, appName string) {
	token := viper.GetString(cliConfig.KEY_TOKEN)

	fmt.Println("Start PUSH")
	remo := addYbRemote(repo, appName, user)
	err := ybPush(remo, user, token)
	cleaRemote(repo)
	uiCheckErr("push.ybPush", err)
	return
}

func repositoryYbIsClean(repo *git.Repository) bool {
	wt, err := repo.Worktree()
	uiCheckErr("repositoryYbIsClean.Worktree", err)

	s, err := wt.Status()
	uiCheckErr("repositoryYbIsClean.Status", err)

	return s.IsClean()
}
func latestTag(repo *git.Repository) (tag string) {
	tagRefs, err := repo.Tags()
	uiCheckErr("latestTag", err)

	var latestTagCommit *object.Commit
	var latestTagName string
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := repo.ResolveRevision(revision)
		uiCheckErr("latestTag", err)

		commit, err := repo.CommitObject(*tagCommitHash)
		uiCheckErr("latestTag", err)

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().Short()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().Short()
		}

		return nil
	})
	uiCheckErr("latestTag", err)

	return latestTagName
}

func repositoryYbCommit(repo *git.Repository) string {
	w, err := repo.Worktree()
	_, err = w.Add(".")
	uiCheckErr("repositoryYbCommit.Add", err)

	_, err = w.Commit(fmt.Sprintf("yb @ %s", time.Now().Format(time.RFC822)), &git.CommitOptions{
		Author: &object.Signature{
			When: time.Now(),
		},
	})
	uiCheckErr("repositoryYbCommit.Commit", err)
	headRef, err := repo.Head()
	uiCheckErr("repositoryYbCommit.Head", err)
	return headRef.String()[:10]
}

func latestCommit(repo *git.Repository) string {
	headRef, err := repo.Head()
	uiCheckErr("repositoryYbCommit.Head", err)
	return headRef.String()[:10]
}

// getRepo open Repository and
// if yottab.Remote not exist, add the Remote to Repository
func getYbRepo(path, appName, user string) (repo *git.Repository) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			fmt.Printf("Repasitory %s Create at %s", appName, path)
			repo, err = git.PlainInit(path, false)
			if err != nil {
				return
			}
		} else {
			return
		}
		fmt.Printf("Repasitory %s open at %s", appName, path)
	}
	return
}

type stdWriter struct {
	io.Writer
}

func (in *stdWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func ybPush(remo *git.Remote, user, pass string) error {
	// progress := sideband.NewDemuxer(sideband.Sideband, os.Stdout)
	err := remo.Push(&git.PushOptions{
		RemoteName: ybRemoteName,
		Auth: &http.BasicAuth{
			Username: user,
			Password: pass,
		},
		Progress: os.Stdout,
	})
	return err
}

// addRemote Add a new remote, with the default fetch refspec
func addYbRemote(repo *git.Repository, appName, user string) *git.Remote {
	// check exist YbRemote
	// Temproray add remote
	_, err := repo.Remote(ybRemoteName)
	if err == nil {
		err = repo.DeleteRemote(ybRemoteName)
		uiCheckErr("Delete old remote", err)
	}
	remo, err := repo.CreateRemote(&config.RemoteConfig{
		Name: ybRemoteName,
		URLs: []string{
			getRemoteURL(user, appName)},
	})
	uiCheckErr("Create remote", err)
	return remo
}

func cleaRemote(repo *git.Repository) error {
	err := repo.DeleteRemote(ybRemoteName)
	uiCheckErr("Delete old remote", err)
	uiCheckErr("Create remote", err)
	return err
}

func getRemoteURL(user, app string) string {
	return fmt.Sprintf(gitYbPath, user, app)
}
