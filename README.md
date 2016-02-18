Pared down example of how to do version injection and methods for golang.  Eventually will be looking at making available the whole `Makefile` and process which performs additional things useful for bootstrapping golang projects.

## Quick steps

Repository needs to have a tag on it.

Using this repository as an example.

```
root@3e32bcca8e8e:/go/src/github.com/mhoglan/golang_version_injection# git clone https://github.com/mhoglan/golang_version_injection.git .
root@3e32bcca8e8e:/go/src/github.com/mhoglan/golang_version_injection# ls
LICENSE  Makefile  README.md  generate_version_info.sh	main.go  scripts  textfile_constants.go
```

`generate_version_info.sh` script can output `keyvalue` pairs of the workspace, useful for sourcing in the `Makefile`

```
root@3e32bcca8e8e:/go/src/github.com/mhoglan/golang_version_injection# ./generate_version_info.sh keyvalue
BRANCH=master
BUILD_DATE=20160218.164903
BUILD_LABEL=projectname-v0.0.1-0-ga
COMMITS=0
DIRTY=false
GIT_DESCRIBE=v0.0.1-0-g0034474
GIT_SHA1=g0034474
LABEL=ga
VERSION=v0.0.1
VERSION_INFO_JSON='{ "version_info": { "branch": "master", "build_date": "20160218.164903", "build_label": "projectname-v0.0.1-0-ga", "commits": "0", "dirty": "false", "git_describe": "v0.0.1-0-g0034474", "git_sha1": "g0034474", "label": "ga", "version": "v0.0.1" } }'
```

Or output `json` which is useful for human readable or machine parseable

```
root@3e32bcca8e8e:/go/src/github.com/mhoglan/golang_version_injection# ./generate_version_info.sh json
{
    "version_info": {
        "branch": "master",
        "build_date": "20160218.164907",
        "build_label": "projectname-v0.0.1-0-ga",
        "commits": "0",
        "dirty": "false",
        "git_describe": "v0.0.1-0-g0034474",
        "git_sha1": "g0034474",
        "label": "ga",
        "version": "v0.0.1"
    }
}
```

The `Makefile` wraps the `go generate` and `go install` commands

Example of injecting a value via the `-ldflags` is also shown

```
root@3e32bcca8e8e:/go/src/github.com/mhoglan/golang_version_injection# make go-install
go generate github.com/mhoglan/golang_version_injection/...
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL=projectname-v0.0.1-0-ga" github.com/mhoglan/golang_version_injection/...
```

Check the resultant binary and can see the `json` version information has been injected

```
root@3e32bcca8e8e:/go/src/github.com/mhoglan/golang_version_injection# /go/bin/golang_version_injection version
{
  "version_info": {
    "branch": "master",
    "build_date": "20160218.164912",
    "build_label": "projectname-v0.0.1-0-ga",
    "commits": "0",
    "dirty": "false",
    "git_describe": "v0.0.1-0-g0034474",
    "git_sha1": "g0034474",
    "label": "ga",
    "version": "v0.0.1"
  }
}
```

## Version Information
Information regarding the environment and version of the project can be injected into the golang binary.  The binary can then display this information upon being asked.  This is preferable over relying on secondary files packaged with the build, timestamps or md5sums of the binary  

There are several different ways that information can be injected into a golang project.  Will first describe the way that is being done with this project, and then alternative ways that may be useful in other situations.

These strategies are agnostic of the the project and will work in a vast majority of GIT based projects.  The injection described is golang specific, but same philosophy can be employed with other languages.

By having the detailed build information available from the binary, health checks can be setup on servers to check if services are based on clean builds or not.  Then actions can be taken for how to get a server back to a clean deploy state to satisfy the health check.

### Generation
The script `generate_version_info.sh` is responsible for generating the output to be included in builds of the project.  It is preferable to have a separate script generating this information instead of doing it directly in the `Makefile`.  This allows projects to alter just the script and have more complex logic that is not possible inside the `Makefile`.

This script can be called with `keyvalue` or `json` parameter to output different formats.   

