# configmap-watcher-example

This repo contains some examples of loading and watching Kubernetes configmaps.  While most of the examples use the `knative.dev/pkg/configmap` package, which I found by chance and thought both useful and very straightforward (borderline insultingly easy) to use, the `watch` example just uses the vanilla [client-go](github.com/kubernetes/client-go) package.  All code is free to borrow/steal/copy/do what you want with it.

There are four sample directories, each containing a different example:
- [load](load) contains an example showing possibly the simplest method I've yet found for loading a configmap using the `configmap.Load()` method.
- [watch](watch) contains an example showing how to implement a simple and fairly robust watcher using [client-go](github.com/kubernetes/client-go).  This watcher has a default value to be used if the config map does not exist at startup or gets deleted later, and it automatically updates its information based on changes to the configmap (this is pretty much the point of using watchers).
- [kn-watch](kn-watch) contains an example showing how to use knative's `InformedWatcher` to watch a configmap for changes and handle them dynamically.
- [watch-with-defaults](watch-with-defaults) contains an example showing how to use the same `InformedWatcher` type to both watch a configmap and to supply defaults in case the map cannot be loaded or is deleted while the application is running.  This is basically a somewhat simpler version of my own `watch` implementation showing how you can reuse the knative package to simplify your code, if you like.

There is also a `build` script which will simply build four docker images, one for each example, and push them to the container registry specified by the `REGISTRY` environment variable.  You must be logged into this registry for the script to succeed.

Last, in the [configs](configs) directory there is a `clusters.properties` file you can use to create a configmap for trying out these examples, by running the command `kubectl create cm clusters-config --from-env-file=clusters.properties`.  The examples are hardcoded to use `clusters-config` as the name of the configmap, just change that in the code and rebuild the docker images if you want to use a different name.
