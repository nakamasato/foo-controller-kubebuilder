# foo-controller

- [実践入門 Kubernetesカスタムコントローラへの道](https://www.amazon.co.jp/-/en/gp/product/B0851QCR81/ref=ppx_yo_dt_b_d_asin_title_o00?ie=UTF8&psc=1)
- https://github.com/govargo/foo-controller-kubebuilder

## Version

- go: go1.20.6
- kubebuilder: 3.12.0
- kustomize: v5.1.1
- cert-manager: v1.13.0

## Prerequisite

```
go version
go version go1.20.6 darwin/amd64
```

install kubebuilder

```bash
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
```

```
kubebuilder version
Version: main.version{KubeBuilderVersion:"3.12.0", KubernetesVendor:"1.27.1", GitCommit:"b48f95cd5384eadcdfd02a47a02910f72ddc7ea8", BuildDate:"2023-09-06T06:04:11Z", GoOs:"darwin", GoArch:"arm64"}
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

go 1.20
```

Initialize a kubebuilder project with plugin `go/v4`

```
kubebuilder init --domain example.com --owner nakamasato
```

<details><summary>tree</summary>

```bash
tree
.
├── Dockerfile
├── Makefile
├── PROJECT
├── README.md
├── cmd
│   └── main.go
├── config
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   └── manager_config_patch.yaml
│   ├── manager
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
└── hack
    └── boilerplate.go.txt

7 directories, 24 files
```

</details>

### 2. Create API resource and controller

```
kubebuilder create api --group samplecontroller --version v1alpha1 --kind Foo --controller --resource
```

<details><summary>result</summary>

```bash
INFO[0000] Writing kustomize manifests for you to edit...
INFO[0000] Writing scaffold for you to edit...
INFO[0000] api/v1alpha1/foo_types.go
INFO[0000] api/v1alpha1/groupversion_info.go
INFO[0000] internal/controller/suite_test.go
INFO[0000] internal/controller/foo_controller.go
INFO[0000] Update dependencies:
$ go mod tidy
INFO[0000] Running make:
$ make generate
mkdir -p /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin
test -s /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen && /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen --version | grep -q v0.13.0 || \
        GOBIN=/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
go: downloading sigs.k8s.io/controller-tools v0.13.0
/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

</details>

```
make manifests
```

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
go run ./cmd/main.go
2023-09-27T07:42:37+09:00       INFO    setup   starting manager
2023-09-27T07:42:37+09:00       INFO    starting server {"kind": "health probe", "addr": "[::]:8081"}
2023-09-27T07:42:37+09:00       INFO    controller-runtime.metrics      Starting metrics server
2023-09-27T07:42:37+09:00       INFO    controller-runtime.metrics      Serving metrics server  {"bindAddress": ":8080", "secure": false}
2023-09-27T07:42:37+09:00       INFO    Starting EventSource    {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo", "source": "kind source: *v1alpha1.Foo"}
2023-09-27T07:42:37+09:00       INFO    Starting Controller     {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo"}
2023-09-27T07:42:37+09:00       INFO    Starting workers        {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo", "worker count": 1}
```

### 4. Implement reconciliation logic (foo_controller.go and main.go)

`internal/controller/foo_controller.go`:

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

```diff
  type FooReconciler struct {
  	client.Client
  	Scheme   *runtime.Scheme
+  	Recorder record.EventRecorder
  }
```

```diff
  //+kubebuilder:rbac:groups=samplecontroller.example.com,resources=foos,verbs=get;list;watch;create;update;patch;delete
  //+kubebuilder:rbac:groups=samplecontroller.example.com,resources=foos/status,verbs=get;update;patch
  //+kubebuilder:rbac:groups=samplecontroller.example.com,resources=foos/finalizers,verbs=update
+ //+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete
+ //+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
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

Update `cmd/main.go`

```diff
  if err = (&controllers.FooReconciler{
      Client:   mgr.GetClient(),
      Scheme:   mgr.GetScheme(),
+     Recorder: mgr.GetEventRecorderFor("foo-controller"),
  }).SetupWithManager(mgr); err != nil {
      setupLog.Error(err, "unable to create controller", "controller", "Foo")
      os.Exit(1)
  }
```

```
make fmt
```

Run the operator

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

```yaml
apiVersion: samplecontroller.example.com/v1alpha1
kind: Foo
metadata:
  labels:
    app.kubernetes.io/name: foo
    app.kubernetes.io/instance: foo-sample
    app.kubernetes.io/part-of: foo-controller-kubebuilder
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: foo-controller-kubebuilder
  name: foo-sample
spec:
  deploymentName: "foo-sample"
  replicas: 2
```

```
kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
```

Result:

```
kubectl get deploy
NAME         READY   UP-TO-DATE   AVAILABLE   AGE
foo-sample   2/2     2            2           21s
```

```
kubectl get event | grep foo/
30s         Normal   Updated                   foo/foo-sample                     Update foo.status.AvailableReplicas: 1
28s         Normal   Updated                   foo/foo-sample                     Update foo.status.AvailableReplicas: 2
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
test -s /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen && /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen --version | grep -q v0.13.0 || \
        GOBIN=/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/kustomize edit set image controller=nakamasato/foo-controller:kubebuilder
/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/kustomize build config/default | kubectl apply -f -
# Warning: 'patchesStrategicMerge' is deprecated. Please use 'patches' instead. Run 'kustomize edit fix' to update your Kustomization automatically.
namespace/foo-controller-kubebuilder-system created
customresourcedefinition.apiextensions.k8s.io/foos.samplecontroller.example.com unchanged
serviceaccount/foo-controller-kubebuilder-controller-manager created
role.rbac.authorization.k8s.io/foo-controller-kubebuilder-leader-election-role created
clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-manager-role created
clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-metrics-reader created
clusterrole.rbac.authorization.k8s.io/foo-controller-kubebuilder-proxy-role created
rolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/foo-controller-kubebuilder-proxy-rolebinding created
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
NAME         READY   UP-TO-DATE   AVAILABLE   AGE
foo-sample   2/2     2            2           7s
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
INFO[0000] Writing kustomize manifests for you to edit...
INFO[0000] Writing scaffold for you to edit...
INFO[0000] api/v1alpha1/foo_webhook.go
INFO[0000] api/v1alpha1/webhook_suite_test.go
INFO[0007] Update dependencies:
$ go mod tidy
INFO[0008] Running make:
$ make generate
test -s /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen && /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen --version | grep -q v0.13.0 || \
        GOBIN=/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new Webhook and generate the manifests with:
$ make manifests
```

</details>

```
make manifests
```

## 6. Implement admission webhook
### 6.1. default

Set replicas to optional in `api/v1alpha1/foo_types.go`

```diff
-  	// +kubebuilder:validation:Required
+  	// +kubebuilder:validation:Optional

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


We'll implement the [Validator](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.2/pkg/webhook/admission#Validator) interface:

```go
type Validator interface {
	runtime.Object

	// ValidateCreate validates the object on creation.
	// The optional warnings will be added to the response as warning messages.
	// Return an error if the object is invalid.
	ValidateCreate() (warnings Warnings, err error)

	// ValidateUpdate validates the object on update. The oldObj is the object before the update.
	// The optional warnings will be added to the response as warning messages.
	// Return an error if the object is invalid.
	ValidateUpdate(old runtime.Object) (warnings Warnings, err error)

	// ValidateDelete validates the object on deletion.
	// The optional warnings will be added to the response as warning messages.
	// Return an error if the object is invalid.
	ValidateDelete() (warnings Warnings, err error)
}
```

The return value was updated in [controller-runtime@v0.15.0](https://github.com/kubernetes-sigs/controller-runtime/releases/tag/v0.15.0) ([⚠️ feat: new features about support warning with webhook #2014](https://github.com/kubernetes-sigs/controller-runtime/pull/2014)) from [[Feature Request]: Support "Warning" for Validation Webhook #1896](https://github.com/kubernetes-sigs/controller-runtime/issues/1896)

This is because Kubernets supports `warning` message in response for Admission webhook since [1.19](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.19.md#deprecation-warnings) [ref](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#response):

> Admission webhooks can optionally return warning messages that are returned to the requesting client in HTTP Warning headers with a warning code of 299. Warnings can be sent with allowed or rejected admission responses.


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
func (r *Foo) ValidateCreate() (admission.Warnings, error) {
	foolog.Info("validate create", "name", r.Name)

	return r.validateFoo()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
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

func (r *Foo) validateFoo() (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := r.validateDeploymentName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: "samplecontroller.example.com", Kind: "Foo"}, r.Name, allErrs)
}
```

### 6.3. Execute

1. prepare certificate

    ```
    kubectl create ns cert-manager
    kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.13.0/cert-manager.yaml
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

    Uncomment the codes for `[WEBHOOK]` and `[CERTMANAGER]`:

    ```yaml
    resources:
    - ../crd
    - ../rbac
    - ../manager
    # [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
    # crd/kustomization.yaml
    - ../webhook
    # [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'. 'WEBHOOK' components are required.
    - ../certmanager
    # [PROMETHEUS] To enable prometheus monitor, uncomment all sections with 'PROMETHEUS'.
    #- ../prometheus

    patchesStrategicMerge:
    # Protect the /metrics endpoint by putting it behind auth.
    # If you want your controller-manager to expose the /metrics
    # endpoint w/o any authn/z, please comment the following line.
    - manager_auth_proxy_patch.yaml



    # [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
    # crd/kustomization.yaml
    - manager_webhook_patch.yaml

    # [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'.
    # Uncomment 'CERTMANAGER' sections in crd/kustomization.yaml to enable the CA injection in the admission webhooks.
    # 'CERTMANAGER' needs to be enabled to use ca injection
    - webhookcainjection_patch.yaml

    # [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER' prefix.
    # Uncomment the following replacements to add the cert-manager CA injection annotations
    replacements:
     - source: # Add cert-manager annotation to ValidatingWebhookConfiguration, MutatingWebhookConfiguration and CRDs
         kind: Certificate
         group: cert-manager.io
         version: v1
         name: serving-cert # this name should match the one in certificate.yaml
         fieldPath: .metadata.namespace # namespace of the certificate CR
       targets:
         - select:
             kind: ValidatingWebhookConfiguration
           fieldPaths:
             - .metadata.annotations.[cert-manager.io/inject-ca-from]
           options:
             delimiter: '/'
             index: 0
             create: true
         - select:
             kind: MutatingWebhookConfiguration
           fieldPaths:
             - .metadata.annotations.[cert-manager.io/inject-ca-from]
           options:
             delimiter: '/'
             index: 0
             create: true
         - select:
             kind: CustomResourceDefinition
           fieldPaths:
             - .metadata.annotations.[cert-manager.io/inject-ca-from]
           options:
             delimiter: '/'
             index: 0
             create: true
     - source:
         kind: Certificate
         group: cert-manager.io
         version: v1
         name: serving-cert # this name should match the one in certificate.yaml
         fieldPath: .metadata.name
       targets:
         - select:
             kind: ValidatingWebhookConfiguration
           fieldPaths:
             - .metadata.annotations.[cert-manager.io/inject-ca-from]
           options:
             delimiter: '/'
             index: 1
             create: true
         - select:
             kind: MutatingWebhookConfiguration
           fieldPaths:
             - .metadata.annotations.[cert-manager.io/inject-ca-from]
           options:
             delimiter: '/'
             index: 1
             create: true
         - select:
             kind: CustomResourceDefinition
           fieldPaths:
             - .metadata.annotations.[cert-manager.io/inject-ca-from]
           options:
             delimiter: '/'
             index: 1
             create: true
     - source: # Add cert-manager annotation to the webhook Service
         kind: Service
         version: v1
         name: webhook-service
         fieldPath: .metadata.name # namespace of the service
       targets:
         - select:
             kind: Certificate
             group: cert-manager.io
             version: v1
           fieldPaths:
             - .spec.dnsNames.0
             - .spec.dnsNames.1
           options:
             delimiter: '.'
             index: 0
             create: true
     - source:
         kind: Service
         version: v1
         name: webhook-service
         fieldPath: .metadata.namespace # namespace of the service
       targets:
         - select:
             kind: Certificate
             group: cert-manager.io
             version: v1
           fieldPaths:
             - .spec.dnsNames.0
             - .spec.dnsNames.1
           options:
             delimiter: '.'
             index: 1
             create: true
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
    2023-09-26T23:24:37Z    INFO    controller-runtime.certwatcher  Updated current TLS certificate
    2023-09-26T23:24:37Z    INFO    controller-runtime.webhook      Serving webhook server  {"host": "", "port": 9443}
    2023-09-26T23:24:37Z    INFO    controller-runtime.certwatcher  Starting certificate watcher
    I0926 23:24:37.521249       1 leaderelection.go:250] attempting to acquire leader lease foo-controller-kubebuilder-system/78433493.example.com...
    I0926 23:24:37.543740       1 leaderelection.go:260] successfully acquired lease foo-controller-kubebuilder-system/78433493.example.com
    2023-09-26T23:24:37Z    DEBUG   events  foo-controller-kubebuilder-controller-manager-679b458967-kpxts_afc3db92-c403-445c-ae60-02d0288fad33 became leader     {"type": "Normal", "object": {"kind":"Lease","namespace":"foo-controller-kubebuilder-system","name":"78433493.example.com","uid":"bc33e17a-5070-41d3-891e-34e93df99437","apiVersion":"coordination.k8s.io/v1","resourceVersion":"157918"}, "reason": "LeaderElection"}
    2023-09-26T23:24:37Z    INFO    Starting EventSource    {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo", "source": "kind source: *v1alpha1.Foo"}
    2023-09-26T23:24:37Z    INFO    Starting EventSource    {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo", "source": "kind source: *v1.Deployment"}
    2023-09-26T23:24:37Z    INFO    Starting Controller     {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo"}
    2023-09-26T23:24:37Z    INFO    Starting workers        {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo", "worker count": 1}
    ```

1. Apply `config/samples/samplecontroller_v1alpha1_foo.yaml` without `replicas`

    ```yaml
    spec:
      # replicas: 2 # to check default admission webhook
    ```

    ```
    kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
    foo.samplecontroller.example.com/foo-sample created
    ```

1. Check replica is 1

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
kubectl delete -f https://github.com/jetstack/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```
## 7. Add new API version

```
kubebuilder create api --group samplecontroller --version v1beta1 --kind Foo --resource --controller=False
```

<details><summary>result</summary>

```
INFO[0000] Writing kustomize manifests for you to edit...
INFO[0000] Writing scaffold for you to edit...
INFO[0000] api/v1beta1/foo_types.go
INFO[0000] api/v1beta1/groupversion_info.go
INFO[0000] Update dependencies:
$ go mod tidy
INFO[0000] Running make:
$ make generate
test -s /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen && /Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen --version | grep -q v0.13.0 || \
        GOBIN=/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.13.0
/Users/m.naka/repos/nakamasato/foo-controller-kubebuilder/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
Next: implement your new API and generate the manifests (e.g. CRDs,CRs) with:
$ make manifests
```

</details>


※ At this point, if you run `make manifests`, you'll see `github.com/nakamasato/foo-controller-kubebuilder/api/v1alpha1:-: CRD for Foo.samplecontroller.example.com has no storage version` error as it's not implemented.

## 8. Create Conversion Webhook

Conversion: API version compatibility

Whether applying `apps/v1alpha` or `apps/v1`, Kubernetes controllers ensure to provide objects with the same functionality through reconciliation. This is achieved by Kubernetes performing conversions among multiple versions in the background.

In controller-runtime, we can utilize [conversion](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/conversion) package for implementation.

[Convertible](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/conversion#Convertible) interface:

```go
type Convertible interface {
	runtime.Object
	ConvertTo(dst Hub) error
	ConvertFrom(src Hub) error
}
```

[Hub](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/conversion#Hub) interface:

```go
type Hub interface {
	runtime.Object
	Hub()
}
```

> Hub marks that a given type is the hub type for conversion. This means that all conversions will first convert to the hub type, then convert from the hub type to the destination type. All types besides the hub type should implement Convertible.


Implementation steps:

1. Copy `api/v1alpha1/foo_types.go` to `api/v1beta1/foo_types.go`
    ```diff
    diff api/v1{alpha1,beta1}/foo_types.go
    17c17
    < package v1alpha1
    ---
    > package v1beta1
    ```
1. Add new field for conversion to `FooSpec`

    ```diff
      type FooSpec struct {
         ...
    +    // +kubebuilder:validation:Optional
    +    // the new field for conversion
    +    Foo string `json:"foo"`
      }
    ```

1. Add `//+kubebuilder:storageversion` to `alphav1/foo_types.go`
    > Since we’ll have more than one version, we’ll need to mark a storage version. This is the version that the Kubernetes API server uses to store our data.

    ```diff
      //+kubebuilder:object:root=true
      //+kubebuilder:subresource:status
    + //+kubebuilder:storageversion

      // Foo is the Schema for the foos API
      type Foo struct {
      	metav1.TypeMeta   `json:",inline"`
      	metav1.ObjectMeta `json:"metadata,omitempty"`

      	Spec   FooSpec   `json:"spec,omitempty"`
      	Status FooStatus `json:"status,omitempty"`
      }
    ```

1. Implement conversion.
    1. `api/v1alpha1/foo_conversion.go`:
        ```go
        package v1alpha1

        func (*Foo) Hub() {}
        ```
    1. `api/v1beta1/foo_conversion.go`: `dst.Spec.Foo = src.Spec.DeploymentName` this is the main point of the conversion.
        ```go
        package v1beta1

        import (
        	"sigs.k8s.io/controller-runtime/pkg/conversion"

        	samplecontrollerv1alpha1 "github.com/nakamasato/foo-controller-kubebuilder/api/v1alpha1"
        )

        func (src *Foo) ConvertTo(dstRaw conversion.Hub) error { // v1beta1 -> v1alpha1
        	dst := dstRaw.(*samplecontrollerv1alpha1.Foo)

        	// ObjectMeta
        	dst.ObjectMeta = src.ObjectMeta

        	// Spec
        	dst.Spec.DeploymentName = src.Spec.DeploymentName
        	dst.Spec.Replicas = src.Spec.Replicas

        	// Status
        	dst.Status.AvailableReplicas = src.Status.AvailableReplicas

        	return nil
        }

        func (dst *Foo) ConvertFrom (srcRaw conversion.Hub) error { // v1alpha1 -> v1beta1
        	src := srcRaw.(*samplecontrollerv1alpha1.Foo)

        	// Set the new field in v1beta1 from existing field in v1alphav1
        	dst.Spec.Foo = src.Spec.DeploymentName

        	// ObjectMeta
        	dst.ObjectMeta = src.ObjectMeta

        	// Spec
        	dst.Spec.DeploymentName = src.Spec.DeploymentName
        	dst.Spec.Replicas = src.Spec.Replicas

        	// Status
        	dst.Status.AvailableReplicas = src.Status.AvailableReplicas

        	return nil
        }
        ```
1. Update `config/crd/kustomization.yaml`
    ```yaml
    patches:
    - path: patches/webhook_in_foos.yaml
    - path: patches/cainjection_in_foos.yaml
    ```

1. Build & Push

    ```
    export IMG=nakamasato/foo-controller:kubebuilder-conversion
    make docker-build docker-push
    ```
1. Install cert-manager if uninstalled
    ```
    kubectl create ns cert-manager
    kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.13.0/cert-manager.yaml
    ```
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

    Sometimes the following error occurs:

    ```
    Error from server (InternalError): error when creating "STDIN": Internal error occurred: failed calling webhook "webhook.cert-manager.io": failed to call webhook: Post "https://cert-manager-webhook.cert-manager.svc:443/mutate?timeout=10s": tls: failed to verify certificate: x509: certificate signed by unknown authority
    Error from server (InternalError): error when creating "STDIN": Internal error occurred: failed calling webhook "webhook.cert-manager.io": failed to call webhook: Post "https://cert-manager-webhook.cert-manager.svc:443/mutate?timeout=10s": tls: failed to verify certificate: x509: certificate signed by unknown authority
    make: *** [deploy] Error 1
    ```

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

    ```
    2023-09-28T00:11:32Z    INFO    controller-runtime.certwatcher  Updated current TLS certificate
    2023-09-28T00:11:32Z    INFO    controller-runtime.webhook      Serving webhook server  {"host": "", "port": 9443}
    2023-09-28T00:11:32Z    INFO    controller-runtime.certwatcher  Starting certificate watcher
    I0928 00:11:32.273264       1 leaderelection.go:250] attempting to acquire leader lease foo-controller-kubebuilder-system/78433493.example.com...
    I0928 00:11:32.277137       1 leaderelection.go:260] successfully acquired lease foo-controller-kubebuilder-system/78433493.example.com
    2023-09-28T00:11:32Z    DEBUG   events  foo-controller-kubebuilder-controller-manager-8c9494f56-jt25m_fcc12c04-d25a-44a7-84ec-6bbb69419093 became leader      {"type": "Normal", "object": {"kind":"Lease","namespace":"foo-controller-kubebuilder-system","name":"78433493.example.com","uid":"382e67cb-33f3-455f-a673-8bd7ef2f88de","apiVersion":"coordination.k8s.io/v1","resourceVersion":"244175"}, "reason": "LeaderElection"}
    2023-09-28T00:11:32Z    INFO    Starting EventSource    {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo", "source": "kind source: *v1alpha1.Foo"}
    2023-09-28T00:11:32Z    INFO    Starting EventSource    {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo", "source": "kind source: *v1.Deployment"}
    2023-09-28T00:11:32Z    INFO    Starting Controller     {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo"}
    2023-09-28T00:11:32Z    INFO    Starting workers        {"controller": "foo", "controllerGroup": "samplecontroller.example.com", "controllerKind": "Foo", "worker count": 1}
    ```

1. Check: apply `v1beta1` -> get `v1beta1` and `v1alpha1`

    Apply `v1beta1`

    ```yaml
    spec:
      deploymentName: "foo-sample"
    ```

    ```
    kubectl apply -f config/samples/samplecontroller_v1beta1_foo.yaml
    foo.samplecontroller.k8s.io/foo-sample created
    ```

    confirm `beta1`

    ```
    kubectl get foo.v1beta1.samplecontroller.example.com foo-sample -o jsonpath='{.spec}'
    {"deploymentName":"foo-sample","foo":"foo-sample","replicas":1}
    ```

    Check `alpha1` -> no foo

    ```
    kubectl get foo.v1alpha1.samplecontroller.example.com foo-sample -o jsonpath='{.spec}'
    {"deploymentName":"foo-sample","replicas":1}
    ```

1. Check: apply `v1alpha1` -> get `v1beta1` and `v1alpha1`

    ```
    kubectl apply -f config/samples/samplecontroller_v1alpha1_foo.yaml
    foo.samplecontroller.example.com/foo-sample created

    kubectl get foo.v1alpha1.samplecontroller.example.com foo-sample -o jsonpath='{.spec}'
    {"deploymentName":"foo-sample","replicas":2}

    kubectl get foo.v1beta1.samplecontroller.example.com foo-sample -o jsonpath='{.spec}'
    {"deploymentName":"foo-sample","foo":"foo-sample","replicas":2}
    ```

1. Clean up

    ```
    kubectl delete -f config/samples/samplecontroller_v1beta1_foo.yaml
    make undeploy
    kubectl delete -f https://github.com/jetstack/cert-manager/releases/download/v1.13.0/cert-manager.yaml
    ```