```
~ ❯ ./generate_version_info.sh keyvalue
BRANCH=master
BUILD_DATE=20150804.194521
BUILD_LABEL=project_name-v1.0.2-0-ga-dirty
COMMITS=0
DIRTY=true
GIT_DESCRIBE=v1.0.2-0-g0e28d10-dirty
GIT_SHA1=g0e28d10
LABEL=ga
VERSION=v1.0.2
VERSION_INFO_JSON='{ "version_info": { "branch": "master", "build_date": "20150804.194521", "build_label": "project_name-v1.0.2-0-ga-dirty", "commits": "0", "dirty": "true", "git_describe": "v1.0.2-0-g0e28d10-dirty", "git_sha1": "g0e28d10", "label": "ga", "version": "v1.0.2" } }'
~ ❯ ./generate_version_info.sh json
{
    "version_info": {
        "branch": "master",
        "build_date": "20150804.194524",
        "build_label": "project_name-v1.0.2-0-ga-dirty",
        "commits": "0",
        "dirty": "true",
        "git_describe": "v1.0.2-0-g0e28d10-dirty",
        "git_sha1": "g0e28d10",
        "label": "ga",
        "version": "v1.0.2"
    }
}
```

`keyvalue` is useful for sourcing the output in other scripts, such as the `Makefile`.  This makes the variables available to be used directly. 

`json` is useful for providing a structured output that can be machine interpreted.

The script has been made to where any ENV vars that are defined with a prefix of `EXPORT_VAR_` will be included in the output.  This allows callers of the script to include information dynamically, such as a build server wanting to include additional information of the build environment.

```
~ ❯ EXPORT_VAR_BUILD_SERVER=build-server-1 ./generate_version_info.sh json
{
    "version_info": {
        "branch": "master",
        "build_date": "20150804.194949",
        "build_label": "project_name-v1.0.2-0-ga-dirty",
        "build_server": "build-server-1",
        "commits": "0",
        "dirty": "true",
        "git_describe": "v1.0.2-0-g0e28d10-dirty",
        "git_sha1": "g0e28d10",
        "label": "ga",
        "version": "v1.0.2"
    }
}
```

#### GIT Information
Various details of the GIT workspace are used to determine the version information and state of the build.  With these we can generate a unique version string that describes the build.  We can also include enough information to know if the build is reproducible and how to reproduce.

By utilizing `git tag` we can avoid having to manage a version file manually and eliminate the merge conflicts that result from doing so.   The version will be determined at build time which is a better source of truth.

`git describe` can give us the following information

* `version`
 * the most recent `git tag` in the commit tree of the workspace (what version this workspace is derived from)
* `commits` 
 * how many commits ahead is the workspace from the `version` 
* `sha1`
 * the latest commit hash of the workspace 
 * it is not enough to know just the `commits` as different branches could be ahead by the same number of `commits` 
* `dirty`
 * whether the workspace contains any staged files which have not been committed
 * does not include whether there are unknown (untracked) files in the workspace 

This is close to enough information to identify the workspace state.   

Untracked files can change the build output of a project.  The `dirty` flag can be updated to also include if there are any untracked files by checking `git ls-files`

```
if [[ $(git ls-files --directory --exclude-standard --others -t) ]]; then
    DIRTY=true
fi
```

The branch that was checked out can be found with `git rev-parse`

```
BRANCH=$(git rev-parse --abbrev-ref HEAD)
``` 

A `label` value can be generated based on `branch` to allow identifying the type of build.  

Assign the `sha1` as the `label` to any build with an unexpected `branch` value.  

Separating the `label` from the `branch` like this allows marketing / facing designations to be used for builds, and should the logic need to change for how a particular type of build is made, it is just a matter of tagging the build with the right `label` value.  

Simple mapping for now.

```
# Generate a well known label based on branch
# Unknown branches will go by their SHA1
case $BRANCH in
  master)
    LABEL="ga"
    ;;
  develop)
    LABEL="dev"
    ;;
  release/*)
    LABEL="rc"
    ;;
  *)
    LABEL=$SHA1
esac
```

Taking all of this, a build artifact name can be generated

```
{version}-{commits}-{label}[-dirty]
# examples
# clean workspace
v1.0.2-0-ga
# dirty workspace
v1.0.2-0-ga-dirty
```

Ideally all deployed builds would be clean.  Checking the artifact name and the information generated by `generate_version_info.sh` will quickly identify if so, and how to reproduce the build if possible.

With clean builds, the `commits` value can be treated as a build number / revision.  After all, the goal is to build every commit right?

Should be noted that not all builds are reproducible, if the workspace was dirty, those files could be lost as they were never committed or tracked.  Thus important to know if the build came from a dirty workspace.

