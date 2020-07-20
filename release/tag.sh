#!/usr/bin/env bash

GITHUB_ROOT="$HOME/github/free5gc"
BITBUCKET_ROOT='.'
TAG=v2020-07-21-01

OLDIFS="$IFS"

IFS=$'\n'
modules=`git config --file .gitmodules --get-regexp '\.path' | awk '{ print $2 }'`

for module in ${modules};
do
	GITHUB_PATH="${GITHUB_ROOT}/${module}"
	echo "==== Start ${GITHUB_PATH} ===="
	cd ${GITHUB_PATH}
	git tag ${TAG}
	git push origin ${TAG}
	echo "==== end ===="
	cd ${GITHUB_ROOT}
	git add ${module}
done

IFS=$OLDIFS
