# Healthz

[![Build Status](https://travis-ci.org/cafebazaar/healthz.svg)](https://travis-ci.org/cafebazaar/healthz) [![GoDoc](https://godoc.org/github.com/cafebazaar/healthz?status.svg)](https://godoc.org/github.com/cafebazaar/healthz)

![Screenshot of HTML Report - Healthz](https://github.com/cafebazaar/healthz/raw/master/demo/Healthz.png)

## Examples
  * [healthz.NewHandler()](https://godoc.org/github.com/cafebazaar/healthz#example-NewHandler)

## Features
* Overall Health Calculation using Subcomponents' Severity and Health Levels
* HTML Report
* JSON Report
* [Kubernetes Liveness and Readiness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/)
  * Reports the service as *live* when the overall health is over `Error`
  * Reports the service as *ready* when the overall health is at least on `Normal`

## TODO
* gRPC Report
