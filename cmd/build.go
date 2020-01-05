package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cliConfig "github.com/yottab/cli/config"
	ybApi "github.com/yottab/proto-api/proto"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	ybRemoteName    = "io.YOTTAb.git"
	pushLogIDFormat = "%s:%s" // imageName:imageTa
)

func pushRepository(cmd *cobra.Command, args []string) {
	baseFolder := cmd.Flag("path").Value.String() // Repository path
	imageTag := cmd.Flag("tag").Value.String()    // Repository tag
	imageName := cmd.Flag("name").Value.String()  // Repository name

	checkExistDockerfile(baseFolder)

	push(imageName, imageTag, baseFolder)

	client := grpcConnect()
	defer client.Close()
	req := new(ybApi.ImgBuildReq)
	req.RepositoryName = imageName
	req.RepositoryTag = imageTag

	_, err := client.V2().ImgBuild(client.Context(), req)
	uiCheckErr("Could not Build the Repository", err)
	log.Print("Build started!\r\nWaiting for builder log to get ready...")
	time.Sleep(20 * time.Second)
	getPushLog(imageName, imageTag)
	log.Printf("Enter this command to see more:\r\n$: yb push log --name=%s --tag=%s\r\n", imageName, imageTag)
}

func pushLog(cmd *cobra.Command, args []string) {
	imageTag := cmd.Flag("tag").Value.String()   // Repository tag
	imageName := cmd.Flag("name").Value.String() // Repository name
	getPushLog(imageName, imageTag)
}

func getPushLog(appName, appTag string) {
	id := getRequestIdentity(
		fmt.Sprintf(pushLogIDFormat, appName, appTag))
	client := grpcConnect()
	defer client.Close()
	logClient, err := client.V2().ImgBuildLog(context.Background(), id)
	uiCheckErr(fmt.Sprintf("Could not get build log right now!\nTry again in a few soconds using:\n$yb push log --name=%s --tag=%s", appName, appTag), err)
	uiImageLog(logClient)

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

func push(imageName, imageTag, repoPath string) {
	user := viper.GetString(cliConfig.KEY_USER)
	token := viper.GetString(cliConfig.KEY_TOKEN)

	repo, err := getYbRepo(repoPath, imageName, user)
	if err != nil {
		log.Fatal("getYbRepo: ", err)
	}

	if repositoryYbIsClean(repo) == false {
		repositoryYbCommit(repo, imageTag)
	}

	fmt.Println("Start PUSH")
	err = ybPush(repo, user, token)
	uiCheckErr("push.ybPush", err)
}

func repositoryYbIsClean(repo *git.Repository) bool {
	wt, err := repo.Worktree()
	uiCheckErr("repositoryYbIsClean.Worktree", err)

	s, err := wt.Status()
	uiCheckErr("repositoryYbIsClean.Status", err)

	return s.IsClean()
}

func repositoryYbCommit(repo *git.Repository, tag string) {
	w, err := repo.Worktree()
	_, err = w.Add(".")
	uiCheckErr("repositoryYbCommit.Add", err)

	_, err = w.Commit(tag, &git.CommitOptions{
		Author: &object.Signature{
			When: time.Now(),
		},
	})
	uiCheckErr("repositoryYbCommit.Commit", err)
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
	return fmt.Sprintf("https://git.yottab.io/%s/%s.git", user, app)
}
