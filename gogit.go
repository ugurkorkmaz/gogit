package gogit

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	_github    string = "https://github.com/%s/%s/archive/%s.tar.gz"
	_gitlab    string = "https://gitlab.com/%s/%s/repository/archive.tar.gz?ref=%s"
	_bitbucket string = "https://bitbucket.org/%s/%s/get/%s.tar.gz"
)

const _regex string = `(?m)(?:(\w+):)?([\w\.\-]+)/([\w\.\-]+)(?:[#]([\w\.\-]+))?`

var (
	ErrInvalidServer     = errors.New("invalid git server")
	ErrInvalidUrl        = errors.New("invalid url")
	ErrUserNotFoundInUrl = errors.New("user not found in url")
	ErrRepoNotFoundInUrl = errors.New("repo not found in url")
)

type git struct {
	// Base address of the Git server.
	host string
	// Your username on the Git server.
	user string
	// Your repo on the Git server.
	repo string
	// The name of the version or branch of the repo.
	types string
	// The name of the file downloaded to the local.
	file string
}

func New() *git {
	return &git{}
}
func (git *git) Run(prefix string, dir string) error {

	if err := git.Parser(prefix); err != nil {
		return err
	}
	url, err := git.GetDownloadURL()
	if err != nil {
		return err
	}
	if err := git.Download(url); err != nil {
		return err
	}
	if err := git.Extract(dir); err != nil {
		return err
	}
	return nil
}

func (git *git) GetDownloadURL() (string, error) {
	if git.host == "github" {
		return fmt.Sprintf(_github, git.user, git.repo, git.types), nil
	}
	if git.host == "gitlab" {
		return fmt.Sprintf(_gitlab, git.user, git.repo, git.types), nil
	}
	if git.host == "bitbucket" {
		return fmt.Sprintf(_bitbucket, git.user, git.repo, git.types), nil
	}
	return "", ErrInvalidServer
}
func (git *git) Parser(url string) error {

	match := regexp.MustCompile(_regex).FindStringSubmatch(url)

	if len(match) == 0 {
		return ErrInvalidUrl
	}
	if match[1] != "" {
		git.host = match[1]
	}
	if match[2] != "" {
		git.user = match[2]
	}
	if match[3] != "" {
		git.repo = match[3]
	}
	if match[4] != "" {
		git.types = match[4]
	}
	if git.types == "" {
		git.types = "master"
	}
	if git.host == "" {
		git.host = "github"
	}
	if git.user == "" {
		return ErrUserNotFoundInUrl
	}
	if git.repo == "" {
		return ErrRepoNotFoundInUrl
	}
	git.file = git.repo + "-" + git.types + ".tar.gz"
	return nil
}
func (git *git) Download(path string) error {
	request, err := http.Get(path)
	if err != nil {
		return err
	}
	defer request.Body.Close()
	if request.StatusCode != 200 {
		return errors.New(strings.ToLower(http.StatusText(request.StatusCode)))
	} else {

		file, err := os.Create(git.file)
		if err != nil {
			return errors.New("")
		}
		defer file.Close()
		_, err = io.Copy(file, request.Body)
		if err != nil {
			return err
		}
	}
	return nil
}
func (git *git) Extract(dir string) error {
	file, err := os.Open(git.file)
	if err != nil {
		return err
	}
	err = Untar(file)
	if err != nil {
		return err
	}
	file.Close()
	err = os.Remove(git.file)
	if err != nil {
		return err
	}
	oldpath := filepath.Join(git.repo + "-" + git.types)
	newpath := filepath.Join(dir)
	err = os.Chmod(newpath, 0777)
	if err != nil {
		return err
	}
	err = os.Rename(oldpath, newpath)
	if err != nil {
		return err
	}
	return nil
}

func Untar(r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(header.Name)

		// check the child/dir type

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
