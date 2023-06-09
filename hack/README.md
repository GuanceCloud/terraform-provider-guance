# Guance Community hack GuideLines

## What is this?

The hack directory contains scripts used to build, test, package, and release this project.

Many Kubernetes projects use this pattern. We use it in Guance Cloud to ensure that all contributors use the same build and release process.

## How to use it?

`make` will call the mage step to complete the workflow. Follow commands is supported.

```bash
Targets:
  build:install    run installation on the provider to local
  dev:all          run all dev tasks
  dev:fmt          run format for all the code
  dev:lint         run lint for all the code
  gen:doc          run generator over the documentation
  test:acc         run acceptance test for specified resource
```

## References

* [Kubernetes hack GuideLines](https://github.com/kubernetes/kubernetes/tree/v1.26.1/hack)
* [Dagger hack Overview](https://github.com/dagger/dagger/tree/main/hack)
