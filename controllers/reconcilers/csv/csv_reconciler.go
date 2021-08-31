package csv

import (
	"context"
	alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"github.com/RHEcosystemAppEng/dbaas-operator/controllers/reconcilers"
	"github.com/go-logr/logr"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type Reconciler struct {
	client client.Client
	logger logr.Logger
}

func NewReconciler(client client.Client, logger logr.Logger) reconcilers.PlatformReconciler {
	return &Reconciler{
		client: client,
		logger: logger,
	}
}

func (r *Reconciler) Reconcile(ctx context.Context, cr *alpha1.DBaaSPlatform, status *alpha1.DBaaSPlatformStatus) (alpha1.PlatformsInstlnStatus, error) {

	return alpha1.ResultSuccess, nil
}

func (r *Reconciler) Cleanup(ctx context.Context, cr *alpha1.DBaaSPlatform) (alpha1.PlatformsInstlnStatus, error) {

	list := &v1alpha1.ClusterServiceVersionList{}
	opts := &client.ListOptions{
		Namespace: reconcilers.INSTALL_NAMESPACE,
	}
	err := r.client.List(ctx, list, opts)
	if err != nil && !errors.IsNotFound(err) {
		return alpha1.ResultFailed, err
	}

	for _, csv := range list.Items {
		if csv.Namespace == reconcilers.INSTALL_NAMESPACE && strings.HasPrefix(csv.Name, "crunchy-bridge-operator.") {
			err := r.client.Delete(ctx, &csv)
			if err != nil && !errors.IsNotFound(err) {
				return alpha1.ResultFailed, err
			}
			return alpha1.ResultInProgress, nil
		} else if csv.Namespace == reconcilers.INSTALL_NAMESPACE && strings.HasPrefix(csv.Name, "mongodb-atlas-kubernetes.") {
			err := r.client.Delete(ctx, &csv)
			if err != nil && !errors.IsNotFound(err) {
				return alpha1.ResultFailed, err
			}
			return alpha1.ResultInProgress, nil
		}
	}

	return alpha1.ResultSuccess, nil
}

func GetDBaaSOperatorCSV(namespace string, ctx context.Context, serverClient k8sclient.Client) (*v1alpha1.ClusterServiceVersion, error) {

	list := &v1alpha1.ClusterServiceVersionList{}
	opts := &client.ListOptions{
		Namespace: namespace,
	}
	err := serverClient.List(ctx, list, opts)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	for _, csv := range list.Items {
		// Operator CSV
		if csv.Namespace == namespace && strings.HasPrefix(csv.Name, "dbaas-operator.") {
			return &csv, nil

		}
	}

	return nil, nil
}
