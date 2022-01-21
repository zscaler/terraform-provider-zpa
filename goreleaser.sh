#!/bin/bash

git tag -a v2.0.5 -m "Introduced ZPA Provider v2.0.5"
git push origin v2.0.5
goreleaser release