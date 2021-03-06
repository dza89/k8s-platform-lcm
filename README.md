# Kubernetes platform lifecycle management

![Build](https://github.com/arminc/k8s-platform-lcm/workflows/test/badge.svg?branch=master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/bf91530a165c4113b1691353ae51ae1d)](https://app.codacy.com/manual/armin/k8s-platform-lcm?utm_source=github.com&utm_medium=referral&utm_content=arminc/k8s-platform-lcm&utm_campaign=Badge_Grade_Dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/arminc/k8s-platform-lcm)](https://goreportcard.com/report/github.com/arminc/k8s-platform-lcm)
[![CodeFactor](https://www.codefactor.io/repository/github/arminc/k8s-platform-lcm/badge)](https://www.codefactor.io/repository/github/arminc/k8s-platform-lcm)
[![codecov](https://codecov.io/gh/arminc/k8s-platform-lcm/branch/master/graph/badge.svg)](https://codecov.io/gh/arminc/k8s-platform-lcm)

<img src="https://github.com/arminc/k8s-platform-lcm/raw/master/assets/logo.png" width="100">

----

This project helps you keep track of all your software and tools that are used or running in and around your Kubernetes platform. It helps you with part of the lifecycle management to keep your software up to data for feature completeness, security or compliance reasons.

## Features

-   [x] Keep track of versions of all the running containers (including init containers) inside the Kubernetes
-   [x] Keep track of new image versions. Supporting Quay, Gcr, Docker hub, Jfrog Artifactory by default
-   [x] Works with private registries and private images
-   [x] Allow overriding of the registry to search latest versions from another registry
-   [x] Keep track of image vulnerabilities using Jfrog Xray
-   [x] Possibility to provide local tool versions (like terraform) and find the new versions on GitHub
-   [x] Keep track of Helm chart deployments and track new versions of the charts
-   [x] Present the information command line
-   [x] Present the information trough a web UI
-   [x] Export prometheus metrics
-   [x] Filter out vulnerabilties which are accepted

## Help (how to run)

For all the configuration options please have a look at the [exampleConfig.yaml](exampleConfig.yaml).

When running lcm you can provide certain flags which are not available in the config. The application assumes there is a config.yaml available in the same folder.

```bash
./lcm --help
usage: lcm [<flags>]

Kubernetes platform lifecycle management

Flags:
  --help                  Show context-sensitive help (also try --help-long and --help-man).
  --version               Show application version.
  --config="config.yaml"  Provide the path to the config file. Default is config.yaml which is in the same folder as lcm
  --local                 Run locally, default expected behavior is to run in the Kubernetes cluster
  --verbose               Show more information. This overrides the config setting
  --debug                 Show debug information, debug includes verbose. This overrides the config setting
  --jsonLogging           Log in json format
  --logFile=LOGFILE       Log file path
  --server                Start the server
  --metrics               Start the metric server (runs on port 9572)
  --vul                   Print all vulnerabilities at the end
```

**Note:** If you are using `--server` option please make sure the `templates` and `static` folder are next to the binary so it can serve the page.

### Docker

Docker image is available at **arminc/lcm:VERSION** or **arminc/lcm:latest**. It is packaged with the template and css. Run it as following, and add any necessary flags you want or use the yaml file.

```bash
docker run -it -v $(pwd)/config.yaml:/config.yaml -p 7321:7321 arminc/lcm:latest --local --server
```

## Example output

### Command Line

```bash
+---------------------------------------+-------------------+----------+-------+
|                 IMAGE                 |      VERSION      |  LATEST  | CVES  |
+---------------------------------------+-------------------+----------+-------+
| library/alpine                        |      3.10.1       |  3.10.3  | ERROR |
| openpolicyagent/kube-mgmt             |        0.9        |   0.10   | 0     |
| openpolicyagent/opa                   |      0.14.1       |  0.15.1  | 0     |
| velero/velero                         |      v1.1.0       |  v1.2.0  | 0     |
+---------------------------------------+-------------------+----------+-------+
+----------------------------+------------+----------+
|           CHART            |  VERSION   |  LATEST  |
+----------------------------+------------+----------+
| opa                        |   0.12.0   |  1.13.1  |
| velero                     |   2.5.0    |  2.7.0   |
+----------------------------+------------+----------+
+---------------------+---------+----------+
|        TOOL         | VERSION |  LATEST  |
+---------------------+---------+----------+
| derailed/popeye     | v0.4.1  |  v0.5.0  |
| hashicorp/terraform | 0.11.14 | v0.12.18 |
+---------------------+---------+----------+
```

### Metric output

```text
chart_info{chart="polaris",latestVersion="1.1.0",version="0.10.1"} 0
image_info{image="storageos/csi-provisioner",latestVersion="v1.4.0",registry="docker.io",version="v1.4.0"} 1
tool_info{latestVersion="v0.12.26",tool="hashicorp/terraform",version="0.11.14"} 0
```

### Web UI

<img src="assets/screenshot.png" width="800">