### Injection
Capturing the version information is a good first step, and allows proper naming of artifacts and could even dump the information to a file at build time to be included in the packaging of the artifact.  

However artifacts can be renamed and secondary files can be changed or lost in distribution.   Binaries can be replaced by development builds, and the secondary files around it not changed to reflect the state.   

Also in a golang world, having to add in packaging to just include a resource file with the binary puts a burden on simple projects which produce a single static compiled binary.

The version information ideally would come from the binary itself.   

Included in the project is a utility script `./scripts/includetxt.go` which will convert a text file to a golang source file `textfile_constants.go` with the text defined as a constant.  The constant can be used by the golang project and will be included in the binary output.

By calling the `generate_version_info.sh` script and dumping the result to a `version_info` file, the utility script can be used by a `go generate` command to convert the `version_info` file and make the JSON string available as a constant.

The `go generate` command specified in `main.go`

```
//go:generate go run scripts/includetxt.go version_info
```

Now in the `main.go` we can define if the `version` parameter is passed in, the JSON constant can be outputted

```
case "version":
    buf := new(bytes.Buffer)
    json.Indent(buf, []byte(version_info), "", "  ")
    fmt.Println(buf)
```

From the binary output

```
~ ❯ ./target/project_name version
{
  "version_info": {
    "branch": "master",
    "build_date": "20150804.183531",
    "build_label": "project_name-v1.0.2-0-ga-dirty",
    "commits": "0",
    "dirty": "true",
    "git_describe": "v1.0.2-0-g0e28d10-dirty",
    "git_sha1": "g0e28d10",
    "label": "ga",
    "version": "v1.0.2"
  }
}
```

### Alternative Injection Methods
The above injection strategy involved using `go generate` to act as a preprocessor command to convert a text file to a constant that will be included in the source.  

There are other ways that information can be injected into a golang binary at build time.  Some of which the `Makefile` included allow if wanted.  Most strategies revolve around using `-ldflags` or transformation into source files.

#### LDFLAGS (static values)
One of the first ways seen to inject values is by utilizing the `-ldflags -X symbol value` construct to pass in key value pairs.  There is a varying level of intricacy that be done using this.

Example of the symbol `main.BUILD_LABEL` being set to `project_name-v1.0.2-0-ga-dirty`

```
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL project_name-v1.0.2-0-ga-dirty" github.com/mhoglan/project_name/...
```

Multiple key value pairs can be passed in by repeating the `-X symbol value` construct in the `-ldflags` definition

```
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL project_name-v1.0.2-0-ga-dirty -X main.TEST test" github.com/mhoglan/project_name/...
```

Symbol names map to `importpath.name` and can be displayed with `go tool nm`

#### LDFLAGS (dynamic values)
Using static definitions of symbols to be injected works well in simple scenarios.  To introduce a new symbol to be injected requires multiple places to be updated though, and as the number of symbols grows, the same logic is repeated over and over.  Adding in some intelligence into the flow can allow this to be more dynamic, working in a coding by convention style.
 
The included `Makefile` performs dynamic injection for all variables (ENV and local) defined that begin with the prefix `GO_VAR_`.  Combining this with the `keyvalue` option of the `generate_version_info.sh` script, and information of the build environment can be injected.

```
# Vars for go phase
# All vars which being with prefix will be included in ldflags
# Defaulting to full static build
GO_VARIABLE_PREFIX              = GO_VAR_
GO_VAR_BUILD_LABEL              := $(BUILD_LABEL)
GO_LDFLAGS                      = $(foreach v,$(filter $(GO_VARIABLE_PREFIX)%, $(.VARIABLES)),-X main.$(patsubst $(GO_VARIABLE_PREFIX)%,%,$(v)) $(value $(value v)))
GO_BUILD_FLAGS                  = -a -tags netgo -installsuffix nocgo -ldflags "$(GO_LDFLAGS)"
```

Callers to the `Makefile` can override / include additional values to be injected.  Must use the `-e` parameter to `make` to have ENV vars override defined vars in the `Makefile` or the vars can be passed in the override position with the targets.

