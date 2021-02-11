# configmap-watcher-example

This repo contains some examples of loading and watching Kubernetes configmaps.  As of this point in time, all of the examples use the `knative.dev/pkg/configmap` package, which I found by chance and thought both useful and very straightforward (borderline insultingly easy) to use.  All code is free to borrow/steal/copy/do what you want with it.

There are three directories, each containing a different example:
- [load](load) contains an example showing possibly the simplest method I've yet found for loading a configmap using the `configmap.Load()` method.
- [watch](watch) contains an example showing how to use knative's `InformedWatcher` to watch a configmap for changes and handle them dynamically.
- [watch-with-defaults](watch-with-defaults) contains an example showing how to use the same `InformedWatcher` type to both watch a configmap and to supply defaults in case the map cannot be loaded or is deleted while the application is running.

There is also a `build` script which will simply build three docker images, one for each example, and push them to the container registry specified by the `REGISTRY` environment variable.  You must be logged into this registry for the script to succeed.

Last, there is a `clusters.properties` file you can use to create a configmap for trying out these examples, by running the command `kubectl create cm clusters-config --from-env-file=clusters.properties`.  The examples are hardcoded to use `clusters-config` as the name of the configmap, just change that in the code and rebuild the docker images if you want to use a different name.
