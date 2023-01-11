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

   Visit [localhost:8259](http://localhost:8259)

### Package

#### Package RPM and DEB files

1. Use make target to run a docker images that builds RPM and DEB files.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}
    make package

    ```

1. The results will be in the `${GIT_REPOSITORY_DIR}/target` directory.
   Example:

    ```console
    tree ${GIT_REPOSITORY_DIR}/target

    ```

#### Test DEB package on Ubuntu

1. Determine if `dashboard` is installed.
   Example:

    ```console
    apt list --installed | grep dashboard

    ```

1. :pencil2: Install `dashboard`.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}/target
    sudo apt install ./dashboard-0.0.0.deb

    ```

1. Run command.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    dashboard

    ```

   Visit [localhost:8259](http://localhost:8259)

1. Remove `dashboard` from system.
   Example:

    ```console
    sudo apt-get remove dashboard

    ```

#### Test RPM package on Centos

1. Determine if `dashboard` is installed.
   Example:

    ```console
    yum list installed | grep dashboard

    ```

1. :pencil2: Install `dashboard`.
   Example:

    ```console
    cd ${GIT_REPOSITORY_DIR}/target
    sudo yum install ./dashboard-0.0.0.rpm

    ```

1. Run command.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    dashboard

    ```

1. Remove `dashboard` from system.
   Example:

    ```console
    sudo yum remove dashboard

    ```

### Make documents

Make documents visible at
[hub.senzing.com/dashboard](https://hub.senzing.com/dashboard).

1. Identify repository.
   Example:

    ```console
    export GIT_ACCOUNT=senzing
    export GIT_REPOSITORY=dashboard
    export GIT_ACCOUNT_DIR=~/${GIT_ACCOUNT}.git
    export GIT_REPOSITORY_DIR="${GIT_ACCOUNT_DIR}/${GIT_REPOSITORY}"

    ```

1. Make documents.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    dashboard docs --dir ${GIT_REPOSITORY_DIR}/docs

    ```

### How to update static files

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
        --output-document ${GIT_REPOSITORY_DIR}/dashboard/static/css/bootstrap.min.css \
        https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css

    wget \
        --output-document ${GIT_REPOSITORY_DIR}/dashboard/static/css/bootstrap.min.css.map \
        https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css.map

    wget \
        --output-document ${GIT_REPOSITORY_DIR}/dashboard/static/js/bootstrap.bundle.min.js \
        https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.bundle.min.js

   ```

1. JQuery.
   Example:

    ```console
    wget \
        --output-document ${GIT_REPOSITORY_DIR}/dashboard/static/js/jquery.min.js \
        https://code.jquery.com/jquery-3.6.3.min.js

    wget \
        --output-document ${GIT_REPOSITORY_DIR}/dashboard/static/css/jquery.dataTables.min.css \
        https://cdn.datatables.net/1.13.1/css/jquery.dataTables.min.css

    wget \
        --output-document ${GIT_REPOSITORY_DIR}/dashboard/static/js/jquery.dataTables.min.js \
        https://cdn.datatables.net/1.13.1/js/jquery.dataTables.min.js

   ```

1. Dashboard.
   Example:

    ```console
        export BOOTSTRAP_VERSION=5.2.3

        wget \
        --output-document /tmp/bootstrap-${BOOTSTRAP_VERSION}-examples.zip \
        https://github.com/twbs/bootstrap/releases/download/v${BOOTSTRAP_VERSION}/bootstrap-${BOOTSTRAP_VERSION}-examples.zip

        unzip \
            /tmp/bootstrap-${BOOTSTRAP_VERSION}-examples.zip \
            -d /tmp

        cp --force /tmp/bootstrap-${BOOTSTRAP_VERSION}-examples/dashboard/dashboard.css  ${GIT_REPOSITORY_DIR}/dashboard/static/css/dashboard.css
        cp --force /tmp/bootstrap-${BOOTSTRAP_VERSION}-examples/dashboard/dashboard.js   ${GIT_REPOSITORY_DIR}/dashboard/static/js/dashboard.js

    ```

1. Feather icons.
   Example:

    ```console
    wget \
        --output-document ${GIT_REPOSITORY_DIR}/dashboard/static/js/feather.min.js \
        https://cdn.jsdelivr.net/npm/feather-icons/dist/feather.min.js

    ```