```
# will not work because ENV vars do not override by default and GO_VAR_BUILD_LABEL already defined
~ ❯ GO_VAR_BUILD_LABEL=test make go-install
go generate github.com/mhoglan/project_name/...
go get github.com/mhoglan/project_name/...
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL project_name-v1.0.2-0-ga-dirty" github.com/mhoglan/project_name/...

# flag to have ENV vars override
~ ❯ GO_VAR_BUILD_LABEL=test make -e go-install
go generate github.com/mhoglan/project_name/...
go get github.com/mhoglan/project_name/...
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL test" github.com/mhoglan/project_name/...

# pass vars with targets to explicitly override
~ ❯ make go-install GO_VAR_BUILD_LABEL=test
go generate github.com/mhoglan/project_name/...
go get github.com/mhoglan/project_name/...
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL test" github.com/mhoglan/project_name/...

# add additional vars to be injected
~ ❯ GO_VAR_TEST=test make go-install
go generate github.com/mhoglan/project_name/...
go get github.com/mhoglan/project_name/...
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL project_name-v1.0.2-0-ga-dirty -X main.TEST test" github.com/mhoglan/project_name/...
```

If using the container based build environments, any variables defined on the command line with `make` will be captured and passed to the subsequent `make` calls inside.

```
~ ❯ make build GO_VAR_TEST=test
docker run --rm --entrypoint /bin/sh  -v /root/project_name:/go/src/github.com/mhoglan/project_name -v /root/project_name/target:/export -w /go/src/github.com/mhoglan/project_name golang:1.4 -c "make restoredep GO_VAR_TEST=test && make go-install GO_VAR_TEST=test && cp /go/bin/project_name /export/project_name"
go get github.com/tools/godep
godep restore github.com/mhoglan/project_name/...
go generate github.com/mhoglan/project_name/...
go get github.com/mhoglan/project_name/...
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL project_name-v1.0.4-2-ga-dirty -X main.TEST test" github.com/mhoglan/project_name/...
```

Environment variables that begin with the `EXPORT_VAR_` prefix will be injected into the container environment.  Useful for having build environments inject additional information into the version information.

```
~ ❯ EXPORT_VAR_JENKINS_BUILDSERVER=buildserver-1 make build
docker run --rm --entrypoint /bin/sh -e EXPORT_VAR_JENKINS_BUILDSERVER -v /root/project_name:/go/src/github.com/mhoglan/project_name -v /root/project_name/target:/export -w /go/src/github.com/mhoglan/project_name golang:1.4 -c "make restoredep && make go-install && cp /go/bin/project_name /export/project_name"
go get github.com/tools/godep
godep restore github.com/mhoglan/project_name/...
go generate github.com/mhoglan/project_name/...
go get github.com/mhoglan/project_name/...
go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL project_name-v1.0.4-2-ga" github.com/mhoglan/project_name/...
~ ❯ ./target/project_name version
{
  "version_info": {
    "branch": "master",
    "build_date": "20150806.161238",
    "build_label": "project_name-v1.0.4-2-ga",
    "commits": "2",
    "dirty": "false",
    "git_describe": "v1.0.4-2-g9946295",
    "git_sha1": "g9946295",
    "jenkins_buildserver": "buildserver-1",
    "label": "ga",
    "version": "v1.0.4"
  }
}
```

Being able to pass in dynamic variables makes it easier to extend later.  If the golang source does not contain the variable, then it is just not used.  

As the key value pairs outputted from `generate_version_info.sh` are sourced in the `Makefile`, any project can update the output of the script to have new values injected.  Callers can also use the ENV or arguments to the `make` command to inject new values.

While the `Makefile` alleviates the need to hard code the assignments in `-ldflags`, the symbol names must match exactly as defined in the golang source (case sensitive) and there is currently not a way to define dynamic symbols in golang 1.4 (looks like maybe in golang 1.5);  

#### LDFLAGS (structured values)
In both cases of using static or dynamic definitions of `-ldflags -X`, the golang source must be updated to utilize the injected symbol.  

