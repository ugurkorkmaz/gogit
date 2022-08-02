<h1 align="center">Gogit</h1>

<p align="center">
  <img alt="Github top language" src="https://img.shields.io/github/languages/top/extendswork/gogit?color=56BEB8">

  <img alt="Github language count" src="https://img.shields.io/github/languages/count/extendswork/gogit?color=56BEB8">

  <img alt="Repository size" src="https://img.shields.io/github/repo-size/extendswork/gogit?color=56BEB8">

  <img alt="License" src="https://img.shields.io/github/license/extendswork/gogit?color=56BEB8">
</p>

<br>

## :dart: About ##

``gogit`` makes copies of git repositories.This is much quicker than using git clone, because you're not downloading the entire git history.


## :rocket: Technologies ##

The following tools were used in this project:

- [Golang](https://go.dev/)

## :checkered_flag: Starting ##

```bash
# Install gogit
$ go install github.com/extendswork/gogit/cmd/gogit

```
Default git server GitHub is used.
In addition, you can use Gitlab and Bitbucket.
You can download release versions and versioned versions

```bash	
# Start gogit
gogit user/repo my-project
gogit github:user/repo my-project
gogit github:user/repo#v1.0.0 my-project
gogit github:user/repo#master my-project

gogit gitlab:user/repo my-project
gogit gitlab:user/repo#v1.0.0 my-project
gogit gitlab:user/repo#master my-project

gogit bitbucket:user/repo my-project
gogit bitbucket:user/repo#v1.0.0 my-project
gogit bitbucket:user/repo#master my-project
```	

## :memo: License ##

This project is under license from MIT. For more details, see the [LICENSE](LICENSE.md) file.


Made with :heart: by <a href="https://github.com/extendswork" target="_blank">Extends Work</a>

&#xa0;
## See also
- [Degit](https://github.com/Rich-Harris/degit) by [Rich Harris](https://github.com/Rich-Harris)
<br/>&#xa0;
<br/>
<a href="#top">Back to top</a>
