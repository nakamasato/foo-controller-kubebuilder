import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	samplecontrollerv1alpha1 "github.com/nakamasato/foo-controller-kubebuilder/api/v1alpha1"
)

// ConvertTo converts this Foo to the Hub version (v1alpha1)
func (src *Foo) ConvertTo(dstRaw conversion.Hub) error {
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

// ConvertFrom converts from the Hub version (v1alpha1) to this version
func (dst *Foo) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*samplecontrollerv1alpha1.Foo)

	// just copy DeploymentName value to Foo
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