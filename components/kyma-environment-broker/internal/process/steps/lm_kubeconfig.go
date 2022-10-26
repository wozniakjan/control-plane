package steps

import (
	"context"
	"fmt"
	"time"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SyncKubeconfig step ensures desired state of kubeconfig secret for lifecycle manager
type syncKubeconfig struct {
	k8sClient client.Client
	cleanup   bool
}

func SyncKubeconfig(k8sClient client.Client) syncKubeconfig {
	return syncKubeconfig{
		k8sClient: k8sClient,
	}
}

func DeleteKubeconfig(k8sClient client.Client) syncKubeconfig {
	return syncKubeconfig{
		k8sClient: k8sClient,
		cleanup:   true,
	}
}

func (_ syncKubeconfig) Name() string {
	return "Sync_Kubeconfig"
}

func (s syncKubeconfig) Run(o internal.Operation, log logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	if s.cleanup {
		return s.ensureDeleted(o, log)
	} else {
		return s.ensureExists(o, log)
	}
	return o, 0, nil
}

func (s syncKubeconfig) ensureDeleted(o internal.Operation, log logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	secret := initSecret(o)
	if err := s.k8sClient.Delete(context.Background(), secret); err != nil && !errors.IsNotFound(err) {
		log.Errorf("failed to delete kubeconfig secret %v/%v for lifecycle manager: %v", secret.Namespace, secret.Name, err)
		return o, time.Minute, nil
	}
	return o, 0, nil
}

func (s syncKubeconfig) ensureExists(o internal.Operation, log logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	secret := initSecret(o)
	if err := s.k8sClient.Create(context.Background(), secret); errors.IsAlreadyExists(err) {
		if err := s.k8sClient.Update(context.Background(), secret); err != nil {
			log.Errorf("failed to update kubeconfig secret %v/%v for lifecycle manager: %v", secret.Namespace, secret.Name, err)
			return o, time.Minute, nil
		}
	} else if err != nil {
		log.Errorf("failed to create kubeconfig secret %v/%v for lifecycle manager: %v", secret.Namespace, secret.Name, err)
		return o, time.Minute, nil
	}
	return o, 0, nil
}

func initSecret(o internal.Operation) *corev1.Secret {
	// TODO: define common things such as namespace and labels with resource kyma.operator.kyma-project.io/v1alpha1 in one place
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "kcp-system",
			Name:      fmt.Sprintf("kubeconfig-%v", o.ShootName), // TODO: consider something else than shoot
			Labels: map[string]string{
				"kyma-project.io/instance-id":        o.InstanceID,
				"kyma-project.io/runtime-id":         o.RuntimeID,
				"kyma-project.io/broker-plan-id":     o.ProvisioningParameters.PlanID,
				"kyma-project.io/global-account-id":  o.GlobalAccountID,
				"operator.kyma-project.io/kyma-name": o.ShootName, // TODO: sync with kyma resource naming
			},
		},
		StringData: map[string]string{
			"config": o.Kubeconfig,
		},
	}
}

// NOTE: adapter for upgrade_kyma which is currently not using shared staged_manager
type syncKubeconfigUpgradeKyma struct {
	syncKubeconfig
}

func SyncKubeconfigUpgradeKyma(k8sClient client.Client) syncKubeconfigUpgradeKyma {
	return syncKubeconfigUpgradeKyma{SyncKubeconfig(k8sClient)}
}

func (s syncKubeconfigUpgradeKyma) Run(o internal.UpgradeKymaOperation, logger logrus.FieldLogger) (internal.UpgradeKymaOperation, time.Duration, error) {
	o2, w, err := s.syncKubeconfig.Run(o.Operation, logger)
	return internal.UpgradeKymaOperation{o2}, w, err
}
