package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/emreyalnizlar/trendyol-operator/api/v1"
)

// TrendyolApplicationReconciler reconciles a TrendyolApplication object
type TrendyolApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.trendyol.com,resources=trendyolapplications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.trendyol.com,resources=trendyolapplications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.trendyol.com,resources=trendyolapplications/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *TrendyolApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. CR'ı getir
	var app appv1alpha1.TrendyolApplication
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Deployment nesnesi oluştur
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Spec.Namespace,
		},
	}

	// 3. Deployment oluştur ya da güncelle
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, deploy, func() error {
		replicas := int32(1)
		if app.Spec.Replicas != nil {
			replicas = *app.Spec.Replicas
		}

		deploy.Spec = appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": app.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": app.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    app.Name,
							Image:   app.Spec.Image,
							Command: app.Spec.Command,
							Args:    app.Spec.Arguments,
							Env: func() []corev1.EnvVar {
								var envs []corev1.EnvVar
								for k, v := range app.Spec.Environment {
									envs = append(envs, corev1.EnvVar{Name: k, Value: v})
								}
								return envs
							}(),
						},
					},
				},
			},
		}

		// Pull secret ekleme
		if app.Spec.PullSecret != "" {
			deploy.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
				{Name: app.Spec.PullSecret},
			}
		}

		return controllerutil.SetControllerReference(&app, deploy, r.Scheme)
	})

	if err != nil {
		logger.Error(err, "unable to create or update Deployment")
		return ctrl.Result{}, err
	}

	// 4. Status Guncelleme
	app.Status.Phase = "Ready"
	app.Status.DeployedAs = deploy.Name
	app.Status.Namespace = deploy.Namespace

	if err := r.Status().Update(ctx, &app); err != nil {
		logger.Error(err, "unable to update status")
		return ctrl.Result{}, err
	}

	logger.Info("Reconciled successfully", "name", app.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TrendyolApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.TrendyolApplication{}).
		Complete(r)
}
