package kops

import (
	"fmt"

	kopsapi "k8s.io/kops/pkg/apis/kops"

	"github.com/pkg/errors"
	kcontrolplanev1alpha1 "github.com/topfreegames/kubernetes-kops-operator/apis/controlplane/v1alpha1"
	kinfrastructurev1alpha1 "github.com/topfreegames/kubernetes-kops-operator/apis/infrastructure/v1alpha1"
)

func GetSubnetFromKopsControlPlane(kcp *kcontrolplanev1alpha1.KopsControlPlane) (*kopsapi.ClusterSubnetSpec, error) {
	if kcp.Spec.KopsClusterSpec.Subnets == nil {
		return nil, errors.Wrap(errors.Errorf("SubnetNotFound"), "subnet not found in KopsControlPlane")
	}
	subnet := kcp.Spec.KopsClusterSpec.Subnets[0]
	return &subnet, nil
}

func GetRegionFromKopsSubnet(subnet kopsapi.ClusterSubnetSpec) (*string, error) {
	if subnet.Region != "" {
		return &subnet.Region, nil
	}

	if subnet.Zone != "" {
		zone := subnet.Zone
		region := zone[:len(zone)-1]
		return &region, nil
	}

	return nil, errors.Wrap(errors.Errorf("RegionNotFound"), "couldn't get region from KopsControlPlane")
}

func GetAutoScalingGroupNameFromKopsMachinePool(kmp kinfrastructurev1alpha1.KopsMachinePool) (*string, error) {
	if _, ok := kmp.Spec.KopsInstanceGroupSpec.NodeLabels["kops.k8s.io/instance-group-name"]; !ok {
		return nil, fmt.Errorf("failed to retrieve igName from KopsMachinePool %s", kmp.GetName())
	}

	if kmp.Spec.ClusterName == "" {
		return nil, fmt.Errorf("failed to retrieve clusterName from KopsMachinePool %s", kmp.GetName())
	}

	asgName := fmt.Sprintf("%s.%s", kmp.Spec.KopsInstanceGroupSpec.NodeLabels["kops.k8s.io/instance-group-name"], kmp.Spec.ClusterName)

	return &asgName, nil
}
