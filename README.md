# visgostruct

## Overview

Visualize relations of structs in golang by PlantUML.

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
    visgostruct -f -t ~/go/workspace/visgostruct/*.go
    # comments/tags shown in note
    visgostruct -f -c -t -n ~/go/workspace/visgostruct/*.go
    # specify including or excluding struct's name with regexp
    visgostruct -i 'Information$' -e '^Field' ~/go/workspace/visgostruct/*.go
    # specify root struct and extract its children, descendant
    visgostruct -r 'StructInformation' -f -n ~/go/workspace/visgostruct/*.go
    # print definitions as CSV
    visgostruct -r 'StructInformation' -csv ~/go/workspace/visgostruct/*.go
    # print definitions as TSV
    visgostruct -r 'StructInformation' -tsv ~/go/workspace/visgostruct/*.go
    # print definitions in flat style
    visgostruct -r 'StructInformation' -csv -flat ~/go/workspace/visgostruct/*.go

## Requirements

github.com/urfave/cli

    go get github.com/urfave/cli

## Installation

    go build
    go install
