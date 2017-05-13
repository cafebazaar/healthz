/*
Package healthz is a status reporter library for microservices. It provides
a `Handler` which allows components in the app to store their status, and
exposes these reports in multiple formats:
 * HTML Report (`/`)
 * JSON Report (`/json`)
 * Liveness (`/liveness`): HTTP 200 only if the overall health is at least on `Unknown`
 * Readiness (`/readiness`): HTTP 200 only if the overall health is at least on `Normal`
*/
package healthz
