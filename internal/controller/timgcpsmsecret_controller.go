package controller

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sort"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	secretsv1alpha1 "github.com/renatoruis/timgcpsm-operator/api/v1alpha1"
	gsm "github.com/renatoruis/timgcpsm-operator/internal/secretmanager"
)

// TimGcpSmSecretReconciler reconciles a TimGcpSmSecret object
type TimGcpSmSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	GSM    *gsm.Client
}

func decodeFormat(spec secretsv1alpha1.TimGcpSmSecretSpec) string {
	if spec.DecodeFormat == "json" {
		return "json"
	}
	return "text"
}

// +kubebuilder:rbac:groups=secrets.tim.operator,resources=timgcpsmsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secrets.tim.operator,resources=timgcpsmsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=secrets.tim.operator,resources=timgcpsmsecrets/finalizers,verbs=update
// +kubebuilder:rbac:groups=secrets.tim.operator,resources=timgcpsmsecretconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups=secrets.tim.operator,resources=timgcpsmclusterconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile syncs GCP Secret Manager → Kubernetes Secret and optionally restarts a Deployment when data changes.
func (r *TimGcpSmSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	obj := &secretsv1alpha1.TimGcpSmSecret{}
	err := r.Get(ctx, req.NamespacedName, obj)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("TimGcpSmSecret not found; ignoring")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get TimGcpSmSecret")
		return ctrl.Result{}, err
	}

	syncInterval := r.parseSyncInterval(obj.Spec.SyncInterval)

	projectID, err := r.resolveProjectID(ctx, obj)
	if err != nil {
		logger.Error(err, "failed to resolve GCP project")
		return r.handleError(ctx, obj, syncInterval, err, "GcpProjectResolutionFailed")
	}

	if r.GSM == nil {
		err := fmt.Errorf("Secret Manager client is not configured")
		return r.handleError(ctx, obj, syncInterval, err, "GcpSecretManagerClientMissing")
	}

	version := obj.Spec.SecretVersion
	if version == "" {
		version = "latest"
	}

	secretData, err := r.GSM.GetSecretData(ctx, projectID, obj.Spec.SecretID, version, decodeFormat(obj.Spec), obj.Spec.SecretKey)
	if err != nil {
		logger.Error(err, "failed to read secret from GCP Secret Manager")
		return r.handleError(ctx, obj, syncInterval, err, "GcpSecretFetchFailed")
	}

	newHash := calculateHash(secretData)
	oldHash := obj.Status.SecretHash

	logger.V(1).Info("hash comparison", "oldHash", oldHash, "newHash", newHash, "changed", oldHash != newHash)

	namespace := obj.Spec.Namespace
	if namespace == "" {
		namespace = obj.Namespace
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      obj.Spec.SecretName,
			Namespace: namespace,
		},
	}

	secretExists := true
	err = r.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace}, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			secretExists = false
		} else {
			logger.Error(err, "failed to get Secret")
			return ctrl.Result{}, err
		}
	}

	secretChanged := obj.Status.SecretHash != newHash

	secretDataBytes := make(map[string][]byte)
	for k, v := range secretData {
		secretDataBytes[k] = []byte(v)
	}

	if !secretExists {
		secret.Data = secretDataBytes
		secret.Type = corev1.SecretTypeOpaque
		if err := ctrl.SetControllerReference(obj, secret, r.Scheme); err != nil {
			logger.Error(err, "failed to set controller reference")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, secret); err != nil {
			logger.Error(err, "failed to create Secret")
			return ctrl.Result{}, err
		}
		logger.Info("created Secret", "name", secret.Name, "namespace", secret.Namespace)
	} else if secretChanged {
		secret.Data = secretDataBytes
		secret.Type = corev1.SecretTypeOpaque
		if err := r.Update(ctx, secret); err != nil {
			logger.Error(err, "failed to update Secret")
			return ctrl.Result{}, err
		}
		logger.Info("updated Secret", "name", secret.Name, "namespace", secret.Namespace)
	} else {
		logger.Info("Secret data unchanged, skipping update", "name", secret.Name, "namespace", secret.Namespace)
	}

	if secretChanged && obj.Spec.DeploymentName != "" {
		if err := r.restartDeployment(ctx, obj.Spec.DeploymentName, namespace); err != nil {
			logger.Error(err, "failed to restart Deployment")
			return ctrl.Result{}, err
		}
		logger.Info("restarted Deployment", "name", obj.Spec.DeploymentName, "namespace", namespace)
	}

	now := metav1.Now()
	obj.Status.LastSyncTime = &now
	obj.Status.SecretHash = newHash
	obj.Status.RetryCount = 0
	obj.Status.LastError = ""
	obj.Status.Conditions = []metav1.Condition{
		{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: now,
			Reason:             "SecretSynced",
			Message:            "Secret successfully synced from GCP Secret Manager",
		},
	}

	if err := r.Status().Update(ctx, obj); err != nil {
		logger.Error(err, "failed to update TimGcpSmSecret status")
		return ctrl.Result{}, err
	}

	logger.Info("sync ok, requeue", "syncInterval", syncInterval)
	return ctrl.Result{RequeueAfter: syncInterval}, nil
}

