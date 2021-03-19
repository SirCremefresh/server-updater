#!/bin/sh

UPDATE_LEVEL="$1"

if [ -z "$1" ]
  then
    echo "The update level needs to be suplied in the arguments"
    exit 1
fi

if [ "$UPDATE_LEVEL" != "MAJOR" ] && [ "$UPDATE_LEVEL" != "MINOR" ] && [ "$UPDATE_LEVEL" != "PATCH" ] 
  then
    echo "The update level needs to be one of the following values. MAJOR, MINOR or PATCH"
    exit 1
fi

#Get the highest tag number
VERSION=`git describe --abbrev=0 --tags`
VERSION=${VERSION:-'0.0.0'}

#Get number parts
MAJOR="${VERSION%%.*}"; VERSION="${VERSION#*.}"
MINOR="${VERSION%%.*}"; VERSION="${VERSION#*.}"
PATCH="${VERSION%%.*}"; VERSION="${VERSION#*.}"

#Increase version
if [ "$UPDATE_LEVEL" = "MAJOR" ] 
  then
    MAJOR=$((MAJOR+1))
    MINOR=0
    PATCH=0
fi
if [ "$UPDATE_LEVEL" = "MINOR" ] 
  then
    MINOR=$((MINOR+1))
    PATCH=0
fi
if [ "$UPDATE_LEVEL" = "PATCH" ] 
  then
    PATCH=$((PATCH+1))
fi


#Create new tag
NEW_TAG="$MAJOR.$MINOR.$PATCH"
echo "Updating to $NEW_TAG"

read -r "publish tag (y/N)?" CONT
if [ "$CONT" = "y" ]; then
  echo "publishing"
  git tag $NEW_TAG
  git push --tag
else
  echo "not publishing"
fi

