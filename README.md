# foo-controller

following [実践入門 Kubernetesカスタムコントローラへの道](https://www.amazon.co.jp/-/en/gp/product/B0851QCR81/ref=ppx_yo_dt_b_d_asin_title_o00?ie=UTF8&psc=1)

# Version

- kubebuilder: 2.2.0
- cert-manager: v0.11.0

# 5章

## How to create

### 1. Init

```
go version
go version go1.15.5 darwin/amd64
```

```
mkdir -p $GOPATH/src/github.com/nakamasato/foo-controller-kubebuilder
```

```
os=$(go env GOOS) # darwin
arch=$(go env GOARCH) # amd64
curl -sL https://go.kubebuilder.io/dl/2.2.0/$os/$arch | tar -xz -C /tmp/
```

```
sudo mv /tmp/kubebuilder_2.2.0_darwin_amd64 /usr/local/kubebuilder
Password:
```

```
export PATH=$PATH:/usr/local/kubebuilder/bin
```

```
kubebuilder version
Version: version.Version{KubeBuilderVersion:"2.2.0", KubernetesVendor:"1.15.5", GitCommit:"0824a139f59e109c9e418a0b6e71a53c6e9e144f", BuildDate:"2019-11-20T00:00:24Z", GoOs:"unknown", GoArch:"unknown"}
```

kustomize

```
curl -o /usr/local/kubebuilder/bin/kustomize -sL https://go.kubebuilder.io/kustomize/$os/$arch
```

```
kustomize version
Version: {KustomizeVersion:v2.0.4-kubebuilder-special-edition GitCommit:f9c631e9eec7a2d6e46eb9e1bf5122f68b97d12d BuildDate:1970-01-01T00:00:00Z GoOs:darwin GoArch:amd64}
```

```
cd ~/go/src/github.com/nakamasato/foo-controller-kubebuilder
```

```
kubebuilder init --domain k8s.io
```

```
[20-11-21 13:28:37] masato-naka at ip-192-168-31-162 in ~/go/src/github.com/nakamasato/foo-controller-kubebuilder
○ tree
.
├── Dockerfile
├── Makefile
├── PROJECT
├── bin
│   └── manager
├── config
│   ├── certmanager
│   │   ├── certificate.yaml
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   ├── manager_webhook_patch.yaml
│   │   └── webhookcainjection_patch.yaml
│   ├── manager
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   ├── rbac
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── kustomization.yaml
│   │   ├── leader_election_role.yaml
│   │   ├── leader_election_role_binding.yaml
│   │   └── role_binding.yaml
│   └── webhook
│       ├── kustomization.yaml
│       ├── kustomizeconfig.yaml
│       └── service.yaml
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── main.go

9 directories, 29 files
```

```
go mod download
```

### 2. API Object and Controller

```
kubebuilder create api --group samplecontroller --version v1alpha1 --kind Foo
```

### 3. Define API Object Foo (types.go)
check the commit
### 4. Implement reconciele (controller.go)
check the commit

### 5. Edit main.go
check the commit

### 6. Execute

```
make install
go: creating new go.mod: module tmp
go: found sigs.k8s.io/controller-tools/cmd/controller-gen in sigs.k8s.io/controller-tools v0.2.4
/Users/masato-naka/go/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build config/crd | kubectl apply -f -
Warning: apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use apiextensions.k8s.io/v1 CustomResourceDefinition
customresourcedefinition.apiextensions.k8s.io/foos.samplecontroller.k8s.io created
```

Run with go

```
make run
± make run
go: creating new go.mod: module tmp
go: found sigs.k8s.io/controller-tools/cmd/controller-gen in sigs.k8s.io/controller-tools v0.2.4
/Users/masato-naka/go/bin/controller-gen object:headerFile=./hack/boilerplate.go.txt paths="./..."
go fmt ./...
go vet ./...
/Users/masato-naka/go/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go run ./main.go
2020-11-21T15:04:47.856+0900    INFO    controller-runtime.metrics      metrics server is starting to listen    {"addr": ":8080"}
2020-11-21T15:04:47.856+0900    INFO    setup   starting manager
2020-11-21T15:04:47.857+0900    INFO    controller-runtime.manager      starting metrics server {"path": "/metrics"}
2020-11-21T15:04:48.057+0900    INFO    controller-runtime.controller   Starting EventSource    {"controller": "foo", "source": "kind source: /, Kind="}
2020-11-21T15:04:48.157+0900    INFO    controller-runtime.controller   Starting EventSource    {"controller": "foo", "source": "kind source: /, Kind="}
2020-11-21T15:04:48.157+0900    INFO    controller-runtime.controller   Starting Controller     {"controller": "foo"}
2020-11-21T15:04:48.261+0900    INFO    controller-runtime.controller   Starting workers        {"controller": "foo", "worker count": 1}
2020-11-21T15:09:04.424+0900    INFO    controllers.Foo fetching Foo Resource   {"foo": "default/foo-sample"}
2020-11-21T15:09:04.424+0900    INFO    controllers.Foo finding existing Deployments for Foo resource       {"foo": "default/foo-sample"}
2020-11-21T15:09:04.468+0900    ERROR   controllers.Foo unable to fetch Deployment      {"foo": "default/foo-sample", "error": "Deployment.apps \"example-foo\" not found"}
github.com/go-logr/zapr.(*zapLogger).Error
        /Users/masato-naka/go/pkg/mod/github.com/go-logr/zapr@v0.1.0/zapr.go:128
github.com/nakamasato/foo-controller-kubebuilder/controllers.(*FooReconciler).Reconcile
        /Users/masato-naka/go/src/github.com/nakamasato/foo-controller-kubebuilder/controllers/foo_controller.go:119
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).reconcileHandler
        /Users/masato-naka/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.4.0/pkg/internal/controller/controller.go:256
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).processNextWorkItem
        /Users/masato-naka/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.4.0/pkg/internal/controller/controller.go:232
sigs.k8s.io/controller-runtime/pkg/internal/controller.(*Controller).worker
        /Users/masato-naka/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.4.0/pkg/internal/controller/controller.go:211
k8s.io/apimachinery/pkg/util/wait.JitterUntil.func1
        /Users/masato-naka/go/pkg/mod/k8s.io/apimachinery@v0.0.0-20190913080033-27d36303b655/pkg/util/wait/wait.go:152
k8s.io/apimachinery/pkg/util/wait.JitterUntil
        /Users/masato-naka/go/pkg/mod/k8s.io/apimachinery@v0.0.0-20190913080033-27d36303b655/pkg/util/wait/wait.go:153
k8s.io/apimachinery/pkg/util/wait.Until
        /Users/masato-naka/go/pkg/mod/k8s.io/apimachinery@v0.0.0-20190913080033-27d36303b655/pkg/util/wait/wait.go:88
2020-11-21T15:09:04.468+0900    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "foo", "request": "default/foo-sample"}
2020-11-21T15:09:04.471+0900    INFO    controllers.Foo fetching Foo Resource   {"foo": "default/foo-sample"}
2020-11-21T15:09:04.471+0900    INFO    controllers.Foo finding existing Deployments for Foo resource       {"foo": "default/foo-sample"}
2020-11-21T15:09:04.472+0900    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "foo", "request": "default/foo-sample"}
2020-11-21T15:09:04.528+0900    INFO    controllers.Foo fetching Foo Resource   {"foo": "default/foo-sample"}
2020-11-21T15:09:04.528+0900    INFO    controllers.Foo finding existing Deployments for Foo resource       {"foo": "default/foo-sample"}
2020-11-21T15:09:04.528+0900    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "foo", "request": "default/foo-sample"}
2020-11-21T15:09:04.548+0900    INFO    controllers.Foo fetching Foo Resource   {"foo": "default/foo-sample"}
2020-11-21T15:09:04.548+0900    INFO    controllers.Foo finding existing Deployments for Foo resource       {"foo": "default/foo-sample"}
2020-11-21T15:09:04.548+0900    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "foo", "request": "default/foo-sample"}
2020-11-21T15:09:04.624+0900    INFO    controllers.Foo fetching Foo Resource   {"foo": "default/foo-sample"}
2020-11-21T15:09:04.625+0900    INFO    controllers.Foo finding existing Deployments for Foo resource       {"foo": "default/foo-sample"}
2020-11-21T15:09:04.625+0900    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "foo", "request": "default/foo-sample"}
2020-11-21T15:09:18.027+0900    INFO    controllers.Foo fetching Foo Resource   {"foo": "default/foo-sample"}
2020-11-21T15:09:18.027+0900    INFO    controllers.Foo finding existing Deployments for Foo resource       {"foo": "default/foo-sample"}
2020-11-21T15:09:18.028+0900    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "foo", "request": "default/foo-sample"}
2020-11-21T15:09:18.066+0900    INFO    controllers.Foo fetching Foo Resource   {"foo": "default/foo-sample"}
2020-11-21T15:09:18.067+0900    INFO    controllers.Foo finding existing Deployments for Foo resource       {"foo": "default/foo-sample"}
2020-11-21T15:09:18.067+0900    DEBUG   controller-runtime.controller   Successfully Reconciled {"controller": "foo", "request": "default/foo-sample"}
```

Create resources

```
± kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
foo.samplecontroller.k8s.io/foo-sample created
```

```
kubectl get foo
NAME         AGE
foo-sample   62s
```

```
± kubectl delete -f config/samples/samplecontroller_v1alpha1_foo.yaml
foo.samplecontroller.k8s.io "foo-sample" deleted

± kubectl get foo
No resources found in default namespace.
```

```
kubectl get deploy
NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
```

Run with container

```
export IMG=nakamasato/foo-controller:kubebuilder
```

```
make docker-build
go: creating new go.mod: module tmp
go: found sigs.k8s.io/controller-tools/cmd/controller-gen in sigs.k8s.io/controller-tools v0.2.4
/Users/masato-naka/go/bin/controller-gen object:headerFile=./hack/boilerplate.go.txt paths="./..."
go fmt ./...
go vet ./...
/Users/masato-naka/go/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go test ./... -coverprofile cover.out
?       github.com/nakamasato/foo-controller-kubebuilder        [no test files]
?       github.com/nakamasato/foo-controller-kubebuilder/api/v1alpha1   [no test files]
ok      github.com/nakamasato/foo-controller-kubebuilder/controllers    7.990s  coverage: 0.0% of statements
docker build . -t nakamasato/foo-controller:kubebuilder
Sending build context to Docker daemon  39.39MB
Step 1/14 : FROM golang:1.13 as builder
1.13: Pulling from library/golang
d6ff36c9ec48: Pull complete
c958d65b3090: Pull complete
edaf0a6b092f: Pull complete
80931cf68816: Pull complete
813643441356: Pull complete
799f41bb59c9: Pull complete
16b5038bccc8: Pull complete
Digest: sha256:8ebb6d5a48deef738381b56b1d4cd33d99a5d608e0d03c5fe8dfa3f68d41a1f8
Status: Downloaded newer image for golang:1.13
 ---> d6f3656320fe
Step 2/14 : WORKDIR /workspace
 ---> Running in 24fb80e67520
Removing intermediate container 24fb80e67520
 ---> 6ae6f3f3d2df
Step 3/14 : COPY go.mod go.mod
 ---> c38ffb6f2f02
Step 4/14 : COPY go.sum go.sum
 ---> bb62172177d8
Step 5/14 : RUN go mod download
 ---> Running in a03c35a2cbf6
go: finding cloud.google.com/go v0.38.0
go: finding github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78
go: finding github.com/Azure/go-autorest/autorest v0.9.0
go: finding github.com/Azure/go-autorest/autorest/adal v0.5.0
go: finding github.com/Azure/go-autorest/autorest/date v0.1.0
go: finding github.com/Azure/go-autorest/autorest/mocks v0.2.0
go: finding github.com/Azure/go-autorest/logger v0.1.0
go: finding github.com/Azure/go-autorest/tracing v0.5.0
go: finding github.com/BurntSushi/toml v0.3.1
go: finding github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802
go: finding github.com/NYTimes/gziphandler v0.0.0-20170623195520-56545f4a5d46
go: finding github.com/PuerkitoBio/purell v1.1.1
go: finding github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578
go: finding github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6
go: finding github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
go: finding github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973
go: finding github.com/blang/semver v3.5.0+incompatible
go: finding github.com/client9/misspell v0.3.4
go: finding github.com/coreos/bbolt v1.3.1-coreos.6
go: finding github.com/coreos/etcd v3.3.15+incompatible
go: finding github.com/coreos/go-etcd v2.0.0+incompatible
go: finding github.com/coreos/go-oidc v2.1.0+incompatible
go: finding github.com/coreos/go-semver v0.3.0
go: finding github.com/coreos/go-systemd v0.0.0-20180511133405-39ca1b05acc7
go: finding github.com/coreos/pkg v0.0.0-20180108230652-97fdf19511ea
go: finding github.com/cpuguy83/go-md2man v1.0.10
go: finding github.com/davecgh/go-spew v1.1.1
go: finding github.com/dgrijalva/jwt-go v3.2.0+incompatible
go: finding github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
go: finding github.com/docker/go-units v0.3.3
go: finding github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96
go: finding github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
go: finding github.com/elazarl/goproxy v0.0.0-20170405201442-c4fc26588b6e
go: finding github.com/emicklei/go-restful v2.9.5+incompatible
go: finding github.com/evanphx/json-patch v4.5.0+incompatible
go: finding github.com/fsnotify/fsnotify v1.4.7
go: finding github.com/ghodss/yaml v0.0.0-20150909031657-73d445a93680
go: finding github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
go: finding github.com/go-logr/logr v0.1.0
go: finding github.com/go-logr/zapr v0.1.0
go: finding github.com/go-openapi/analysis v0.19.2
go: finding github.com/go-openapi/errors v0.19.2
go: finding github.com/go-openapi/jsonpointer v0.19.2
go: finding github.com/go-openapi/jsonreference v0.19.2
go: finding github.com/go-openapi/loads v0.19.2
go: finding github.com/go-openapi/runtime v0.19.0
go: finding github.com/go-openapi/spec v0.19.2
go: finding github.com/go-openapi/strfmt v0.19.0
go: finding github.com/go-openapi/swag v0.19.2
go: finding github.com/go-openapi/validate v0.19.2
go: finding github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
go: finding github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
go: finding github.com/golang/groupcache v0.0.0-20180513044358-24b0969c4cb7
go: finding github.com/golang/mock v1.2.0
go: finding github.com/golang/protobuf v1.3.2
go: finding github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c
go: finding github.com/google/go-cmp v0.3.0
go: finding github.com/google/gofuzz v1.0.0
go: finding github.com/google/martian v2.1.0+incompatible
go: finding github.com/google/pprof v0.0.0-20181206194817-3ea8567a2e57
go: finding github.com/google/uuid v1.1.1
go: finding github.com/googleapis/gax-go/v2 v2.0.4
go: finding github.com/googleapis/gnostic v0.3.1
go: finding github.com/gophercloud/gophercloud v0.1.0
go: finding github.com/gorilla/websocket v1.4.0
go: finding github.com/gregjones/httpcache v0.0.0-20170728041850-787624de3eb7
go: finding github.com/grpc-ecosystem/go-grpc-middleware v0.0.0-20190222133341-cfaf5686ec79
go: finding github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
go: finding github.com/grpc-ecosystem/grpc-gateway v1.3.0
go: finding github.com/hashicorp/golang-lru v0.5.1
go: finding github.com/hashicorp/hcl v1.0.0
go: finding github.com/hpcloud/tail v1.0.0
go: finding github.com/imdario/mergo v0.3.6
go: finding github.com/inconshreveable/mousetrap v1.0.0
go: finding github.com/jonboulle/clockwork v0.1.0
go: finding github.com/json-iterator/go v1.1.7
go: finding github.com/jstemmer/go-junit-report v0.0.0-20190106144839-af01ea7f8024
go: finding github.com/kisielk/errcheck v1.2.0
go: finding github.com/kisielk/gotool v1.0.0
go: finding github.com/konsorten/go-windows-terminal-sequences v1.0.1
go: finding github.com/kr/pretty v0.1.0
go: finding github.com/kr/pty v1.1.5
go: finding github.com/kr/text v0.1.0
go: finding github.com/magiconair/properties v1.8.0
go: finding github.com/mailru/easyjson v0.0.0-20190614124828-94de47d64c63
go: finding github.com/matttproud/golang_protobuf_extensions v1.0.1
go: finding github.com/mitchellh/go-homedir v1.1.0
go: finding github.com/mitchellh/mapstructure v1.1.2
go: finding github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
go: finding github.com/modern-go/reflect2 v1.0.1
go: finding github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d
go: finding github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f
go: finding github.com/onsi/ginkgo v1.8.0
go: finding github.com/onsi/gomega v1.5.0
go: finding github.com/pborman/uuid v1.2.0
go: finding github.com/pelletier/go-toml v1.2.0
go: finding github.com/peterbourgon/diskv v2.0.1+incompatible
go: finding github.com/pkg/errors v0.8.1
go: finding github.com/pmezard/go-difflib v1.0.0
go: finding github.com/pquerna/cachecontrol v0.0.0-20171018203845-0dec1b30a021
go: finding github.com/prometheus/client_golang v0.9.2
go: finding github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
go: finding github.com/prometheus/common v0.0.0-20181126121408-4724e9255275
go: finding github.com/prometheus/procfs v0.0.0-20181204211112-1dc9a6cbc91a
go: finding github.com/remyoudompheng/bigfft v0.0.0-20170806203942-52369c62f446
go: finding github.com/russross/blackfriday v1.5.2
go: finding github.com/sirupsen/logrus v1.4.2
go: finding github.com/soheilhy/cmux v0.1.3
go: finding github.com/spf13/afero v1.2.2
go: finding github.com/spf13/cast v1.3.0
go: finding github.com/spf13/cobra v0.0.5
go: finding github.com/spf13/jwalterweatherman v1.0.0
go: finding github.com/spf13/pflag v1.0.3
go: finding github.com/spf13/viper v1.3.2
go: finding github.com/stretchr/objx v0.2.0
go: finding github.com/stretchr/testify v1.3.0
go: finding github.com/tmc/grpc-websocket-proxy v0.0.0-20170815181823-89b8d40f7ca8
go: finding github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8
go: finding github.com/xiang90/probing v0.0.0-20160813154853-07dd2e8dfe18
go: finding github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77
go: finding go.opencensus.io v0.21.0
go: finding go.uber.org/atomic v1.3.2
go: finding go.uber.org/multierr v1.1.0
go: finding go.uber.org/zap v1.9.1
go: finding golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
go: finding golang.org/x/exp v0.0.0-20190312203227-4b39c73a6495
go: finding golang.org/x/image v0.0.0-20190227222117-0694c2d4d067
go: finding golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3
go: finding golang.org/x/mobile v0.0.0-20190312151609-d3739f865fa6
go: finding golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc
go: finding golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
go: finding golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
go: finding golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f
go: finding golang.org/x/text v0.3.2
go: finding golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
go: finding golang.org/x/tools v0.0.0-20190621195816-6e04913cbbac
go: finding golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
go: finding gomodules.xyz/jsonpatch/v2 v2.0.1
go: finding gonum.org/v1/gonum v0.0.0-20190331200053-3d26580ed485
go: finding gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e
go: finding google.golang.org/api v0.4.0
go: finding google.golang.org/appengine v1.5.0
go: finding google.golang.org/genproto v0.0.0-20190502173448-54afdca5d873
go: finding google.golang.org/grpc v1.23.0
go: finding gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
go: finding gopkg.in/fsnotify.v1 v1.4.7
go: finding gopkg.in/inf.v0 v0.9.1
go: finding gopkg.in/natefinch/lumberjack.v2 v2.0.0
go: finding gopkg.in/square/go-jose.v2 v2.2.2
go: finding gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
go: finding gopkg.in/yaml.v2 v2.2.2
go: finding gotest.tools v2.2.0+incompatible
go: finding honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc
go: finding k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
go: finding k8s.io/apiextensions-apiserver v0.0.0-20190918161926-8f644eb6e783
go: finding k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
go: finding k8s.io/apiserver v0.0.0-20190918160949-bfa5e2e684ad
go: finding k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
go: finding k8s.io/code-generator v0.0.0-20190912054826-cd179ad6a269
go: finding k8s.io/component-base v0.0.0-20190918160511-547f6c5d7090
go: finding k8s.io/gengo v0.0.0-20190822140433-26a664648505
go: finding k8s.io/klog v0.4.0
go: finding k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
go: finding k8s.io/utils v0.0.0-20190801114015-581e00157fb1
go: finding modernc.org/cc v1.0.0
go: finding modernc.org/golex v1.0.0
go: finding modernc.org/mathutil v1.0.0
go: finding modernc.org/strutil v1.0.0
go: finding modernc.org/xc v1.0.0
go: finding sigs.k8s.io/controller-runtime v0.4.0
go: finding sigs.k8s.io/structured-merge-diff v0.0.0-20190817042607-6149e4549fca
go: finding sigs.k8s.io/testing_frameworks v0.1.2
go: finding sigs.k8s.io/yaml v1.1.0
Removing intermediate container a03c35a2cbf6
 ---> a4ce8c8f4f43
Step 6/14 : COPY main.go main.go
 ---> 06c5dc0b1931
Step 7/14 : COPY api/ api/
 ---> 6c6cef61ab72
Step 8/14 : COPY controllers/ controllers/
 ---> e322aa298308
Step 9/14 : RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go
 ---> Running in 4e0372d99aa8
Removing intermediate container 4e0372d99aa8
 ---> 823ac3190381
Step 10/14 : FROM gcr.io/distroless/static:nonroot
nonroot: Pulling from distroless/static
e59bd8947ac7: Pull complete
Digest: sha256:dfabf95aa48990eec883b8b75462b1304678ebe9f36bdbca7841f5ddbd32af96
Status: Downloaded newer image for gcr.io/distroless/static:nonroot
 ---> aa99000bc55d
Step 11/14 : WORKDIR /
 ---> Running in 57d1de30c727
Removing intermediate container 57d1de30c727
 ---> 47299bbc0d99
Step 12/14 : COPY --from=builder /workspace/manager .
 ---> 387643c7b755
Step 13/14 : USER nonroot:nonroot
 ---> Running in 61e30ffa53e2
Removing intermediate container 61e30ffa53e2
 ---> be55051d8e6c
Step 14/14 : ENTRYPOINT ["/manager"]
 ---> Running in 61cbc65fdf01
Removing intermediate container 61cbc65fdf01
 ---> 533912b99d88
Successfully built 533912b99d88
Successfully tagged nakamasato/foo-controller:kubebuilder
```

Create docker repository

https://hub.docker.com/repository/docker/nakamasato/foo-controller

```
make docker-push
```

prepare config/manager/kustomization.yaml

```
make deploy
```

```
kubectl get deploy -n foo-controller-kubebuilder-system
NAME                                            READY   UP-TO-DATE   AVAILABLE   AGE
foo-controller-kubebuilder-controller-manager   1/1     1            1           60s
```

```
kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
foo.samplecontroller.k8s.io/foo-sample created
```

clean up

```
kubectl delete -f config/samples/samplecontroller_v1alpha1_foo.yaml
foo.samplecontroller.k8s.io "foo-sample" deleted

kubectl delete -k config/default/
```

# 7章

## Create Admission webhook

### Create webhook

```
kubebuilder create webhook --group samplecontroller --version v1alpha1 --kind Foo --defaulting --programmatic-validation 
Writing scaffold for you to edit...
api/v1alpha1/foo_webhook.go
```

### Implement Mutating and Validating

change `api/v1alpha1/foo_webhook.go` and `api/v1alpha1/foo_types.go`

### Execute

1. prepare certificate

```
kubectl create ns cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.11.0/cert-manager.yaml
```

```
kubectl get pod -n cert-manager                                                                      
NAME                                       READY   STATUS    RESTARTS   AGE
cert-manager-7746446996-9clhd              1/1     Running   0          63s
cert-manager-cainjector-59b88c94bc-t48p5   1/1     Running   0          64s
cert-manager-webhook-6b8f46fbfc-lk56p      1/1     Running   0          63s
```

1. Build Docker image

```
export IMG=nakamasato/foo-controller:kubebuilder-webhook
make docker-build
```

```
make docker-push
```

1. Update `config/default/kustomization.yaml`

1. Deploy

```
make deploy
```

1. Check certificate

```
kubectl get certificate -n foo-controller-kubebuilder-system
NAME                                      READY   SECRET                AGE
foo-controller-kubebuilder-serving-cert   True    webhook-server-cert   112s
```

```
± kubectl get deploy -n foo-controller-kubebuilder-system
NAME                                            READY   UP-TO-DATE   AVAILABLE   AGE
foo-controller-kubebuilder-controller-manager   1/1     1            1           2m16s

± kubectl get pod -n foo-controller-kubebuilder-system
NAME                                                             READY   STATUS    RESTARTS   AGE
foo-controller-kubebuilder-controller-manager-65dc97ddfd-lqdd2   2/2     Running   1          2m23s
```

1. Apply `config/samples/samplecontroller_v1alpha1_foo.yaml` without `replicas: 1`

```
kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml 
foo.samplecontroller.k8s.io/foo-sample created
```

1. Check replica is 1

```
kubectl get foo foo-sample -o yaml | yq r - spec.replicas
1
```

1. Apply `config/samples/samplecontroller_v1alpha1_foo.yaml` with too long name

```yaml
  deploymentName: aaaaaaaaaaaa...aaaaaaaaaaaaa # more than 253 letters
```

```
kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
Error from server (Foo.samplecontroller.k8s.io "foo-sample" is invalid: spec.deploymentName: Invalid value: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa": must be no more than 253 characters): error when applying patch:
{"metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":"{\"apiVersion\":\"samplecontroller.k8s.io/v1alpha1\",\"kind\":\"Foo\",\"metadata\":{\"annotations\":{},\"name\":\"foo-sample\",\"namespace\":\"default\"},\"spec\":{\"deploymentName\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\"}}\n"}},"spec":{"deploymentName":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
to:
Resource: "samplecontroller.k8s.io/v1alpha1, Resource=foos", GroupVersionKind: "samplecontroller.k8s.io/v1alpha1, Kind=Foo"
Name: "foo-sample", Namespace: "default"
for: "config/samples/samplecontroller_v1alpha1_foo.yaml": admission webhook "vfoo.kb.io" denied the request: Foo.samplecontroller.k8s.io "foo-sample" is invalid: spec.deploymentName: Invalid value: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa": must be no more than 253 characters
```

### clean up

```
kubectl delete -f config/samples/samplecontroller_v1alpha1_foo.yaml 
foo.samplecontroller.k8s.io "foo-sample" deleted
```

```
kubectl delete -k config/default 
namespace "foo-controller-kubebuilder-system" deleted
Warning: apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use apiextensions.k8s.io/v1 CustomResourceDefinition
customresourcedefinition.apiextensions.k8s.io "foos.samplecontroller.k8s.io" deleted
Warning: admissionregistration.k8s.io/v1beta1 MutatingWebhookConfiguration is deprecated in v1.16+, unavailable in v1.22+; use admissionregistration.k8s.io/v1 MutatingWebhookConfiguration
mutatingwebhookconfiguration.admissionregistration.k8s.io "foo-controller-kubebuilder-mutating-webhook-configuration" deleted
Warning: admissionregistration.k8s.io/v1beta1 ValidatingWebhookConfiguration is deprecated in v1.16+, unavailable in v1.22+; use admissionregistration.k8s.io/v1 ValidatingWebhookConfiguration
validatingwebhookconfiguration.admissionregistration.k8s.io "foo-controller-kubebuilder-validating-webhook-configuration" deleted
role.rbac.authorization.k8s.io "foo-controller-kubebuilder-leader-election-role" deleted
clusterrole.rbac.authorization.k8s.io "foo-controller-kubebuilder-manager-role" deleted
clusterrole.rbac.authorization.k8s.io "foo-controller-kubebuilder-proxy-role" deleted
rolebinding.rbac.authorization.k8s.io "foo-controller-kubebuilder-leader-election-rolebinding" deleted
clusterrolebinding.rbac.authorization.k8s.io "foo-controller-kubebuilder-manager-rolebinding" deleted
clusterrolebinding.rbac.authorization.k8s.io "foo-controller-kubebuilder-proxy-rolebinding" deleted
service "foo-controller-kubebuilder-controller-manager-metrics-service" deleted
service "foo-controller-kubebuilder-webhook-service" deleted
deployment.apps "foo-controller-kubebuilder-controller-manager" deleted
certificate.cert-manager.io "foo-controller-kubebuilder-serving-cert" deleted
issuer.cert-manager.io "foo-controller-kubebuilder-selfsigned-issuer" deleted
```

## Create Conversion Webhook

Conversion: API version compatibility

apps/v1alphaのResourceがApplyされても、apps/v1のResourceがApplyされても、KuberntesのControllerがそれをReconcileし、同じ機能・効果を提供するObjectを保証してくれます。これは、Kubernetesが裏側でConversionを実行することで、Multi Version間の差分を吸収してくれているからです。

### Overview

1. Add new API version
1. Implement Conversion
1. Execute the operator

### Add new API version

```
kubebuilder create api --group samplecontroller --version v1beta1 --kind Foo
Create Resource [y/n]
y
Create Controller [y/n]
n
Writing scaffold for you to edit...
api/v1beta1/foo_types.go
Running make...
go: creating new go.mod: module tmp
go: found sigs.k8s.io/controller-tools/cmd/controller-gen in sigs.k8s.io/controller-tools v0.2.4
/Users/masato-naka/go/bin/controller-gen object:headerFile=./hack/boilerplate.go.txt paths="./..."
go fmt ./...
go vet ./...
go build -o bin/manager main.go
```

### Prepare

1. Copy `v1alpha1/foo_types.go` to `v1beta1/foo_types.go`
1. Add new field for conversion

```
	// +kubebuilder:validation:Optional
	// the new field for conversion
	Foo string `json:"foo"`
```

1. Add `// +kubebuilder:storageversion` to `alphav1/foo_types.go`
1. Implement conversion -> Commit
1. Change `Makefile`
1. Build & Push

```
export IMG=nakamasato/foo-controller:kubebuilder-conversion
make docker-build
```

```
make docker-push
```

1. Deploy -> wait until controller manager starts up

```
make deploy
go: creating new go.mod: module tmp
go: found sigs.k8s.io/controller-tools/cmd/controller-gen in sigs.k8s.io/controller-tools v0.2.4
/Users/masato-naka/go/bin/controller-gen "crd:preserveUnknownFields=false" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && kustomize edit set image controller=nakamasato/foo-controller:kubebuilder-conversion
kustomize build config/default | kubectl apply -f -
namespace/foo-controller-kubebuilder-system created
Warning: apiextensions.k8s.io/v1beta1 CustomResourceDefinition is deprecated in v1.16+, unavailable in v1.22+; use apiextensions.k8s.io/v1 CustomResourceDefinition
customresourcedefinition.apiextensions.k8s.io/foos.samplecontroller.k8s.io created
Warning: admissionregistration.k8s.io/v1beta1 MutatingWebhookConfiguration is deprecated in v1.16+, unavailable in v1.22+; use admissionregistration.k8s.io/v1 MutatingWebhookConfiguration
mutatingwebhookconfiguration.admissionregistration.k8s.io/foo-controller-kubebuilder-mutating-webhook-configuration created
Warning: admissionregistration.k8s.io/v1beta1 ValidatingWebhookConfiguration is deprecated in v1.16+, unavailable in v1.22+; use admissionregistration.k8s.io/v1 ValidatingWebhookConfiguration
validatingwebhookconfiguration.admissionregistration.k8s.io/foo-controller-kubebuilder-validating-webhook-configuration created
role.rbac.authorization.k8s.io/foo-controller-kubebuilder-leader-election-role created
clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-manager-role created
clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-proxy-role created
rolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-proxy-rolebinding created
service/foo-controller-kubebuilder-controller-manager-metrics-service created
service/foo-controller-kubebuilder-webhook-service created
deployment.apps/foo-controller-kubebuilder-controller-manager created
certificate.cert-manager.io/foo-controller-kubebuilder-serving-cert created
issuer.cert-manager.io/foo-controller-kubebuilder-selfsigned-issuer created
```

1. Check

```
kubectl get deployment -n foo-controller-kubebuilder-system
```

```
kubectl get pods -n foo-controller-kubebuilder-system
```

```
kubectl logs -f -c manager -n foo-controller-kubebuilder-system
...
2020-11-28T07:51:59.198Z        INFO    controller-runtime.builder      conversion webhook enabled      {"object": {"metadata":{"creationTimestamp":null},"spec":{"deploymentName":"","replicas":null},"status":{"availableReplicas":0}}}
...
```

1. Check

Apply `betav1`

```
kubectl apply -f config/samples/samplecontroller_v1beta1_foo.yaml
foo.samplecontroller.k8s.io/foo-sample created
```

confirm `foo` is copied from alphav1

```
kubectl get foo.v1beta1.samplecontroller.k8s.io foo-sample -o yaml | yq r - spec
deploymentName: example-foo
foo: example-foo
replicas: 1
```

Check `alphav1` -> no foo

```
kubectl get foo.v1alpha1.samplecontroller.k8s.io foo-sample -o yaml | yq r - spec
deploymentName: example-foo
replicas: 1
```

### Clean up

```
kubectl delete -f config/samples/samplecontroller_v1beta1_foo.yaml 
kubectl delete -k config/default
kubectl delete -f https://github.com/jetstack/cert-manager/releases/download/v0.11.0/cert-manager.yaml
```
