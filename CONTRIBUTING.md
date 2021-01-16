# Contributing
## Release process

In order to standardize the release process, [goreleaser](https://goreleaser.com) has been adopted.

To build and release a new version:
```
git tag vX.X.X && git push --tags
gorelease --rm-dist
```

The binary version is set automatically to the git tag.

Please tag new versions using [semantic versioning](https://semver.org/spec/v2.0.0.html).
