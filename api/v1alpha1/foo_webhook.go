/*
Copyright 2023 nakamasato.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var foolog = logf.Log.WithName("foo-resource")

func (r *Foo) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-samplecontroller-example-com-v1alpha1-foo,mutating=true,failurePolicy=fail,sideEffects=None,groups=samplecontroller.example.com,resources=foos,verbs=create;update,versions=v1alpha1,name=mfoo.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Foo{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Foo) Default() {
	foolog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-samplecontroller-example-com-v1alpha1-foo,mutating=false,failurePolicy=fail,sideEffects=None,groups=samplecontroller.example.com,resources=foos,verbs=create;update,versions=v1alpha1,name=vfoo.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Foo{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateCreate() (admission.Warnings, error) {
	foolog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	foolog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateDelete() (admission.Warnings, error) {
	foolog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
