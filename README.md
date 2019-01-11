# visgostruct

## Overview

Visualize relations of struct in golang by PlantUML.

## Description

CLI tool to parse golang sources, then print definitions and relations of structs by PlantUML(http://plantuml.com/).

## Synopsis
    # so simple output
    visgostruct ~/go/workspace/visgostruct/*.go
    # incude fields definitions
    visgostruct -f ~/go/workspace/visgostruct/*.go
    # include comments
    visgostruct -f -c ~/go/workspace/visgostruct/*.go
    # include tags
    visgostruct -f -c -t ~/go/workspace/visgostruct/*.go
    # comments and tags shown in note
    visgostruct -f -c -n ~/go/workspace/visgostruct/*.go
    # specify including or excluding struct's name with regexp
    visgostruct -i 'Information$' -e '^Field' ~/go/workspace/visgostruct/*.go

## Requirements

github.com/urfave/cli

    go get github.com/urfave/cli

## Installation

    go build
    go install
