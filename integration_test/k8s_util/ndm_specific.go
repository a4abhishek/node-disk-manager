package k8sutil

import (
	"time"

	core_v1 "k8s.io/api/core/v1"
)

// GetNdmPod returns Pod object of node-disk-manager.
// :return: kubernetes.client.models.v1_pod.V1Pod: node-disk-manager Pod object.
func GetNdmPod() (core_v1.Pod, error) {
	// Assumption: NDM pod runs under 'dafault' namespace.
	// Assumption: Pod name starts with string 'node-disk-manager'.
	// Assumption: There is only one node-disk-manager pod
	return GetPod("default", "node-disk-manager")
}

// GetContainerStateInNdmPod returns the state of the first container of the supplied index.
//    :param waitTimeUnit: maximum time duration to get the container's state.
//                       This method does not very strictly obey this param.
//    :return: k8s.io/api/core/v1.ContainerState: state of the container.
func GetContainerStateInNdmPod(waitTimeUnit time.Duration) (core_v1.ContainerState, error) {
	ndmPod, err := GetNdmPod()
	if err != nil {
		return core_v1.ContainerState{}, err
	}
	// Assumption: There is only one container in the node-disk-manager pod
	return GetContainerStateInPod(ndmPod, 0, waitTimeUnit)
}
