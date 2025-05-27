package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Webhook interface'lerini elle tanımladık (alternatif olarak sigs.k8s.io/controller-runtime/pkg/webhook import edilebilir)
type defaulter interface {
	Default()
}

type validator interface {
	ValidateCreate() error
	ValidateUpdate(old runtime.Object) error
	ValidateDelete() error
}

// Webhook interface implementasyonları
var _ defaulter = &TrendyolApplication{}
var _ validator = &TrendyolApplication{}

func (r *TrendyolApplication) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

func SetupTrendyolApplicationWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&TrendyolApplication{}).
		Complete()
}

func (r *TrendyolApplication) Default() {
	if r.Spec.PullSecret == "" {
		r.Spec.PullSecret = "ty-docker-registry"
	}
}

func (r *TrendyolApplication) ValidateCreate() error {
	return r.validateTrendyolApplication()
}

func (r *TrendyolApplication) ValidateUpdate(old runtime.Object) error {
	return r.validateTrendyolApplication()
}

func (r *TrendyolApplication) ValidateDelete() error {
	return nil
}

func (r *TrendyolApplication) validateTrendyolApplication() error {
	var allErrs field.ErrorList

	if r.Spec.Image == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("image"), "image must be specified"))
	}

	if r.Spec.Replicas != nil && (*r.Spec.Replicas < 1 || *r.Spec.Replicas > 50) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("replicas"), r.Spec.Replicas, "replicas must be between 1 and 50"))
	}

	if len(r.Spec.Command) == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec").Child("command"), "command must be specified"))
	}

	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		GroupVersion.WithKind("TrendyolApplication").GroupKind(),
		r.Name, allErrs)
}