func (r *TimGcpSmSecretReconciler) restartDeployment(ctx context.Context, name, namespace string) error {
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, deployment)
	if err != nil {
		return fmt.Errorf("get deployment: %w", err)
	}
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["secrets.tim.operator/restartedAt"] = time.Now().Format(time.RFC3339)
	if err := r.Update(ctx, deployment); err != nil {
		return fmt.Errorf("update deployment: %w", err)
	}
	return nil
}

func (r *TimGcpSmSecretReconciler) resolveProjectID(ctx context.Context, ts *secretsv1alpha1.TimGcpSmSecret) (string, error) {
	if ts.Spec.ProjectID != "" {
		return ts.Spec.ProjectID, nil
	}
	if ts.Spec.GcpSmConfig != "" {
		ns := ts.Spec.GcpSmConfigNamespace
		if ns == "" {
			ns = ts.Namespace
		}
		cfg := &secretsv1alpha1.TimGcpSmSecretConfig{}
		if err := r.Get(ctx, types.NamespacedName{Name: ts.Spec.GcpSmConfig, Namespace: ns}, cfg); err != nil {
			return "", fmt.Errorf("get TimGcpSmSecretConfig %s/%s: %w", ns, ts.Spec.GcpSmConfig, err)
		}
		if cfg.Spec.ProjectID == "" {
			return "", fmt.Errorf("TimGcpSmSecretConfig %s/%s has empty projectId", ns, ts.Spec.GcpSmConfig)
		}
		return cfg.Spec.ProjectID, nil
	}
	if ts.Spec.ClusterConfig != "" {
		return r.projectIDFromClusterConfig(ctx, ts.Spec.ClusterConfig)
	}
	cc := &secretsv1alpha1.TimGcpSmClusterConfig{}
	if err := r.Get(ctx, types.NamespacedName{Name: secretsv1alpha1.DefaultClusterConfigName, Namespace: ""}, cc); err != nil {
		if errors.IsNotFound(err) {
			return "", fmt.Errorf("set spec.projectId, spec.gcpSmConfig, spec.clusterConfig, or create TimGcpSmClusterConfig/%s with spec.projectId", secretsv1alpha1.DefaultClusterConfigName)
		}
		return "", fmt.Errorf("get TimGcpSmClusterConfig %q: %w", secretsv1alpha1.DefaultClusterConfigName, err)
	}
	if cc.Spec.ProjectID == "" {
		return "", fmt.Errorf("TimGcpSmClusterConfig %q has empty projectId", secretsv1alpha1.DefaultClusterConfigName)
	}
	return cc.Spec.ProjectID, nil
}

func (r *TimGcpSmSecretReconciler) projectIDFromClusterConfig(ctx context.Context, name string) (string, error) {
	cc := &secretsv1alpha1.TimGcpSmClusterConfig{}
	if err := r.Get(ctx, types.NamespacedName{Name: name, Namespace: ""}, cc); err != nil {
		return "", fmt.Errorf("get TimGcpSmClusterConfig %q: %w", name, err)
	}
	if cc.Spec.ProjectID == "" {
		return "", fmt.Errorf("TimGcpSmClusterConfig %q has empty projectId", name)
	}
	return cc.Spec.ProjectID, nil
}

func calculateHash(data map[string]string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(data[k]))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (r *TimGcpSmSecretReconciler) parseSyncInterval(interval string) time.Duration {
	if interval == "" {
		return 5 * time.Minute
	}
	duration, err := time.ParseDuration(interval)
	if err != nil {
		return 5 * time.Minute
	}
	if duration < 30*time.Second {
		return 30 * time.Second
	}
	if duration > 1*time.Hour {
		return 1 * time.Hour
	}
	return duration
}

func (r *TimGcpSmSecretReconciler) handleError(ctx context.Context, ts *secretsv1alpha1.TimGcpSmSecret, syncInterval time.Duration, err error, reason string) (ctrl.Result, error) {
	ts.Status.RetryCount++
	ts.Status.LastError = err.Error()
	const maxRetries = 20
	if ts.Status.RetryCount > maxRetries {
		ts.Status.RetryCount = maxRetries
	}
	var backoff time.Duration
	exponent := ts.Status.RetryCount - 1
	if exponent > 8 {
		exponent = 8
	}
	if exponent <= 0 {
		backoff = 10 * time.Second
	} else {
		backoff = time.Duration(10*int64(time.Second)) * (1 << uint(exponent))
	}
	if backoff < 10*time.Second {
		backoff = 10 * time.Second
	}
	if backoff > syncInterval {
		backoff = syncInterval
	}
	if backoff > 5*time.Minute {
		backoff = 5 * time.Minute
	}
	ts.Status.Conditions = []metav1.Condition{
		{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             reason,
			Message:            fmt.Sprintf("Retry %d (max %d): %v", ts.Status.RetryCount, maxRetries, err),
		},
	}
	if updateErr := r.Status().Update(ctx, ts); updateErr != nil {
		return ctrl.Result{RequeueAfter: backoff}, fmt.Errorf("update status: %w (was: %v)", updateErr, err)
	}
	log.FromContext(ctx).Info("retry after error", "retryCount", ts.Status.RetryCount, "backoff", backoff, "error", err.Error())
	return ctrl.Result{RequeueAfter: backoff}, nil
}

// SetupWithManager sets up the controller
func (r *TimGcpSmSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretsv1alpha1.TimGcpSmSecret{}).
		Owns(&corev1.Secret{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)
}
