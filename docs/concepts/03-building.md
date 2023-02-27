## Building your project

Exciting! We are getting somewhere :)

You can run your app by:
```shell
make run
```

And you can build a container for it locally by:
```shell
make build-podman
```

Following information is all about what happens when you run these commands :)

### Go dependencies

We have three `make` targets for working with dependencies:

* `make download-deps` installs dependencies locally (aliased `make prep`)
* `make update-deps` updates dependencies to the newest versions
* `make tidy-deps` cleans up dependencies locally

### The binary

The main binary is called `api` and it serves as api application http server for our service.

It's entry point is `cmd/api/main.go`

### Containerization

Your app will run in the production environment in Container in OpenShift.
So lets package our app in a Container :)

You can build the container with Podman by running `make build-podman`.

There are two phases in the build. It is captured in [`Containerfile`](../../build/Containerfile).

#### First is the build phase.

We use the official Red Hat go-toolset build container to build our projects.
We have manually set the go version of the container, when bumping the go version, it needs to be bumped here.

The build itself is done by copying all project files in `/build` directory in the build container.
Followed by running `make prep build strip` which runs phases:

* Install dependencies
* Building binaries
* Stripping binaries of debug information (to keep them smaller).


```Dockerfile
FROM registry.access.redhat.com/ubi8/go-toolset:1:18 as build
USER 0
RUN mkdir /build
WORKDIR /build
COPY . .
RUN make prep build strip
```

#### Second phase produces the final container image

It just copies the binaries from the `build` container. 

```Dockerfile
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
COPY --from=build /build/api /api
USER 1001
CMD ["/api"]
```

We are done :)