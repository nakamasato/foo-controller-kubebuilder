# foo-controller

- [実践入門 Kubernetesカスタムコントローラへの道](https://www.amazon.co.jp/-/en/gp/product/B0851QCR81/ref=ppx_yo_dt_b_d_asin_title_o00?ie=UTF8&psc=1)
- https://github.com/govargo/foo-controller-kubebuilder

## Version

- go: 1.17.3
- kubebuilder: 3.2.0
- kustomize: v4.4.0
- cert-manager: v1.6.1

## Prerequisite

```
go version
go version go1.17.3 darwin/amd64
```

install kubebuilder

```bash
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
```

```
kubebuilder version
Version: main.version{KubeBuilderVersion:"3.2.0", KubernetesVendor:"1.22.1", GitCommit:"b7a730c84495122a14a0faff95e9e9615fffbfc5", BuildDate:"2021-10-29T18:32:16Z", GoOs:"darwin", GoArch:"amd64"}
```

kustomize

```
kustomize version
{Version:kustomize/v4.4.0 GitCommit:63ec6bdb3d737a7c66901828c5743656c49b60e1 BuildDate:2021-09-27T16:13:36Z GoOs:darwin GoArch:amd64}
```

## Overview

1. Init Kubebuilｄer Project
1. Create API resource `Foo` and its controller
1. Define API resource `Foo` (`types.go`)
1. Implement reconciliation logic (`controller.go` and `main.go`)
1. Run the operator

## Reconciliation logic

1. Fetch Foo object.
1. Delete old `Deployment` managed by `Foo` if exists.
1. Create `Deployment` if `Deployment` managed by `Foo` does not exists. Otherwise, update the existing one.
1. Update `Foo`'s status.

## How to create

### 1. Initialize an operator

Make directory

```bash
mkdir foo-controller-kubebuilder
cd foo-controller-kubebuilder
```

Initialize go module

```bash
go mod init github.com/nakamasato/foo-controller-kubebuilder # replace nakamasato with your GitHub name
go: creating new go.mod: module github.com/nakamasato/foo-controller-kubebuilder
```

this command generates `go.mod`

```
module github.com/nakamasato/foo-controller-kubebuilder

go 1.17
```

Initialize a kubebuilder project

```
kubebuilder init --plugins go/v3 --domain example.com --owner nakamasato
```

<details><summary>tree</summary>

```bash
tree
.
├── Dockerfile
├── Makefile
├── PROJECT
├── config
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   └── manager_config_patch.yaml
│   ├── manager
│   │   ├── controller_manager_config.yaml
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   └── rbac
│       ├── auth_proxy_client_clusterrole.yaml
│       ├── auth_proxy_role.yaml
│       ├── auth_proxy_role_binding.yaml
│       ├── auth_proxy_service.yaml
│       ├── kustomization.yaml
│       ├── leader_election_role.yaml
│       ├── leader_election_role_binding.yaml
│       ├── role_binding.yaml
│       └── service_account.yaml
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
└── main.go

6 directories, 24 files
```

</details>

### 2. Create API resource and controller

```
kubebuilder create api --group samplecontroller --version v1alpha1 --kind Foo --controller --resource
```

<details><summary>result</summary>

```bash
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
api/v1alpha1/foo_types.go
controllers/foo_controller.go
Update dependencies:
$ go mod tidy
Running make:
$ make generate
go: creating new go.mod: module tmp
Downloading sigs.k8s.io/controller-tools/cmd/controller-gen@v0.7.0
go get: installing executables with 'go get' in module mode is deprecated.
        To adjust and download dependencies of the current module, use 'go get -d'.
        To install using requirements of the current module, use 'go install'.
        To install ignoring the current module, use 'go install' with a version,
        like 'go install example.com/cmd@latest'.
        For more information, see https://golang.org/doc/go-get-install-deprecation
        or run 'go help get' or 'go help install'.
go get: added github.com/fatih/color v1.12.0
go get: added github.com/go-logr/logr v0.4.0
go get: added github.com/gobuffalo/flect v0.2.3
go get: added github.com/gogo/protobuf v1.3.2
go get: added github.com/google/go-cmp v0.5.6
go get: added github.com/google/gofuzz v1.1.0
go get: added github.com/inconshreveable/mousetrap v1.0.0
go get: added github.com/json-iterator/go v1.1.11
go get: added github.com/mattn/go-colorable v0.1.8
go get: added github.com/mattn/go-isatty v0.0.12
go get: added github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
go get: added github.com/modern-go/reflect2 v1.0.1
go get: added github.com/spf13/cobra v1.2.1
go get: added github.com/spf13/pflag v1.0.5
go get: added golang.org/x/mod v0.4.2
go get: added golang.org/x/net v0.0.0-20210520170846-37e1c6afe023
go get: added golang.org/x/sys v0.0.0-20210616094352-59db8d763f22
go get: added golang.org/x/text v0.3.6
go get: added golang.org/x/tools v0.1.5
go get: added golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
go get: added gopkg.in/inf.v0 v0.9.1
go get: added gopkg.in/yaml.v2 v2.4.0
go get: added gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
go get: added k8s.io/api v0.22.2
go get: added k8s.io/apiextensions-apiserver v0.22.2
go get: added k8s.io/apimachinery v0.22.2
go get: added k8s.io/klog/v2 v2.9.0
go get: added k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
go get: added sigs.k8s.io/controller-tools v0.7.0
go get: added sigs.k8s.io/structured-merge-diff/v4 v4.1.2
go get: added sigs.k8s.io/yaml v1.2.0
/Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

</details>

### 3. Define API resource `Foo` (types.go)

FooSpec:

```go
type FooSpec struct {
    // +kubebuilder:validation:Required
    // +kubebuilder:validation:Format:=string

    DeploymentName string `json:"deploymentName"`

    // +kubebuilder:validation:Required
    // +kubebuilder:validation:Minimum=0

    // the replicas of deployment which is owned by foo
    Replicas *int32 `json:"replicas"`
}
```

FooStatus:

```go
// FooStatus defines the observed state of Foo
type FooStatus struct {
    // this is equal deployment.status.availableReplicas
    // +optional
    AvailableReplicas int32 `json:"availableReplicas"`
}
```

```
make manifests fmt
```

(optional) Run the operator

```
make install run
```

If you see the following logs, it's running successfully.

```
go fmt ./...
go vet ./...
go run ./main.go
2021-12-13T06:32:46.224+0900    INFO    controller-runtime.metrics      metrics server is starting to listen    {"addr": ":8080"}
2021-12-13T06:32:46.225+0900    INFO    setup   starting manager
2021-12-13T06:32:46.226+0900    INFO    starting metrics server {"path": "/metrics"}
2021-12-13T06:32:46.226+0900    INFO    controller.foo  Starting EventSource    {"reconciler group": "samplecontroller.example.com", "reconciler kind": "Foo", "source": "kind source: /, Kind="}
2021-12-13T06:32:46.226+0900    INFO    controller.foo  Starting Controller     {"reconciler group": "samplecontroller.example.com", "reconciler kind": "Foo"}
2021-12-13T06:32:46.327+0900    INFO    controller.foo  Starting workers        {"reconciler group": "samplecontroller.example.com", "reconciler kind": "Foo", "worker count": 1}
```

### 4. Implement reconciliation logic (foo_controller.go and main.go)

```go
import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	samplecontrollerv1alpha1 "github.com/nakamasato/foo-controller-kubebuilder/api/v1alpha1"
)
```

```go
type FooReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}
```

```go
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
```

```go
func (r *FooReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 1. Get Foo object
	var foo samplecontrollerv1alpha1.Foo
	log.Info("fetching Foo Resource")
	if err := r.Get(ctx, req.NamespacedName, &foo); err != nil {
		log.Error(err, "unable to fetch Foo")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// 2. Delete old Deployment if exists
	if err := r.cleanupOwnedResources(ctx, log, &foo); err != nil {
		log.Error(err, "failed to clean up old Deployment resources for this Foo")
		return ctrl.Result{}, err
	}
	// 3. Create or update Deployment that Foo manages
	deploymentName := foo.Spec.DeploymentName
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: req.Namespace,
		},
	}

	if _, err := ctrl.CreateOrUpdate(ctx, r.Client, deploy, func() error {
		replicas := int32(1)
		if foo.Spec.Replicas != nil {
			replicas = *foo.Spec.Replicas
		}
		deploy.Spec.Replicas = &replicas
		labels := map[string]string{
			"app":        "nginx",
			"controller": req.Name,
		}
		if deploy.Spec.Selector == nil {
			deploy.Spec.Selector = &metav1.LabelSelector{MatchLabels: labels}
		}
		if deploy.Spec.Template.ObjectMeta.Labels == nil {
			deploy.Spec.Template.ObjectMeta.Labels = labels
		}
		containers := []corev1.Container{
			{
				Name:  "nginx",
				Image: "nginx:latest",
			},
		}
		if deploy.Spec.Template.Spec.Containers == nil {
			deploy.Spec.Template.Spec.Containers = containers
		}
		if err := ctrl.SetControllerReference(&foo, deploy, r.Scheme); err != nil {
			log.Error(err, "unable to set ownerReference from Foo to Deployment")
			return err
		}
		return nil
	}); err != nil {
		log.Error(err, "unable to ensure deployment is correct")
		return ctrl.Result{}, err
	}
	// 4. Update Foo's status
	var deployment appsv1.Deployment
	var deploymentNamespacedName = client.ObjectKey{Namespace: req.Namespace, Name: foo.Spec.DeploymentName}
	if err := r.Get(ctx, deploymentNamespacedName, &deployment); err != nil {
		log.Error(err, "unable to fetch Deployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	availableReplicas := deployment.Status.AvailableReplicas
	if availableReplicas == foo.Status.AvailableReplicas {
		return ctrl.Result{}, nil
	}
	foo.Status.AvailableReplicas = availableReplicas
	if err := r.Status().Update(ctx, &foo); err != nil {
		log.Error(err, "unable to update Foo status")
		return ctrl.Result{}, err
	}
	r.Recorder.Eventf(&foo, corev1.EventTypeNormal, "Updated", "Update foo.status.AvailableReplicas: %d", foo.Status.AvailableReplicas)
	return ctrl.Result{}, nil
}

func (r *FooReconciler) cleanupOwnedResources(ctx context.Context, log logr.Logger, foo *samplecontrollerv1alpha1.Foo) error {
	log.Info("finding existing Deployments for Foo resource")

	var deployments appsv1.DeploymentList
	if err := r.List(ctx, &deployments, client.InNamespace(foo.Namespace), client.MatchingFields(
		map[string]string{
			deploymentOwnerKey: foo.Name,
		},
	)); err != nil {
		return err
	}

	for _, deployment := range deployments.Items {
		if deployment.Name == foo.Spec.DeploymentName {
			continue
		}

		if err := r.Delete(ctx, &deployment); err != nil {
			log.Error(err, "failed to delete Deployment resource")
			return err
		}

		log.Info("delete deployment resource: " + deployment.Name)
		r.Recorder.Eventf(foo, corev1.EventTypeNormal, "Deleted", "Deleted deployment %q", deployment.Name)
	}
	return nil
}
```

```go
var (
	deploymentOwnerKey = ".metadata.controller"
	apiGVStr           = samplecontrollerv1alpha1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *FooReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{}, deploymentOwnerKey, func(o client.Object) []string {
		deployment := o.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(deployment)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != "Foo" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplecontrollerv1alpha1.Foo{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
```

Update `main.go`

```go
if err = (&controllers.FooReconciler{
    Client:   mgr.GetClient(),
    Scheme:   mgr.GetScheme(),
    Recorder: mgr.GetEventRecorderFor("foo-controller"),
}).SetupWithManager(mgr); err != nil {
    setupLog.Error(err, "unable to create controller", "controller", "Foo")
    os.Exit(1)
}
```

### 5. Run the operator

```
make install
```

```
make install
/Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/foos.samplecontroller.example.com configured
```

Run with go

```
make run
```

Create sample `Foo`

```
kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
```

Result:

```
kubectl get deploy
NAME   READY   UP-TO-DATE   AVAILABLE   AGE
foo    2/2     2            2           9s
```

```
kubectl get event | grep foo/
34s         Normal   Updated             foo/foo-sample              Update foo.status.AvailableReplicas: 1
32s         Normal   Updated             foo/foo-sample              Update foo.status.AvailableReplicas: 2
```

Clean up

```
kubectl delete -f config/samples/samplecontroller_v1alpha1_foo.yaml
```

Run with container

```
export IMG=nakamasato/foo-controller:kubebuilder
```

Build docker image

```
make docker-build
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

<details>

```bash
/Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && /Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/kustomize edit set image controller=nakamasato/foo-controller:kubebuilder
/Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/kustomize build config/default | kubectl apply -f -
namespace/foo-controller-kubebuilder-system created
customresourcedefinition.apiextensions.k8s.io/foos.samplecontroller.example.com configured
serviceaccount/foo-controller-kubebuilder-controller-manager created
role.rbac.authorization.k8s.io/foo-controller-kubebuilder-leader-election-role created
clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-manager-role created
clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-metrics-reader created
clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-proxy-role created
rolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-proxy-rolebinding created
configmap/foo-controller-kubebuilder-manager-config created
service/foo-controller-kubebuilder-controller-manager-metrics-service created
deployment.apps/foo-controller-kubebuilder-controller-manager created
```

</details>

```
kubectl get deploy -n foo-controller-kubebuilder-system
NAME                                            READY   UP-TO-DATE   AVAILABLE   AGE
foo-controller-kubebuilder-controller-manager   1/1     1            1           60s
```

```
kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
foo.samplecontroller.example.com/foo-sample created
```

check deployment

```
kubectl get deploy
NAME   READY   UP-TO-DATE   AVAILABLE   AGE
foo    2/2     2            2           13s
```

clean up

```
kubectl delete -f config/samples/samplecontroller_v1alpha1_foo.yaml
foo.samplecontroller.example.com "foo-sample" deleted
```

```
make undeploy
```

## 5. Create Admission webhook

```
kubebuilder create webhook --group samplecontroller --version v1alpha1 --kind Foo --defaulting --programmatic-validation
```

<details><summary>result</summary>

```
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
api/v1alpha1/foo_webhook.go
Update dependencies:
$ go mod tidy
Running make:
$ make generate
/Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new Webhook and generate the manifests with:
$ make manifests
```

</details>

## 6. Implement admission webhook
### 6.1. default

Set replicas to optional in `api/v1alpha1/foo_types.go`

```go
    // +kubebuilder:validation:Optional

	// the replicas of deployment which is owned by foo
	Replicas *int32 `json:"replicas"`
```

Default: `api/v1alpha1/foo_webhook.go`

```go
func (r *Foo) Default() {
	foolog.Info("default", "name", r.Name)

	if r.Spec.Replicas == nil {
		r.Spec.Replicas = new(int32)
		*r.Spec.Replicas = 1
	}
}
```

### 6.2. validataion

Validation: deploymentName must not be no more than 253 characters

You can use `// +kubebuilder:validation:MaxLength` with a maker.

```go
import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
    ...
)
...
// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateCreate() error {
	foolog.Info("validate create", "name", r.Name)

	return r.validateFoo()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateUpdate(old runtime.Object) error {
	foolog.Info("validate update", "name", r.Name)

	return r.validateFoo()
}
...
func (r *Foo) validateDeploymentName() *field.Error {
	// depoymentName must be no more than 253 characters.
	if (len(r.Spec.DeploymentName) > 253) {
		return field.Invalid(field.NewPath("spec").Child("deploymentName"), r.Spec.DeploymentName, "must be no more than 253 characters")
	}
	return nil
}

func (r *Foo) validateFoo() error {
	var allErrs field.ErrorList
	if err := r.validateDeploymentName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "samplecontroller.example.com", Kind: "Foo"}, r.Name, allErrs)
}
```

### 6.3. Execute

1. prepare certificate

```
kubectl create ns cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml
```

```
kubectl get pod -n cert-manager
NAME                                      READY   STATUS    RESTARTS   AGE
cert-manager-55658cdf68-5w58w             1/1     Running   0          39s
cert-manager-cainjector-967788869-p4psd   1/1     Running   0          39s
cert-manager-webhook-6668fbb57d-44sjk     1/1     Running   0          39s
```

1. Build Docker image

    ```
    export IMG=nakamasato/foo-controller:kubebuilder-webhook
    make docker-build docker-push
    ```

1. Update `config/default/kustomization.yaml`

    Uncomment the followings:

    ```yaml
    bases:
    - ../webhook
    - ../certmanager
    ...
    patchesStrategicMerge:
    - manager_webhook_patch.yaml
    - webhookcainjection_patch.yaml
    ...
    vars:
    # [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'     prefix.
    - name: CERTIFICATE_NAMESPACE # namespace of the certificate CR
      objref:
        kind: Certificate
        group: cert-manager.io
        version: v1
        name: serving-cert # this name should match the one in certificate.yaml
      fieldref:
        fieldpath: metadata.namespace
    - name: CERTIFICATE_NAME
      objref:
        kind: Certificate
        group: cert-manager.io
        version: v1
        name: serving-cert # this name should match the one in certificate.yaml
    - name: SERVICE_NAMESPACE # namespace of the service
      objref:
        kind: Service
        version: v1
        name: webhook-service
      fieldref:
        fieldpath: metadata.namespace
    - name: SERVICE_NAME
      objref:
        kind: Service
        version: v1
        name: webhook-service
    ```

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
1. Check foo controller
    ```bash
    kubectl get deploy -n foo-controller-kubebuilder-system
    NAME                                            READY   UP-TO-DATE   AVAILABLE   AGE
    foo-controller-kubebuilder-controller-manager   1/1     1            1           29s
    ```
    ```bash
    kubectl get pod -n foo-controller-kubebuilder-system
    NAME                                                             READY   STATUS    RESTARTS   AGE
    foo-controller-kubebuilder-controller-manager-65dc97ddfd-lqdd2   2/2     Running   1          2m23s
    ```
    ```bash
    kubectl logs -n foo-controller-kubebuilder-system -l control-plane=controller-manager
    2021-12-13T00:17:03.690Z        INFO    controller-runtime.webhook      serving webhook server  {"host": "", "port": 9443}
    2021-12-13T00:17:03.689Z        INFO    starting metrics server {"path": "/metrics"}
    2021-12-13T00:17:03.690Z        INFO    controller-runtime.certwatcher  Starting certificate watcher
    I1213 00:17:03.706798       1 leaderelection.go:258] successfully acquired lease foo-controller-kubebuilder-system/78433493.example.com
    2021-12-13T00:17:03.706Z        DEBUG   events  Normal  {"object": {"kind":"ConfigMap","namespace":"foo-controller-kubebuilder-system","name":"78433493.example.com","uid":"c91a0b2f-068b-4865-8406-63a6068aa168","apiVersion":"v1","resourceVersion":"262417"}, "reason": "LeaderElection", "message": "foo-controller-kubebuilder-controller-manager-6cc4ff948-pbgwq_809926f5-3199-4f77-be93-af9a7758bf04 became leader"}
    2021-12-13T00:17:03.707Z        DEBUG   events  Normal  {"object": {"kind":"Lease","namespace":"foo-controller-kubebuilder-system","name":"78433493.example.com","uid":"7f92a5a8-86d0-40f5-8bc7-3f86ab1ebe13","apiVersion":"coordination.k8s.io/v1","resourceVersion":"262418"}, "reason": "LeaderElection", "message": "foo-controller-kubebuilder-controller-manager-6cc4ff948-pbgwq_809926f5-3199-4f77-be93-af9a7758bf04 became leader"}
    2021-12-13T00:17:03.790Z        INFO    controller.foo  Starting EventSource    {"reconciler group": "samplecontroller.example.com", "reconciler kind": "Foo", "source": "kind source: /, Kind="}
    2021-12-13T00:17:03.790Z        INFO    controller.foo  Starting EventSource    {"reconciler group": "samplecontroller.example.com", "reconciler kind": "Foo", "source": "kind source: /, Kind="}
    2021-12-13T00:17:03.790Z        INFO    controller.foo  Starting Controller     {"reconciler group": "samplecontroller.example.com", "reconciler kind": "Foo"}
    2021-12-13T00:17:03.891Z        INFO    controller.foo  Starting workers        {"reconciler group": "samplecontroller.example.com", "reconciler kind": "Foo", "worker count": 1}
    ```

1. Apply `config/samples/samplecontroller_v1alpha1_foo.yaml` without `replicas: 1`

    ```
    kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
    foo.samplecontroller.example.com/foo-sample created
    ```

1. Check replica is 1

    ```yaml
    spec:
      # replicas: 2 # to check default admission webhook
    ```

    ```
    kubectl get foo foo-sample -o jsonpath='{.spec.replicas}'
    1
    ```

1. Apply `config/samples/samplecontroller_v1alpha1_foo.yaml` with too long name

    ```yaml
    spec:
      deploymentName: aaaaaaaaaaaa...aaaaaaaaaaaaa # more than 253 letters
    ```

    ```
    kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
    The Foo "foo-sample" is invalid: spec.deploymentName: Invalid value: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa": must be no more than 253 characters
    ```

### 6.4. clean up

```
kubectl delete -f config/samples/samplecontroller_v1alpha1_foo.yaml
```

```
make undeploy
```

```
```
## 7. Add new API version

```
kubebuilder create api --group samplecontroller --version v1beta1 --kind Foo --resource --controller=False
```

<details><summary>result</summary>

```
Writing kustomize manifests for you to edit...
Writing scaffold for you to edit...
api/v1beta1/foo_types.go
Update dependencies:
$ go mod tidy
Running make:
$ make generate
/Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

</details>

## 8. Create Conversion Webhook

Conversion: API version compatibility

`apps/v1alpha`のResourceがApplyされても、`apps/v1`のResourceがApplyされても、KuberntesのControllerがそれをReconcileし、同じ機能・効果を提供するObjectを保証してくれます。これは、Kubernetesが裏側でConversionを実行することで、Multi Version間の差分を吸収してくれているからです。

1. Copy `v1alpha1/foo_types.go` to `v1beta1/foo_types.go`
    ```diff
    diff api/v1{alpha1,beta1}/foo_types.go
    17c17
    < package v1alpha1
    ---
    > package v1beta1
    ```
1. Add new field for conversion to `FooSpec`

    ```go
    type FooSpec struct {
        ...
	    // +kubebuilder:validation:Optional

	    // the new field for conversion
	    Foo string `json:"foo"`
    }
    ```

1. Add `//+kubebuilder:storageversion` to `alphav1/foo_types.go`
    > Since we’ll have more than one version, we’ll need to mark a storage version. This is the version that the Kubernetes API server uses to store our data.
1. Implement conversion.
    1. `api/v1alpha1/foo_conversion.go`:
        ```go
        package v1alpha1

        func (*Foo) Hub() {}
        ```
    1. `api/v1beta1/foo_conversion.go`:
        ```go
        package v1beta1

        import (
	        "sigs.k8s.io/controller-runtime/pkg/conversion"

	        samplecontrollerv1alpha1 "github.com/nakamasato/foo-controller-kubebuilder/api/v1alpha1"
        )

        func (src *Foo) ConvertTo(dstRaw conversion.Hub) error { // v1beta1 -> v1alpha1
            ...
        }

        func (dst *Foo) ConvertFrom(srcRaw conversion.Hub) error { // v1alpha1 -> v1beta1
            ...
            // just copy DeploymentName value to Foo
	        dst.Spec.Foo = src.Spec.DeploymentName
            ...
        }
        ```
1. Update `config/crd/kustomization.yaml`
    ```yaml
    patchesStrategicMerge:
    - patches/webhook_in_foos.yaml
    - patches/cainjection_in_foos.yaml
    ```

1. Build & Push

    ```
    export IMG=nakamasato/foo-controller:kubebuilder-conversion
    make docker-build docker-push
    ```
1. (install cert-manager if uninstalled)
1. Deploy

    ```
    make deploy
    ```

    <details><summary>result</summary>

    ```
    /Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
    cd config/manager && /Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/kustomize edit set image controller=nakamasato/foo-controller:kubebuilder-conversion
    /Users/masato-naka/repos/nakamasato/foo-controller-kubebuilder/bin/kustomize build config/default | kubectl apply -f -
    namespace/foo-controller-kubebuilder-system created
    customresourcedefinition.apiextensions.k8s.io/foos.samplecontroller.example.com created
    serviceaccount/foo-controller-kubebuilder-controller-manager created
    role.rbac.authorization.k8s.io/foo-controller-kubebuilder-leader-election-role created
    clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-manager-role created
    clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-metrics-reader created
    clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-proxy-role created
    rolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-leader-election-rolebinding created
    clusterrolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-manager-rolebinding created
    clusterrolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-proxy-rolebinding created
    configmap/foo-controller-kubebuilder-manager-config created
    service/foo-controller-kubebuilder-controller-manager-metrics-service created
    service/foo-controller-kubebuilder-webhook-service created
    deployment.apps/foo-controller-kubebuilder-controller-manager created
    certificate.cert-manager.io/foo-controller-kubebuilder-serving-cert created
    issuer.cert-manager.io/foo-controller-kubebuilder-selfsigned-issuer created
    mutatingwebhookconfiguration.admissionregistration.k8s.io/foo-controller-kubebuilder-mutating-webhook-configuration created
    validatingwebhookconfiguration.admissionregistration.k8s.io/foo-controller-kubebuilder-validating-webhook-configuration created
    ```

    </details>

1. Check

    ```
    kubectl get deployment -n foo-controller-kubebuilder-system
    ```

    ```
    kubectl get pods -n foo-controller-kubebuilder-system
    ```

    ```
    kubectl logs -n foo-controller-kubebuilder-system -l control-plane=controller-manager -f
    ```

    `conversion webhook enabled` <- cannot see

1. Check: apply `betav1` -> get `betav1` and `alphav1`

    Apply `betav1`

    ```
    kubectl apply -f config/samples/samplecontroller_v1beta1_foo.yaml
    foo.samplecontroller.k8s.io/foo-sample created
    ```

    confirm `beta1`

    ```
    kubectl get foo.v1beta1.samplecontroller.example.com foo-sample -o jsonpath='{.spec}'
    {"deploymentName":"example-foo","foo":"example-foo","replicas":1}%
    ```

    Check `alpha1` -> no foo

    ```
    kubectl get foo.v1alpha1.samplecontroller.example.com foo-sample -o jsonpath='{.spec}'
    {"deploymentName":"example-foo","replicas":1}
    ```

1. Check: apply `alphav1` -> get `betav1` and `alphav1`

    ```
    kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
    foo.samplecontroller.example.com/foo-sample created

    kubectl get foo.v1alpha1.samplecontroller.example.com foo-sample -o jsonpath='{.spec}'
    {"deploymentName":"sample-foo","replicas":1}%

    kubectl get foo.v1beta1.samplecontroller.example.com foo-sample -o jsonpath='{.spec}'
    {"deploymentName":"sample-foo","foo":"sample-foo","replicas":1}%
    ```

1. Clean up

    ```
    kubectl delete -f config/samples/samplecontroller_v1beta1_foo.yaml
    make undeploy
    kubectl delete -f https://github.com/jetstack/cert-manager/releases/download/v1.6.1/cert-manager.yaml
    ```
