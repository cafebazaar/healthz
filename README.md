# Healthz

[![Build Status](https://travis-ci.org/cafebazaar/healthz.svg)](https://travis-ci.org/cafebazaar/healthz) [![GoDoc](https://godoc.org/github.com/cafebazaar/healthz?status.svg)](https://godoc.org/github.com/cafebazaar/healthz)

![Screenshot of HTML Report - Healthz](https://github.com/cafebazaar/healthz/raw/master/demo/Healthz.png)

* HTML Report
* JSON Report
* [Kubernetes Liveness and Readiness](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/)
  * Reports the service as *live* when the overall health is at least on `Unknown`
  * Reports the service as *ready* when the overall health is at least on `Normal`
