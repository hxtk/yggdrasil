# Verdandi

## Description

Verdandi is a kubernetes-native unprivileged NTP server. In the future it aims to support ANTP.

Time is synchronized with an upstream time server and maintained as local state within the application, using the Kernel monotonic clock to interpolate time readings between synchronizations.

For optimal performance, one optionally grant CAP\_SYS\_NICE to ensure that the process will not be preempted.

## Requirements

1. The application shall implement NTP and ANTP servers.
2. The application shall use an ANTP client to synchronize with upstream (A)NTP servers.
3. The application shall, by default, serve time at a stratum one greater than the upstream server from which the most recent synchronization was derived.
4. The application shall allow the user to specify the stratum, with a minimum value of 2 permitted.
5. The application shall interpolate between synchronizations using the system monotonic clock.
6. The application shall be capable of executing without any root privileges.
7. The application shall be capable of running with a specified niceness, if it is granted CAP\_SYS\_NICE.
8. The application shall allow a configurable amount of time for which it is considered healthy between synchronizations.
9. The application shall expose an Authn/Authz-integrated HTTPS endpoint at the path `/metrics` which shall expose metrics in OpenTelemetry format.
10. Application metrics shall support the tracking of successful synchronizations.
11. Application metrics shall support the tracking of failed synchronizations.
12. Application metrics shall support the tracking of request latency.
13. Application metrics shall support the tracking of server liveness.
14. The application shall expose an Authn/Authz-integrated HTTPS endpoint at the `/healthz` path which shall return with HTTP Status Code 200 if and only if the application has synchronized with an upstream server within the interval specified in (8).
15. The application container shall include an (A)NTP client application to use as a Liveness check.
