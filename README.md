# uv-cli
UVCloud cli 

# How to build using make
```sh
$ make clean
$ make -B build
```
# Prebuilt binaries
Common platform binaries are published on the releases page. this includes linux, osx, windows and even arm binary for arm  linux platforms such as raspbian.
See [Releases](https://github.com/uvcloud/uv-cli/releases).

# Quickstart
[![asciicast](https://asciinema.org/a/193296.png)](https://asciinema.org/a/193296)
See [Documentations](http://docs.uvcloud.ir/quickstart/) for more details.

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
  LINK: controller.uvcloud.ir:8443
  #Configure this variable in Secure Variables:
  UVCLOUD_USER: <username>
  # UVCLOUD_PASSWORD:  ----> It's more secure to be setted from settings -> ci/cd -> variables. 
  APP: rec
  IMAGE: "hub.uvcloud.ir/$UVCLOUD_USER/$APP:$CI_COMMIT_REF_NAME"

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
   - docker login -u "$UVCLOUD_USER" -p "$UVCLOUD_PASSWORD" hub.uvcloud.ir

build-master:

  stage: build
  script:
    - docker build --pull -t "$IMAGE" .
    - docker push "$IMAGE"
  only:
    - master # use `tags` for build when ever a new tag pushed to the repository

deploy: 
  image: hub.uvcloud.ir/uvcloud/uv-cli:v1.0.0-rc7
  stage: deploy
  before_script:
    - echo $UVCLOUD_PASSWORD | uv-cli login -u $UVCLOUD_USER
    - uv-cli app:info -n $APP
  script:
    - "uv-cli app:configSet -n $APP  -i $IMAGE"
    - "uv-cli app:info -n $APP"
```
