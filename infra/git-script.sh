#!/bin/bash

# CRLF
# config in local (only this project)
git config core.autocrlf input

# pre-commit hook
cp ./pre-commit ../.git/hooks/pre-commit

