# Documentation

> Yokai automated documentation, based on [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/).

<!-- TOC -->
* [Public version](#public-version)
* [Contributing](#contributing)
<!-- TOC -->

## Public version

Yokai public documentation is available at [https://ankorstore.github.io/yokai/](https://ankorstore.github.io/yokai/).

## Contributing

Before triggering the [docs workflow](../.github/workflows/docs.yml), you can test your documentation changes by running (from the repository root):

```shell
docker build -t squidfunk/mkdocs-material docs && docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material
```

This will make the documentation available locally on [http://localhost:8000](http://localhost:8000).