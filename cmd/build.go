package cmd

import (
	"fmt"
	"io/ioutil"
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

	confirmF(cmd, "build & push docker image \033[1m%s:%s\033[0m into yottab hub", imageName, imageTag)
	latestTag, commitHash := push(imageName, imageTag, baseFolder)

	//Select tag based on user choice
	imageTag = selectTag(imageTag, flagGitCommitHash, flagGitTag, latestTag, commitHash)

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
		panic("Please set either one of (tag-git, tag-commit) boolean switches!")
	}
	if flagGitTag {
		return latestTag
	}
	if flagGitCommitHash {
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

func push(imageName, imageTag, repoPath string) (commitHash, latestTag string) {
	user := viper.GetString(cliConfig.KEY_USER)
	token := viper.GetString(cliConfig.KEY_TOKEN)

	repo, err := getYbRepo(repoPath, imageName, user)
	if err != nil {
		log.Fatal("getYbRepo: ", err)
	}
	latestTag = repositoryYbLatestTag(repo)
	if repositoryYbIsClean(repo) == false {
		commitHash = repositoryYbCommit(repo, imageTag)
	}

	fmt.Println("Start PUSH")
	err = ybPush(repo, user, token)
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
func repositoryYbLatestTag(repo *git.Repository) (tag string) {
	tagrefs, err := repo.Tags()
	uiCheckErr("repositoryYbLatestTag.Tags", err)

	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		fmt.Println(t)
		tag = t.Name().String()
		return nil
	})
	uiCheckErr("repositoryYbLatestTag.Tags", err)
	return
}

func repositoryYbCommit(repo *git.Repository, tag string) string {
	w, err := repo.Worktree()
	_, err = w.Add(".")
	uiCheckErr("repositoryYbCommit.Add", err)

	hash, err := w.Commit(tag, &git.CommitOptions{
		Author: &object.Signature{
			When: time.Now(),
		},
	})
	uiCheckErr("repositoryYbCommit.Commit", err)
	return hash.String()[:10]
}

// getRepo open Repository and
// if yottab.Remote not exist, add the Remote to Repository
func getYbRepo(path, appName, user string) (repo *git.Repository, err error) {
	repo, err = git.PlainOpen(path)
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

	err = addYbRemote(repo, appName, user)
	return
}

func ybPush(repo *git.Repository, user, pass string) error {
	remo, err := repo.Remote(ybRemoteName)
	if err != nil {
		log.Println("ybPush.Remote: ", err)
		return err
	}

	err = remo.Push(&git.PushOptions{
		RemoteName: ybRemoteName,
		Auth: &http.BasicAuth{
			Username: user,
			Password: pass,
		},
	})

	return err
}

// addRemote Add a new remote, with the default fetch refspec
func addYbRemote(repo *git.Repository, appName, user string) error {
	// check exist YbRemote
	_, err := repo.Remote(ybRemoteName)
	if err == git.ErrRemoteNotFound {
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: ybRemoteName,
			URLs: []string{
				getRemoteURL(user, appName)},
		})
	}

	return err
}

func getRemoteURL(user, app string) string {
	return fmt.Sprintf(gitYbPath, user, app)
}
