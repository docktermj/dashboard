# dashboard

## Prerequisites

## Develop

### Prerequisite software

The following software programs need to be installed:

1. [git](https://github.com/Senzing/knowledge-base/blob/master/HOWTO/install-git.md)
1. [make](https://github.com/Senzing/knowledge-base/blob/master/HOWTO/install-make.md)
1. [docker](https://github.com/Senzing/knowledge-base/blob/master/HOWTO/install-docker.md)

### Clone repository

1. Set these environment variable values:

    ```console
    export GIT_ACCOUNT=docktermj
    export GIT_REPOSITORY=go-fileindex
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"
    ```

1. Follow steps in [clone-repository](https://github.com/docktermj/KnowledgeBase/blob/master/HowTo/clone-repository.md) to install the Git repository.

### Add subcommand

1. Use [Cobra generator](https://github.com/spf13/cobra/blob/master/cobra/README.md)
   to create subcommand skeleton files.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    cobra add exampleSubcommand
    ```

### View godoc

1. View godoc documentation.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    ./godoc-docker.sh
    ```

1. View [localhost:6060/pkg/cmd/](http://localhost:6060/pkg/cmd/).

### View database

1. Run docker container.
   Example:

    ```console
    docker run \
        --env SQLITE_DATABASE=go-fileindex.db \
        --name sqlite-web \
        --publish 4380:8080 \
        --rm \
        --volume ${GIT_REPOSITORY_DIR}:/data \
        coleifer/sqlite-web:latest
    ```

1. Visit [localhost:4380](http://localhost:4380)

### Local development

1. Generate `.go` files.
   In particular, static files.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make generate
    ```

   This makes `box/blob.go`.

1. Build binary.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make build
    ```

1. Run binary.
   Example:

    ```console
    ${GIT_REPOSITORY_DIR}/target/linux/go-fileindex --help
    ```

1. Run service.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    ./target/linux/go-fileindex service
    ```

   View [localhost:3000](http://localhost:3000)

### Build docker image for development

1. **Option #1:** Using `docker` command and GitHub.

    ```console
    sudo docker build \
      --tag senzing/stream-loader \
      https://github.com/senzing/stream-loader.git
    ```

   View service at [localhost:3000](http://localhost:3000)

1. **Option #2:** Using `docker` command and local repository.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    sudo docker build --tag senzing/stream-loader .
    ```

1. **Option #3:** Using `make` command.

    ```console
    cd ${GIT_REPOSITORY_DIR}
    sudo make docker-build
    ```

    Note: `sudo make docker-build-development-cache` can be used to create cached docker layers.
