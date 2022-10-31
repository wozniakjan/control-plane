package steps

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
)

// Set common labels for kyma lifecycle manager
func ApplyLabelsForLM(object client.Object, operation internal.Operation) {
	l := object.GetLabels()
	if l == nil {
		l = make(map[string]string)
	}
	l["kyma-project.io/instance-id"] = operation.InstanceID
	l["kyma-project.io/runtime-id"] = operation.RuntimeID
	l["kyma-project.io/broker-plan-id"] = operation.ProvisioningParameters.PlanID
	l["kyma-project.io/global-account-id"] = operation.GlobalAccountID
	l["operator.kyma-project.io/kyma-name"] = KymaName(operation)
	object.SetLabels(l)
}

func KymaKubeconfigName(operation internal.Operation) string {
	return fmt.Sprintf("kubeconfig-%v", operation.ShootName)
}

func KymaName(operation internal.Operation) string {
	return operation.ShootName
}
