[![Validations](https://github.com/nextlinux/gosbom/actions/workflows/validations.yaml/badge.svg)](https://github.com/nextlinux/gosbom/actions/workflows/validations.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nextlinux/gosbom)](https://goreportcard.com/report/github.com/nextlinux/gosbom)
[![GitHub release](https://img.shields.io/github/release/nextlinux/gosbom.svg)](https://github.com/nextlinux/gosbom/releases/latest)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nextlinux/gosbom.svg)](https://github.com/nextlinux/gosbom)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/nextlinux/gosbom/blob/main/LICENSE)
[![Slack Invite](https://img.shields.io/badge/Slack-Join-blue?logo=slack)](https://nextlinux.com/slack)

A CLI tool and Go library for generating a Software Bill of Materials (SBOM) from container images and filesystems. Exceptional for vulnerability detection when used with a scanner like [Gosbom](https://github.com/nextlinux/gosbom).

#### Execution flow

```mermaid
sequenceDiagram
    participant main as cmd/gosbom/main
    participant cli as cli.New()
    participant root as root.Execute()
    participant cmd as <command>.Execute()

    main->>+cli: 

    Note right of cli: wire ALL CLI commands
    Note right of cli: add flags for ALL commands

    cli-->>-main:  root command 

    main->>+root: 
    root->>+cmd: 
    cmd-->>-root: (error)  

    root-->>-main: (error) 

    Note right of cmd: Execute SINGLE command from USER
```

## Features
- Generates SBOMs for container images, filesystems, archives, and more to discover packages and libraries
- Supports OCI, Docker and [Singularity](https://github.com/sylabs/singularity) image formats
- Linux distribution identification
- Works seamlessly with [Grype](https://github.com/nextlinux/grype) (a fast, modern vulnerability scanner)
- Able to create signed SBOM attestations using the [in-toto specification](https://github.com/in-toto/attestation/blob/main/spec/README.md)
- Convert between SBOM formats, such as CycloneDX, SPDX, and Gosbom's own format.

### Supported Ecosystems

- Alpine (apk)
- C (conan)
- C++ (conan)
- Dart (pubs)
- Debian (dpkg)
- Dotnet (deps.json)
- Objective-C (cocoapods)
- Elixir (mix)
- Erlang (rebar3)
- Go (go.mod, Go binaries)
- Haskell (cabal, stack)
- Java (jar, ear, war, par, sar, nar, native-image)
- JavaScript (npm, yarn)
- Jenkins Plugins (jpi, hpi)
- Linux kernel archives (vmlinz)
- Linux kernel modules (ko)
- Nix (outputs in /nix/store)
- PHP (composer)
- Python (wheel, egg, poetry, requirements.txt)
- Red Hat (rpm)
- Ruby (gem)
- Rust (cargo.lock)
- Swift (cocoapods)

## Installation

**Note**: Currently, Gosbom is built only for Linux, macOS and Windows.

### Recommended
```bash
curl -sSfL https://raw.githubusercontent.com/nextlinux/gosbom/main/install.sh | sh -s -- -b /usr/local/bin
```

... or, you can specify a release version and destination directory for the installation:

```
curl -sSfL https://raw.githubusercontent.com/nextlinux/gosbom/main/install.sh | sh -s -- -b <DESTINATION_DIR> <RELEASE_VERSION>
```

### Chocolatey

The chocolatey distribution of gosbom is community maintained and not distributed by the nextlinux team

```powershell
choco install gosbom -y
```

### Homebrew
```bash
brew install gosbom
```

### Nix

**Note**: Nix packaging of Gosbom is [community maintained](https://github.com/NixOS/nixpkgs/blob/master/pkgs/tools/admin/gosbom/default.nix). Gosbom is available in the [stable channel](https://nixos.wiki/wiki/Nix_channels#The_official_channels) since NixOS `22.05`.

```bash
nix-env -i gosbom
```

... or, just try it out in an ephemeral nix shell:

```bash
nix-shell -p gosbom
```

## Getting started

### SBOM

To generate an SBOM for a container image:

```
gosbom <image>
```

The above output includes only software that is visible in the container (i.e., the squashed representation of the image). To include software from all image layers in the SBOM, regardless of its presence in the final image, provide `--scope all-layers`:

```
gosbom <image> --scope all-layers
```

#### Code example of gosbom as a library

Here is a gist of using gosbom as a library to generate a SBOM for a docker image: [link](https://gist.github.com/wagoodman/57ed59a6d57600c23913071b8470175b).
The execution flow for the example is detailed below.

#### Execution flow examples for the gosbom library

```mermaid
sequenceDiagram
    participant source as source.New(ubuntu:latest)
    participant sbom as sbom.SBOM
    participant catalog as gosbom.CatalogPackages(src)
    participant encoder as gosbom.Encode(sbom, format)

    Note right of source: use "ubuntu:latest" as SBOM input

    source-->>+sbom: add source to SBOM struct
    source-->>+catalog: pass src to generate catalog
    catalog-->-sbom: add cataloging results onto SBOM
    sbom-->>encoder: pass SBOM and format desiered to gosbom encoder
    encoder-->>source: return bytes that are the SBOM of the original input 

    Note right of catalog: cataloger configuration is done based on src
```

### Supported sources

Gosbom can generate a SBOM from a variety of sources:

```
# catalog a container image archive (from the result of `docker image save ...`, `podman save ...`, or `skopeo copy` commands)
gosbom path/to/image.tar

# catalog a Singularity Image Format (SIF) container
gosbom path/to/image.sif

# catalog a directory
gosbom path/to/dir
```

Sources can be explicitly provided with a scheme:

```
docker:yourrepo/yourimage:tag            use images from the Docker daemon
podman:yourrepo/yourimage:tag            use images from the Podman daemon
docker-archive:path/to/yourimage.tar     use a tarball from disk for archives created from "docker save"
oci-archive:path/to/yourimage.tar        use a tarball from disk for OCI archives (from Skopeo or otherwise)
oci-dir:path/to/yourimage                read directly from a path on disk for OCI layout directories (from Skopeo or otherwise)
singularity:path/to/yourimage.sif        read directly from a Singularity Image Format (SIF) container on disk
dir:path/to/yourproject                  read directly from a path on disk (any directory)
file:path/to/yourproject/file            read directly from a path on disk (any single file)
registry:yourrepo/yourimage:tag          pull image directly from a registry (no container runtime required)
```

If an image source is not provided and cannot be detected from the given reference it is assumed the image should be pulled from the Docker daemon.
If docker is not present, then the Podman daemon is attempted next, followed by reaching out directly to the image registry last.


This default behavior can be overridden with the `default-image-pull-source` configuration option (See [Configuration](https://github.com/nextlinux/gosbom#configuration) for more details).

### Default Cataloger Configuration by scan type

##### Image Scanning:
- alpmdb
- rpmdb
- dpkgdb
- apkdb
- portage
- ruby-gemspec
- python-package
- php-composer-installed Cataloger
- javascript-package
- java
- go-module-binary
- dotnet-deps

##### Directory Scanning:
- alpmdb
- apkdb
- dpkgdb
- portage
- rpmdb
- ruby-gemfile
- python-index
- python-package
- php-composer-lock
- javascript-lock
- java
- java-pom
- go-module-binary
- go-mod-file
- rust-cargo-lock
- dartlang-lock
- dotnet-deps
- cocoapods
- conan
- hackage

##### Non Default:
- cargo-auditable-binary

### Excluding file paths

Gosbom can exclude files and paths from being scanned within a source by using glob expressions
with one or more `--exclude` parameters:
```
gosbom <source> --exclude './out/**/*.json' --exclude /etc
```
**Note:** in the case of _image scanning_, since the entire filesystem is scanned it is
possible to use absolute paths like `/etc` or `/usr/**/*.txt` whereas _directory scans_
exclude files _relative to the specified directory_. For example: scanning `/usr/foo` with
`--exclude ./package.json` would exclude `/usr/foo/package.json` and `--exclude '**/package.json'`
would exclude all `package.json` files under `/usr/foo`. For _directory scans_,
it is required to begin path expressions with `./`, `*/`, or `**/`, all of which
will be resolved _relative to the specified scan directory_. Keep in mind, your shell
may attempt to expand wildcards, so put those parameters in single quotes, like:
`'**/*.json'`.

### Output formats

The output format for Gosbom is configurable as well using the
`-o` (or `--output`) option:

```
gosbom <image> -o <format>
```

Where the `formats` available are:
- `json`: Use this to get as much information out of Gosbom as possible!
- `text`: A row-oriented, human-and-machine-friendly output.
- `cyclonedx-xml`: A XML report conforming to the [CycloneDX 1.4 specification](https://cyclonedx.org/specification/overview/).
- `cyclonedx-json`: A JSON report conforming to the [CycloneDX 1.4 specification](https://cyclonedx.org/specification/overview/).
- `spdx-tag-value`: A tag-value formatted report conforming to the [SPDX 2.3 specification](https://spdx.github.io/spdx-spec/v2.3/).
- `spdx-tag-value@2.2`: A tag-value formatted report conforming to the [SPDX 2.2 specification](https://spdx.github.io/spdx-spec/v2.2.2/).
- `spdx-json`: A JSON report conforming to the [SPDX 2.3 JSON Schema](https://github.com/spdx/spdx-spec/blob/v2.3/schemas/spdx-schema.json).
- `spdx-json@2.2`: A JSON report conforming to the [SPDX 2.2 JSON Schema](https://github.com/spdx/spdx-spec/blob/v2.2/schemas/spdx-schema.json).
- `github`: A JSON report conforming to GitHub's dependency snapshot format.
- `table`: A columnar summary (default).
- `template`: Lets the user specify the output format. See ["Using templates"](#using-templates) below.

## Using templates

Gosbom lets you define custom output formats, using [Go templates](https://pkg.go.dev/text/template). Here's how it works:

- Define your format as a Go template, and save this template as a file.

- Set the output format to "template" (`-o template`).

- Specify the path to the template file (`-t ./path/to/custom.template`).

- Gosbom's template processing uses the same data models as the `json` output format — so if you're wondering what data is available as you author a template, you can use the output from `gosbom <image> -o json` as a reference.

**Example:** You could make Gosbom output data in CSV format by writing a Go template that renders CSV data and then running `gosbom <image> -o template -t ~/path/to/csv.tmpl`.

Here's what the `csv.tmpl` file might look like:
```gotemplate
"Package","Version Installed","Found by"
{{- range .Artifacts}}
"{{.Name}}","{{.Version}}","{{.FoundBy}}"
{{- end}}
```

Which would produce output like:
```text
"Package","Version Installed","Found by"
"alpine-baselayout","3.2.0-r20","apkdb-cataloger"
"alpine-baselayout-data","3.2.0-r20","apkdb-cataloger"
"alpine-keys","2.4-r1","apkdb-cataloger"
...
```

Gosbom also includes a vast array of utility templating functions from [sprig](http://masterminds.github.io/sprig/) apart from the default Golang [text/template](https://pkg.go.dev/text/template#hdr-Functions) to allow users to customize the output format.

Lastly, Gosbom has custom templating functions defined in `./gosbom/format/template/encoder.go` to help parse the passed-in JSON structs.

## Multiple outputs

Gosbom can also output _multiple_ files in differing formats by appending
`=<file>` to the option, for example to output Gosbom JSON and SPDX JSON:

```shell
gosbom <image> -o json=sbom.gosbom.json -o spdx-json=sbom.spdx.json
```

## Private Registry Authentication

### Local Docker Credentials
When a container runtime is not present, Gosbom can still utilize credentials configured in common credential sources (such as `~/.docker/config.json`). It will pull images from private registries using these credentials. The config file is where your credentials are stored when authenticating with private registries via some command like `docker login`. For more information see the `go-containerregistry` [documentation](https://github.com/google/go-containerregistry/tree/main/pkg/authn).

An example `config.json` looks something like this:
```json
{
	"auths": {
		"registry.example.com": {
			"username": "AzureDiamond",
			"password": "hunter2"
		}
	}
}
```

You can run the following command as an example. It details the mount/environment configuration a container needs to access a private registry:

```
docker run -v ./config.json:/config/config.json -e "DOCKER_CONFIG=/config" nextlinux/gosbom:latest  <private_image>
```

### Docker Credentials in Kubernetes

Here's a simple workflow to mount this config file as a secret into a container on Kubernetes.

1. Create a secret. The value of `config.json` is important. It refers to the specification detailed [here](https://github.com/google/go-containerregistry/tree/main/pkg/authn#the-config-file). Below this section is the `secret.yaml` file that the pod configuration will consume as a volume. The key `config.json` is important. It will end up being the name of the file when mounted into the pod.

    ```yaml
    # secret.yaml

    apiVersion: v1
    kind: Secret
    metadata:
      name: registry-config
      namespace: gosbom
    data:
      config.json: <base64 encoded config.json>
    ```

   `kubectl apply -f secret.yaml`


2. Create your pod running gosbom. The env `DOCKER_CONFIG` is important because it advertises where to look for the credential file. In the below example, setting `DOCKER_CONFIG=/config` informs gosbom that credentials can be found at `/config/config.json`. This is why we used `config.json` as the key for our secret. When mounted into containers the secrets' key is used as the filename. The `volumeMounts` section mounts our secret to `/config`. The `volumes` section names our volume and leverages the secret we created in step one.

    ```yaml
    # pod.yaml

    apiVersion: v1
    kind: Pod
    metadata:
      name: gosbom-k8s-usage
    spec:
      containers:
        - image: nextlinux/gosbom:latest
          name: gosbom-private-registry-demo
          env:
            - name: DOCKER_CONFIG
              value: /config
          volumeMounts:
          - mountPath: /config
            name: registry-config
            readOnly: true
          args:
            - <private_image>
      volumes:
      - name: registry-config
        secret:
          secretName: registry-config
    ```

   `kubectl apply -f pod.yaml`


3. The user can now run `kubectl logs gosbom-private-registry-demo`. The logs should show the Gosbom analysis for the `<private_image>` provided in the pod configuration.

Using the above information, users should be able to configure private registry access without having to do so in the `grype` or `gosbom` configuration files.  They will also not be dependent on a Docker daemon, (or some other runtime software) for registry configuration and access.

## Format conversion (experimental)

The ability to convert existing SBOMs means you can create SBOMs in different formats quickly, without the need to regenerate the SBOM from scratch, which may take significantly more time.

```
gosbom convert <ORIGINAL-SBOM-FILE> -o <NEW-SBOM-FORMAT>[=<NEW-SBOM-FILE>]
```

This feature is experimental and data might be lost when converting formats. Packages are the main SBOM component easily transferable across formats, whereas files and relationships, as well as other information Gosbom doesn't support, are more likely to be lost.

We support formats with wide community usage AND good encode/decode support by Gosbom. The supported formats are:
- Gosbom JSON
- SPDX 2.2 JSON
- SPDX 2.2 tag-value
- CycloneDX 1.4 JSON
- CycloneDX 1.4 XML

Conversion example:
```sh
gosbom alpine:latest -o gosbom-json=sbom.gosbom.json # generate a gosbom SBOM
gosbom convert sbom.gosbom.json -o cyclonedx-json=sbom.cdx.json  # convert it to CycloneDX
```

## Attestation (experimental)
### Keyless support
Gosbom supports generating attestations using cosign's [keyless](https://github.com/sigstore/cosign/blob/main/KEYLESS.md) signatures.

Note: users need to have >= v1.12.0 of cosign installed for this command to function

To use this feature with a format like CycloneDX json simply run:
```
gosbom attest --output cyclonedx-json <IMAGE WITH OCI WRITE ACCESS>
```
This command will open a web browser and allow the user to authenticate their OIDC identity as the root of trust for the attestation (Github, Google, Microsoft).

After authenticating, Gosbom will upload the attestation to the OCI registry specified by the image that the user has write access to.

You will need to make sure your credentials are configured for the OCI registry you are uploading to so that the attestation can write successfully.

Users can then verify the attestation(or any image with attestations) by running:
```
COSIGN_EXPERIMENTAL=1 cosign verify-attestation <IMAGE_WITH_ATTESTATIONS>
```

Users should see that the uploaded attestation claims are validated, the claims exist within the transparency log, and certificates on the attestations were verified against [fulcio](https://github.com/SigStore/fulcio).
There will also be a printout of the certificates subject `<user identity>` and the certificate issuer URL: `<provider of user identity (Github, Google, Microsoft)>`:
```
Certificate subject:  test.email@testdomain.com
Certificate issuer URL:  https://accounts.google.com
```

### Local private key support

To generate an SBOM attestation for a container image using a local private key:
```
gosbom attest --output [FORMAT] --key [KEY] [SOURCE] [flags]
```

The above output is in the form of the [DSSE envelope](https://github.com/secure-systems-lab/dsse/blob/master/envelope.md#dsse-envelope).
The payload is a base64 encoded `in-toto` statement with the generated SBOM as the predicate. For details on workflows using this command see [here](#adding-an-sbom-to-an-image-as-an-attestation-using-gosbom).



## Configuration

Configuration search paths:

- `.gosbom.yaml`
- `.gosbom/config.yaml`
- `~/.gosbom.yaml`
- `<XDG_CONFIG_HOME>/gosbom/config.yaml`

Configuration options (example values are the default):

```yaml
# the output format(s) of the SBOM report (options: table, text, json, spdx, ...)
# same as -o, --output, and GOSBOM_OUTPUT env var
# to specify multiple output files in differing formats, use a list:
# output:
#   - "json=<gosbom-json-output-file>"
#   - "spdx-json=<spdx-json-output-file>"
output: "table"

# suppress all output (except for the SBOM report)
# same as -q ; GOSBOM_QUIET env var
quiet: false

# same as --file; write output report to a file (default is to write to stdout)
file: ""

# enable/disable checking for application updates on startup
# same as GOSBOM_CHECK_FOR_APP_UPDATE env var
check-for-app-update: true

# allows users to specify which image source should be used to generate the sbom
# valid values are: registry, docker, podman
default-image-pull-source: ""

# a list of globs to exclude from scanning. same as --exclude ; for example:
# exclude:
#   - "/etc/**"
#   - "./out/**/*.json"
exclude: []

# os and/or architecture to use when referencing container images (e.g. "windows/armv6" or "arm64")
# same as --platform; GOSBOM_PLATFORM env var
platform: ""

# set the list of package catalogers to use when generating the SBOM
# default = empty (cataloger set determined automatically by the source type [image or file/directory])
# catalogers:
#   - ruby-gemfile
#   - ruby-gemspec
#   - python-index
#   - python-package
#   - javascript-lock
#   - javascript-package
#   - php-composer-installed
#   - php-composer-lock
#   - alpmdb
#   - dpkgdb
#   - rpmdb
#   - java
#   - apkdb
#   - go-module-binary
#   - go-mod-file
#   - dartlang-lock
#   - rust
#   - dotnet-deps
# rust-audit-binary scans Rust binaries built with https://github.com/Shnatsel/rust-audit
#   - rust-audit-binary
catalogers:

# cataloging packages is exposed through the packages and power-user subcommands
package:

  # search within archives that do contain a file index to search against (zip)
  # note: for now this only applies to the java package cataloger
  # GOSBOM_PACKAGE_SEARCH_INDEXED_ARCHIVES env var
  search-indexed-archives: true

  # search within archives that do not contain a file index to search against (tar, tar.gz, tar.bz2, etc)
  # note: enabling this may result in a performance impact since all discovered compressed tars will be decompressed
  # note: for now this only applies to the java package cataloger
  # GOSBOM_PACKAGE_SEARCH_UNINDEXED_ARCHIVES env var
  search-unindexed-archives: false

  cataloger:
    # enable/disable cataloging of packages
    # GOSBOM_PACKAGE_CATALOGER_ENABLED env var
    enabled: true

    # the search space to look for packages (options: all-layers, squashed)
    # same as -s ; GOSBOM_PACKAGE_CATALOGER_SCOPE env var
    scope: "squashed"

golang:
   # search for go package licences in the GOPATH of the system running Gosbom, note that this is outside the
   # container filesystem and potentially outside the root of a local directory scan
   # GOSBOM_GOLANG_SEARCH_LOCAL_MOD_CACHE_LICENSES env var
   search-local-mod-cache-licenses: false
   
   # specify an explicit go mod cache directory, if unset this defaults to $GOPATH/pkg/mod or $HOME/go/pkg/mod
   # GOSBOM_GOLANG_LOCAL_MOD_CACHE_DIR env var
   local-mod-cache-dir: ""

   # search for go package licences by retrieving the package from a network proxy
   # GOSBOM_GOLANG_SEARCH_REMOTE_LICENSES env var
   search-remote-licenses: false

   # remote proxy to use when retrieving go packages from the network,
   # if unset this defaults to $GOPROXY followed by https://proxy.golang.org
   # GOSBOM_GOLANG_PROXY env var
   proxy: ""

   # specifies packages which should not be fetched by proxy
   # if unset this defaults to $GONOPROXY
   # GOSBOM_GOLANG_NOPROXY env var
   no-proxy: ""

linux-kernel:
   # whether to catalog linux kernel modules found within lib/modules/** directories
   # GOSBOM_LINUX_KERNEL_CATALOG_MODULES env var
   catalog-modules: true

# cataloging file contents is exposed through the power-user subcommand
file-contents:
  cataloger:
    # enable/disable cataloging of secrets
    # GOSBOM_FILE_CONTENTS_CATALOGER_ENABLED env var
    enabled: true

    # the search space to look for secrets (options: all-layers, squashed)
    # GOSBOM_FILE_CONTENTS_CATALOGER_SCOPE env var
    scope: "squashed"

  # skip searching a file entirely if it is above the given size (default = 1MB; unit = bytes)
  # GOSBOM_FILE_CONTENTS_SKIP_FILES_ABOVE_SIZE env var
  skip-files-above-size: 1048576

  # file globs for the cataloger to match on
  # GOSBOM_FILE_CONTENTS_GLOBS env var
  globs: []

# cataloging file metadata is exposed through the power-user subcommand
file-metadata:
  cataloger:
    # enable/disable cataloging of file metadata
    # GOSBOM_FILE_METADATA_CATALOGER_ENABLED env var
    enabled: true

    # the search space to look for file metadata (options: all-layers, squashed)
    # GOSBOM_FILE_METADATA_CATALOGER_SCOPE env var
    scope: "squashed"

  # the file digest algorithms to use when cataloging files (options: "sha256", "md5", "sha1")
  # GOSBOM_FILE_METADATA_DIGESTS env var
  digests: ["sha256"]

# maximum number of workers used to process the list of package catalogers in parallel
parallelism: 1

# cataloging secrets is exposed through the power-user subcommand
secrets:
  cataloger:
    # enable/disable cataloging of secrets
    # GOSBOM_SECRETS_CATALOGER_ENABLED env var
    enabled: true

    # the search space to look for secrets (options: all-layers, squashed)
    # GOSBOM_SECRETS_CATALOGER_SCOPE env var
    scope: "all-layers"

  # show extracted secret values in the final JSON report
  # GOSBOM_SECRETS_REVEAL_VALUES env var
  reveal-values: false

  # skip searching a file entirely if it is above the given size (default = 1MB; unit = bytes)
  # GOSBOM_SECRETS_SKIP_FILES_ABOVE_SIZE env var
  skip-files-above-size: 1048576

  # name-regex pairs to consider when searching files for secrets. Note: the regex must match single line patterns
  # but may also have OPTIONAL multiline capture groups. Regexes with a named capture group of "value" will
  # use the entire regex to match, but the secret value will be assumed to be entirely contained within the
  # "value" named capture group.
  additional-patterns: {}

  # names to exclude from the secrets search, valid values are: "aws-access-key", "aws-secret-key", "pem-private-key",
  # "docker-config-auth", and "generic-api-key". Note: this does not consider any names introduced in the
  # "secrets.additional-patterns" config option.
  # GOSBOM_SECRETS_EXCLUDE_PATTERN_NAMES env var
  exclude-pattern-names: []

# options when pulling directly from a registry via the "registry:" scheme
registry:
  # skip TLS verification when communicating with the registry
  # GOSBOM_REGISTRY_INSECURE_SKIP_TLS_VERIFY env var
  insecure-skip-tls-verify: false
  # use http instead of https when connecting to the registry
  # GOSBOM_REGISTRY_INSECURE_USE_HTTP env var
  insecure-use-http: false

  # credentials for specific registries
  auth:
      # the URL to the registry (e.g. "docker.io", "localhost:5000", etc.)
      # GOSBOM_REGISTRY_AUTH_AUTHORITY env var
    - authority: ""
      # GOSBOM_REGISTRY_AUTH_USERNAME env var
      username: ""
      # GOSBOM_REGISTRY_AUTH_PASSWORD env var
      password: ""
      # note: token and username/password are mutually exclusive
      # GOSBOM_REGISTRY_AUTH_TOKEN env var
      token: ""
      # - ... # note, more credentials can be provided via config file only

# generate an attested SBOM
attest:
  # path to the private key file to use for attestation
  # GOSBOM_ATTEST_KEY env var
  key: "cosign.key"

  # password to decrypt to given private key
  # GOSBOM_ATTEST_PASSWORD env var, additionally responds to COSIGN_PASSWORD
  password: ""

log:
  # use structured logging
  # same as GOSBOM_LOG_STRUCTURED env var
  structured: false

  # the log level; note: detailed logging suppress the ETUI
  # same as GOSBOM_LOG_LEVEL env var
  level: "error"

  # location to write the log file (default is not to have a log file)
  # same as GOSBOM_LOG_FILE env var
  file: ""
```

### Adding an SBOM to an image as an attestation using Gosbom

`gosbom attest --output [FORMAT] --key [KEY] [SOURCE] [flags]`

SBOMs themselves can serve as input to different analysis tools. [Grype](https://github.com/nextlinux/grype), a vulnerability scanner CLI tool from Nextlinux, is one such tool. Publishers of container images can use attestations to enable their consumers to trust Gosbom-generated SBOM descriptions of those container images. To create and provide these attestations, image publishers can run `gosbom attest` in conjunction with the [cosign](https://github.com/sigstore/cosign) tool to attach SBOM attestations to their images.

#### Example attestation
Note for the following example replace `docker.io/image:latest` with an image you own. You should also have push access to
its remote reference. Replace `$MY_PRIVATE_KEY` with a private key you own or have generated with cosign.

```bash
gosbom attest --key $MY_PRIVATE_KEY -o spdx-json docker.io/image:latest > image_latest_sbom_attestation.json
cosign attach attestation --attestation image_latest_sbom_attestation.json docker.io/image:latest
```

Verify the new attestation exists on your image.

```bash
cosign verify-attestation --key $MY_PUBLIC_KEY --type spdxjson docker.io/image:latest | jq '.payload | @base64d | .payload | fromjson | .predicate'
```

You should see this output along with the attached SBOM:

```
Verification for docker.io/image:latest --
The following checks were performed on each of these signatures:
  - The cosign claims were validated
  - The signatures were verified against the specified public key
  - Any certificates were verified against the Fulcio roots.
```

Consumers of your image can now trust that the SBOM associated with your image is correct and from a trusted source.

The SBOM can be piped to Grype:


```bash
cosign verify-attestation --key $MY_PUBLIC_KEY --type spdxjson docker.io/image:latest | jq '.payload | @base64d | .payload | fromjson | .predicate' | grype
```
