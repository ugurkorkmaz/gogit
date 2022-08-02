package gogit

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

const Github string = "https://github.com/"
const Gitlab string = "https://gitlab.com/"
const Bitbucket string = "https://bitbucket.org/"

type Git struct {
	Host string
	User string
	Repo string
	Type string
	File string
}

func New() *Git {
	return &Git{}
}
func (git *Git) Run(prefix string, dir string) {
	err := git.Parser(prefix)
	if err != nil {
		log.Fatal(err)
	}
	url, err := git.GetURL()
	if err != nil {
		log.Fatal(err)
	}
	err = git.Download(url)
	if err != nil {
		log.Fatal(err)
	}
	err = git.Extract(dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Downloaded and extracted " + git.User + "/" + git.Repo + " to " + dir)
}

func (git *Git) GetURL() (string, error) {
	var url string
	var orl string = git.User + "/" + git.Repo
	switch git.Host {
	case "github":
		url = Github + orl + "/archive/" + git.Type + ".tar.gz"
	case "gitlab":
		url = Gitlab + orl + "/repository/archive.tar.gz?ref=" + git.Type
	case "bitbucket":
		url = Bitbucket + orl + "/get/" + git.Type + ".tar.gz"
	default:
		return "", errors.New("invalid server")
	}
	return url, nil
}
func (git *Git) Parser(url string) error {

	match := regexp.MustCompile(`(?m)(?:(\w+):)?([\w\.\-]+)/([\w\.\-]+)(?:[#]([\w\.\-]+))?`).FindStringSubmatch(url)

	if len(match) == 0 {
		return errors.New("invalid url")
	}
	if match[1] != "" {
		git.Host = match[1]
	}
	if match[2] != "" {
		git.User = match[2]
	}
	if match[3] != "" {
		git.Repo = match[3]
	}
	if match[4] != "" {
		git.Type = match[4]
	}
	if git.Type == "" {
		git.Type = "master"
	}
	if git.Host == "" {
		git.Host = "github"
	}
	if git.User == "" {
		return errors.New("user is empty")
	}
	if git.Repo == "" {
		return errors.New("repo is empty")
	}
	git.File = git.Repo + "-" + git.Type + ".tar.gz"
	return nil
}
func (git *Git) Download(path string) error {
	request, err := http.Get(path)
	if err != nil {
		return err
	}
	defer request.Body.Close()
	if request.StatusCode != 200 {
		return errors.New("git server error code: " + request.Status + " ")
	} else {

		file, err := os.Create(git.File)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, request.Body)
		if err != nil {
			return err
		}
	}
	return nil
}
func (git *Git) Extract(dir string) error {
	file, err := os.Open(git.File)
	if err != nil {
		return err
	}
	err = Untar(file)
	if err != nil {
		return err
	}
	file.Close()
	err = os.Remove(git.File)
	if err != nil {
		return err
	}
	oldpath := filepath.Join(git.Repo + "-" + git.Type)
	newpath := filepath.Join(dir)

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
