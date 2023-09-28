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

func (dst *Foo) ConvertFrom(srcRaw conversion.Hub) error { // v1alpha1 -> v1beta1
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
