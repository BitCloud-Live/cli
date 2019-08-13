# yb-cli
YOTTAb cli 

# How to build using make
```sh
$ make clean
$ make -B build
```
# Prebuilt binaries
Common platform binaries are published on the releases page. this includes linux, osx, windows and even arm binary for arm  linux platforms such as raspbian.
## Quick installation
The easiest way to install the latest release on Mac or Linux is with the following script:

```
curl -L http://install.yottab.io | bash
```
### Other installation methods:
#### OSX 
```sh
$ wget https://github.com/yottab/cli/releases/download/v${LATEST_VERSION}/yb-${LATEST_VERSION}-darwin-amd64 -O /usr/local/bin/yb
```
#### Linux
```sh
$ sudo wget https://github.com/yottab/cli/releases/download/v${LATEST_VERSION}/yb-v${LATEST_VERSION}-linux-amd64 -O /usr/local/bin/yb
$ sudo chmod +x /usr/local/bin/yb
```  
#### Windows 
Just grap the latest executive file in release page and run it from cmd.  
See [Releases](https://github.com/yottab/cli/releases).

## Update
YOTTAb cli can be updated in place using the following command:  
```sh
$ yb update
```  
# Quickstart
[![asciicast](https://asciinema.org/a/243119.svg)](https://asciinema.org/a/243119)
See [Documentations](http://docs.yottab.io/quickstart/) for more details.

# Gitlab integration
Test, build & deploy can be automated using gitlab ci.  
An example of .gitlab-ci.yaml configuration file is as follow:
```yaml
# This file is a template, and might need editing before it works on your project.
# Official docker image.
image: docker:latest
services:
  - docker:dind
  
stages:
  - test
  - build
  - deploy

variables:
  LINK: controller.yottab.io:443
  #Configure this variable in Secure Variables:
  YOTTAb_USER: <username>
  # YOTTAb_PASSWORD:  ----> It's more secure to be setted from settings -> ci/cd -> variables. 
  APP: rec
  IMAGE: "hub.yottab.io/$YOTTAb_USER/$APP:$CI_COMMIT_REF_NAME"

# Test stages
test1:
  stage: test
  before_script: []
  script:
    - echo run tests

test2:
  stage: test
  before_script: []
  script:
    - echo run tests


# Build stages
before_script:
   - docker login -u "$YOTTAB_USER" -p "$YOTTAB_PASSWORD" hub.yottab.io

build-master:

  stage: build
  script:
    - docker build --pull -t "$IMAGE" .
    - docker push "$IMAGE"
  only:
    - master # use `tags` for build when ever a new tag pushed to the repository

deploy: 
  image: hub.yottab.io/yottab-library/cli:${LATEST_VERSION} # use an existing git tag for the LATEST_VERSION variable.
  stage: deploy
  before_script:
    - yb login -u $YOTTAB_USER --P $YOTTAB_PASSWORD
    - yb application $APP
  script:
    - "yb update application $APP -i $IMAGE"
    - "yb application $APP"
```
