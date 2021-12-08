package upgrade_kyma

import (
	"fmt"

	"github.com/kyma-project/control-plane/components/kyma-environment-broker/internal"
)

func ForKyma2(op internal.UpgradeKymaOperation) bool {
	fmt.Println("DEBUG_DELETE conditions runtime version", op.RuntimeVersion)
	fmt.Println("DEBUG_DELETE conditions ForKyma2", op.RuntimeVersion.MajorVersion == 2)
	return op.RuntimeVersion.MajorVersion == 2
}

func ForKyma1(op internal.UpgradeKymaOperation) bool {
	fmt.Println("DEBUG_DELETE conditions ForKyma1", op.RuntimeVersion.MajorVersion == 1)
	return op.RuntimeVersion.MajorVersion == 1
}

func WhenBTPOperatorCredentialsNotProvided(op internal.UpgradeKymaOperation) bool {
	fmt.Println("DEBUG_DELETE WhenBTPOperatorCredentialsNotProvided", op.ProvisioningParameters.ErsContext.SMOperatorCredentials == nil)
	return op.ProvisioningParameters.ErsContext.SMOperatorCredentials == nil
}

func WhenBTPOperatorCredentialsProvided(op internal.UpgradeKymaOperation) bool {
	fmt.Println("DEBUG_DELETE WhenBTPOperatorCredentialsProvided", op.ProvisioningParameters.ErsContext.SMOperatorCredentials != nil)
	return op.ProvisioningParameters.ErsContext.SMOperatorCredentials != nil
}
