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
    export GIT_REPOSITORY=dashboard
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Follow steps in [clone-repository](https://github.com/docktermj/KnowledgeBase/blob/master/HowTo/clone-repository.md) to install the Git repository.

1. View [localhost:6060/pkg/cmd/](http://localhost:6060/pkg/cmd/).

### Local development

1. Run code.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make run

    ```

### Local web development

1. To view the user interface,
   with no underlying enablement.
   Example:

    ```console
    docker run \
        --publish 8259:80 \
        --volume ${GIT_REPOSITORY_DIR}/dashboard/static:/usr/share/nginx/html:ro \
        nginx

    ```

    Visit [localhost:8259](http://localhost:8259)

### Build binary

1. Build binary.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make build

    ```

1. Run binary.
   Example:

    ```console
    ${GIT_REPOSITORY_DIR}/target/linux/dashboard --help

    ```

1. Run service.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    ./target/linux/dashboard

    ```

   View [localhost:8259](http://localhost:8259)

### Get static files

1. Set these environment variable values:

    ```console
    export GIT_ACCOUNT=docktermj
    export GIT_REPOSITORY=dashboard
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Bootstrap (with Popper).
   Example:

    ```console
    wget \
        --output-document ${GIT_REPOSITORY_DIR}/static/css/bootstrap.min.css \
        https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css

    wget \
        --output-document ${GIT_REPOSITORY_DIR}/static/js/bootstrap.bundle.min.js \
        https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.bundle.min.js

   ```

1. JQuery.
   Example:

    ```console
    wget \
        --output-document ${GIT_REPOSITORY_DIR}/static/js/jquery.min.js \
        https://code.jquery.com/jquery-3.6.3.min.js

    wget \
        --output-document ${GIT_REPOSITORY_DIR}/static/css/jquery.dataTables.min.css \
        https://cdn.datatables.net/1.13.1/css/jquery.dataTables.min.css

    wget \
        --output-document ${GIT_REPOSITORY_DIR}/static/js/jquery.dataTables.min.js \
        https://cdn.datatables.net/1.13.1/js/jquery.dataTables.min.js

   ```