One way to work around this is to define a symbol which contains a string representation of a structure, such as JSON.  Golang symbols must be string values, so there is not a way to cast the value injected into a golang type directly.  It is also better to use a [supported serialization format](http://golang.org/pkg/encoding/) as that is what they are intended for, going to and from textual representations.

JSON is a good choice for structured text.  However, you have to properly escape the JSON, mainly the escaping of quotes.  This can be complicated to ensure proper JSON structure, and not interfere with the `-ldflags` and other build flags.  Both of these caveats would be true of most textual encoding formats.

Here we introduce a variable named `VERSION_INFO` that is outputted by the binary with the parameter `version_info`, and is pretty printed as JSON with the parameter `version_info_json`.

```
# main.go
case "version_info":
    fmt.Println(VERSION_INFO)
case "version_info_json":
    buf := new(bytes.Buffer)
    json.Indent(buf, []byte(VERSION_INFO), "", "  ")
    fmt.Println(buf)
```

See the differences of injected non-escaped / escaped JSON values

```
# Non escaped JSON value
/g/s/g/T/project_name ❯ go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL project_name-v1.0.2-0-ga-dirty -X main.VERSION_INFO '{ "version_info": { "branch": "master", "build_date": "20150804.213554", "build_label": "project_name-v1.0.2-0-ga-dirty", "commits": "0", "dirty": "true", "git_describe": "v1.0.2-0-g0e28d10-dirty", "git_sha1": "g0e28d10", "label": "ga", "version": "v1.0.2" } }'" github.com/mhoglan/project_name/...
# Print raw string ; invalid JSON
/g/s/g/T/project_name ❯ /go/bin/project_name version_info
{ version_info: { branch: master, build_date: 20150804.213554, build_label: project_name-v1.0.2-0-ga-dirty, commits: 0, dirty: true, git_describe: v1.0.2-0-g0e28d10-dirty, git_sha1: g0e28d10, label: ga, version: v1.0.2 } }
# Cannot pretty print invalid JSON (blank output)
/g/s/g/T/project_name ❯ /go/bin/project_name version_info_json

# Escaped JSON value
/g/s/g/T/project_name ❯ go install -a -tags netgo -installsuffix nocgo -ldflags "-X main.BUILD_LABEL project_name-v1.0.2-0-ga-dirty -X main.VERSION_INFO '{ \"version_info\": { \"branch\": \"master\", \"build_date\": \"20150804.213554\", \"build_label\": \"project_name-v1.0.2-0-ga-dirty\", \"commits\": \"0\", \"dirty\": \"true\", \"git_describe\": \"v1.0.2-0-g0e28d10-dirty\", \"git_sha1\": \"g0e28d10\", \"label\": \"ga\", \"version\": \"v1.0.2\" } }'" github.com/mhoglan/project_name/...
# Print raw string
/g/s/g/T/project_name ❯ /go/bin/project_name version_info
{ "version_info": { "branch": "master", "build_date": "20150804.213554", "build_label": "project_name-v1.0.2-0-ga-dirty", "commits": "0", "dirty": "true", "git_describe": "v1.0.2-0-g0e28d10-dirty", "git_sha1": "g0e28d10", "label": "ga", "version": "v1.0.2" } }
# Pretty print JSON
/g/s/g/T/project_name ❯ /go/bin/project_name version_info_json
{
  "version_info": {
    "branch": "master",
    "build_date": "20150804.213554",
    "build_label": "project_name-v1.0.2-0-ga-dirty",
    "commits": "0",
    "dirty": "true",
    "git_describe": "v1.0.2-0-g0e28d10-dirty",
    "git_sha1": "g0e28d10",
    "label": "ga",
    "version": "v1.0.2"
  }
}
```

If the goal is to have a JSON document injected, it seems better to use the `go generate` strategy as it is cleaner and less prone for errors in interfering with build flags and escaping values.  However, it can be useful to inject JSON values using `-ldflags -X` if desired.

#### GO GENERATE
Golang 1.4 introduced the `generate` command which allows preprocessing.  This is useful for generating code from IDLs like thrift or protocol buffers.  It is also useful for triggering golang scripts to perform encoding / transformation of resources into golang source files.   

See the strategy used in the Injection section for an example of how `go generate` is used to inject a JSON document into the golang binary.

This seems the best way to embed a JSON document.

#### Embedding Resources
When embedding a single textual resource, a strategy as described with the `go generate` to transform it works well.  When you have multiple resources, such as the assets of a web server, which can be a mixture of textual and binary objects, then it gets a bit more complicated to manage.   

There are however some utilities that have been written to help tackle this problem.  

* go-bindata
 * https://github.com/jteeuwen/go-bindata
* go-rice
 * https://github.com/GeertJohan/go.rice/
* esc
 * https://github.com/mjibson/esc

None of these were used in a proof of concept as the primary goal in this project was to inject the version information into the binary.  These utilities are aimed primarily at embedding resource folders into the binary.  


