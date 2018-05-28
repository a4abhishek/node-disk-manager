package ndmutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/openebs/node-disk-manager/integration_test/minikube_adm"

	"github.com/openebs/node-disk-manager/integration_test/k8s_util"

	"io/ioutil"

	"github.com/golang/glog"
	. "github.com/openebs/node-disk-manager/integration_test/common"
	core_v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
)

const (
	// NdmYAML is the name of Configuration file for node-disk-manager
	NdmYAML = "node-disk-manager.yaml"
	// NdmTestYAMLPath is the directory name where we shall keep temporary configuration file
	NdmTestYAMLPath = "/tmp/"
	// NdmTestYAMLName is the name of temporary configuration file
	NdmTestYAMLName = "NDM_Test_" + NdmYAML
	// NdmNamespace is the namespace of the node-disk-manager
	NdmNamespace = "default"
)

var (
	// MaxTry is the number of tries that the functions in this package make. (Not all functions)
	MaxTry int = 5
	// WaitTimeUnit is the unit time duration that is mostly used by the functions in this package
	// to wait for after a try (Not always applicable)
	WaitTimeUnit time.Duration = 1 * time.Second
)

// GetNDMDir returns the path to the node-disk-manager repository
func GetNDMDir() string {
	return path.Join(os.Getenv("GOPATH"), "src/github.com/openebs/node-disk-manager/")
}

// GetNDMBinDir returns the path of bin folder of node-disk-manager repository
func GetNDMBinDir() string {
	return path.Join(GetNDMDir(), "bin/")
}

// GetDockerImageName returns the docker image name of node-disk-manager
// that will build when we build the project
func GetDockerImageName() string {
	return "openebs/node-disk-manager-" + GetenvFallback("XC_ARCH", "amd64")
}

// GetDockerImageTag returns the docker tag of the node-disk-manager docker image
// that wiil be used when we build the project
func GetDockerImageTag() string {
	tag, err := ExecCommand("git describe --tags --always")
	if err != nil {
		fmt.Printf("Error while getting docker tag name. Error: %+v\n", err)
		return ""
	}
	return strings.TrimSpace(tag)
}

// GetNDMConfigurationFileName returns value of the constant NdmYAML
// which is the name of Configuration file for node-disk-manager
func GetNDMConfigurationFileName() string {
	return NdmYAML
}

// GetNDMTestConfigurationDir returns value of the constant NdmTestYAMLPath
// which is the directory name where we shall keep temporary configuration file
func GetNDMTestConfigurationDir() string {
	return NdmTestYAMLPath
}

// GetNDMTestConfigurationFileName returns value of the constant NdmTestYAMLName
// which is the name of temporary configuration file
func GetNDMTestConfigurationFileName() string {
	return NdmTestYAMLName
}

// GetNDMNamespace returns value of the constant NdmNamespace
// which is the namespace of the node-disk-manager
func GetNDMNamespace() string {
	return NdmNamespace
}

// GetNDMConfigurationFilePath returns the path of Configuration file for node-disk-manager
// i.e. GetNDMDir() and GetNDMConfigurationFileName() joined together
func GetNDMConfigurationFilePath() string {
	return path.Join(GetNDMDir(), GetNDMConfigurationFileName())
}

// GetNDMTestConfigurationFilePath returns the path of Test Configuration file for node-disk-manager
// i.e. GetNDMTestConfigurationDir() and GetNDMTestConfigurationFileName() joined together
func GetNDMTestConfigurationFilePath() string {
	return path.Join(GetNDMTestConfigurationDir(), GetNDMTestConfigurationFileName())
}

// TODO: Check the pod current status like we do in `kubectl describe`
// Example: Check if all volumes are mounted correctly
// func ValidateNDMPod(ndmPod v1.Pod) bool {}

// ValidateNdmLog checks the supplied log and checks for any error in the log.
// :param string log: log of node-disk-manager-xxx Pod
// :return: bool: `true` if log contains no error otherwise return `false`.
func ValidateNdmLog(log string) bool {
	// Definition of this function should grow as program grows
	if strings.Contains(log, "Started the controller") {
		return true
	}
	return false
}

// GetNDMLogAndValidate extracts log of node-disk-manager and then validate
// return the validation status and error occured during process
func GetNDMLogAndValidate() (bool, error) {
	// Getting the log
	ndmPod, err := k8sutil.GetNdmPod()
	if err != nil {
		return false, err
	}

	log, err := k8sutil.GetLog(ndmPod.Name, ndmPod.Namespace)
	if err != nil {
		return false, err
	}

	return ValidateNdmLog(log), nil
}

// YAMLPrepare reads and parse the configuration file into v1beta1.DaemonSet and changes some fields
// so that it can be applied to create node-disk-manager daemonset from recently built image.
// Then it returns that v1beta1.DaemonSet Structure
func YAMLPrepare() (v1beta1.DaemonSet, error) {
	// Prepare the yaml
	yamlBytes, err := ioutil.ReadFile(GetNDMConfigurationFilePath())
	if err != nil {
		return v1beta1.DaemonSet{}, err
	}

	// Get the DaemonSet Struct
	ds, err := k8sutil.GetDaemonsetStructFromYamlBytes(yamlBytes)
	if err != nil {
		return v1beta1.DaemonSet{}, err
	}

	// Assumption: In following two statements it is assumed that
	// this pod has only one container

	// Set image name
	fmt.Println("Image name:", GetDockerImageName()+":"+GetDockerImageTag()) // Remove me
	ds.Spec.Template.Spec.Containers[0].Image = GetDockerImageName() + ":" + GetDockerImageTag()

	// set imagePullPolicy
	ds.Spec.Template.Spec.Containers[0].ImagePullPolicy = "IfNotPresent"

	return ds, nil
}

// YAMLPrepareAndWriteInTempPath reads and parse the configuration file and changes some fields
// so that it can be applied to create node-disk-manager daemonset from recently built image.
// Then it saves that configuration to temp directory which will be cleaned by calling Clean()
func YAMLPrepareAndWriteInTempPath() error {
	dsManifest, err := YAMLPrepare()
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(dsManifest)
	if err != nil {
		return err
	}

	yamlBytes, err := ConvertJSONtoYAML(jsonBytes)
	if err != nil {
		return err
	}
	ioutil.WriteFile(GetNDMTestConfigurationFilePath(), yamlBytes, 0644)
	return nil
}

// PrepaeAndApplyYAML prepares and applies the node-disk-manager configuration
func PrepaeAndApplyYAML() error {
	dsManifest, err := YAMLPrepare()
	if err != nil {
		return err
	}

	_, err = k8sutil.ApplyDSFromManifestStruct(dsManifest)
	return err
}

// MatchDisksOutsideAndInside takes output of `lsblk -bro name,size,type,mountpoint` on the host
// and output of `ndm device list` inside the pod. Then it matches the two
func MatchDisksOutsideAndInside() (bool, error) {
	lsblkOutput, err := ExecCommand("lsblk -bro name,size,type,mountpoint")
	if err != nil {
		return false, err
	}

	// Remove Me
	fmt.Println("`lsblk` output:")
	fmt.Println(lsblkOutput)

	ndmPod, err := k8sutil.GetNdmPod()
	if err != nil {
		return false, err
	}

	ndmOutput, err := k8sutil.ExecToPod("ndm device list", ndmPod.Name, ndmPod.Namespace)
	if err != nil {
		return false, err
	}

	// Remove Me
	fmt.Println("`ndm device list` output:")
	fmt.Println(ndmOutput)

	pattern := regexp.MustCompile("[^a-zA-Z0-9./\\[\\]]+")
	lsblkOutput = pattern.ReplaceAllString(lsblkOutput, "")

	ndmOutput = pattern.ReplaceAllString(ndmOutput, "")
	ndmOutput = strings.Replace(ndmOutput, "null", "", -1)
	ndmOutput = strings.Replace(ndmOutput, "/etc/hosts/", "/", -1)
	ndmOutput = strings.Replace(ndmOutput, "/etc/hosts", "/", -1)

	fmt.Println("After conversion:")
	fmt.Println("lsblk:", strings.ToLower(lsblkOutput))
	fmt.Println("ndm:", strings.ToLower(ndmOutput))
	return strings.ToLower(lsblkOutput) == strings.ToLower(ndmOutput), nil
}

// WaitTillDefaultNSisReady busy waits until default namespace is ready
// or number of try exceeds MaxTry (at least once).
// Each try is make after waiting for time period of WaitTimeUnit
func WaitTillDefaultNSisReady() {
	// Ensuring minimum value for the arguments

	// Making sure that it tries at least once to apply the YAML
	maxTry := MaxTry
	if maxTry < 1 {
		maxTry = 1
	}
	waitTimeUnit := WaitTimeUnit
	if waitTimeUnit < 0*time.Second {
		waitTimeUnit = 0 * time.Second
	}

	ndmNS := core_v1.Namespace{}
	i := 0
	for i < maxTry {
		namespaces, err := k8sutil.GetAllNamespacesCoreV1NamespaceArray()
		if err == nil {
			for _, namespace := range namespaces {
				if namespace.Name == GetNDMNamespace() {
					ndmNS = namespace
					break
				}
			}
		} else {
			fmt.Printf("Try: %d - Error getting namespaces. Error: %+v\n", i, err)
			time.Sleep(WaitTimeUnit)
			continue
		}

		// If namespace is not even created
		if reflect.DeepEqual(ndmNS, core_v1.Namespace{}) {
			fmt.Printf("Try: %d - Waiting as Namespace %q has not been created yet.\n", i, GetNDMNamespace())
			time.Sleep(WaitTimeUnit)
			continue
		}

		if IsNSinGoodPhase(ndmNS) {
			return
		}
		fmt.Printf("Try: %d - Waiting as Namespace %q is in %q phase\n", i, ndmNS.Name, ndmNS.Status.Phase)
	}

	// Final Check
	if reflect.DeepEqual(ndmNS, core_v1.Namespace{}) {
		glog.Fatalf("Namespace %q didn't came up in %v", GetNDMNamespace(), time.Duration(MaxTry)*WaitTimeUnit)
	} else if !IsNSinGoodPhase(ndmNS) {
		glog.Fatalf("Namespace %q is still in bad phase: %q after %v", GetNDMNamespace(), ndmNS.Status.Phase, time.Duration(MaxTry)*WaitTimeUnit)
	}
	fmt.Printf("Namespace %q is ready\n", ndmNS.Name)
}

// WaitTillNDMisUp busy waits until the pod is up or number of try exceeds MaxTry.
// Each try is make after waiting for WaitTimeUnit number of seconds.
func WaitTillNDMisUp() {
	// Ensuring minimum value for the arguments

	// Making sure that it tries at least once to apply the YAML
	maxTry := MaxTry
	if maxTry < 1 {
		maxTry = 1
	}
	waitTimeUnit := WaitTimeUnit
	if waitTimeUnit < 0*time.Second {
		waitTimeUnit = 0 * time.Second
	}

	var err error
	ndmPod := core_v1.Pod{}
	podState := core_v1.ContainerState{}
	i := 0
	for ; i < maxTry; i++ { // Since we have many continue statements so we have to increment here only
		ndmPod, err = k8sutil.GetNdmPod()
		if err != nil {
			fmt.Printf("Try: %d - Error getting NDM pod. Error: %+v\n", i, err)
			time.Sleep(WaitTimeUnit)
			continue
		}

		podState, err = k8sutil.GetContainerStateInNdmPod(1 * time.Minute)
		if err != nil {
			fmt.Printf("Try: %d - Error getting container state of NDM pod. Error: %+v\n", i, err)
			time.Sleep(WaitTimeUnit)
			continue
		}

		if podState.Terminated != nil {
			glog.Fatalf("Pod terminated unexpectedly, Reason: %q. Details: %+v", podState.Terminated.Reason, podState.Terminated)
		}

		if podState.Waiting != nil {
			if IsPodStateWait(podState.Waiting.Reason) {
				fmt.Printf("Waiting as pod-state: %q. Details: %+v\n", podState.Waiting.Reason, *podState.Waiting)
				time.Sleep(WaitTimeUnit)
				continue
			} else if !IsPodStateGood(podState.Waiting.Reason) {
				glog.Fatalf("Pod is in bad state: %q. Details: %+v", podState.Waiting.Reason, *podState.Waiting)
			}
		}

		if podState.Running == nil {
			// At this point all states are None,
			// so just showing phase is enough
			fmt.Printf("Waiting as pod-phase: %q\n", k8sutil.GetPodPhase(ndmPod))
			time.Sleep(WaitTimeUnit)
		} else {
			break
		}
	}

	// Final Check
	if reflect.DeepEqual(ndmPod, core_v1.Pod{}) {
		glog.Fatalf("NDM-Pod didn't came up in %v", time.Duration(MaxTry)*WaitTimeUnit)
	} else if podState.Running == nil {
		glog.Fatalf("Pod %q is still not in \"Running\" state after %v", ndmPod.Name, time.Duration(MaxTry)*WaitTimeUnit)
	}
	fmt.Printf("Pod %q is up.\n", ndmPod.Name)
}

// Clean is intended to clean the residue of the testing.
// It should be run at the very end of the test.
// CAUTION: it calls `minikubeadm.ClearContainers`
// which removes all Docker Containers in your machine.
func Clean() {
	// Check minikube status and delete if minikube is running
	fmt.Print("Checking minikube status...")
	minikubeStatus, err := minikubeadm.CheckStatus()
	if err != nil {
		fmt.Printf("Error occured while checking status of minikube. Error: %+v\n", err)
	}
	if state, ok := minikubeStatus["minikube"]; ok && (state == "Running" || state == "Stopped") {
		fmt.Print("Deleting minikube...")
		err = minikubeadm.Teardown()
		if err != nil {
			fmt.Printf("Error while deleting minikube. Error: %+v\n", err)
		}
	} else {
		fmt.Print("Machine not present.")
	}

	// Remove all docker containers
	fmt.Print("Removing docker containers...")
	err = minikubeadm.ClearContainers()
	if err != nil {
		fmt.Printf("Error occured when deleting containers. Error: %+v\n", err)
	}

	fmt.Printf("Removing %q...\n", GetNDMTestConfigurationFileName())
	if _, err = os.Stat(GetNDMTestConfigurationFilePath()); os.IsNotExist(err) {
		fmt.Printf("%q not present\n", GetNDMTestConfigurationFileName())
	} else {
		err = os.Remove(GetNDMTestConfigurationFilePath())
		if err != nil {
			fmt.Printf("Error occured while removing NDM's temporary configuration file. Error: %+v\n", err)
		} else {
			fmt.Printf("%q removed.\n", GetNDMTestConfigurationFilePath())
		}
	}
}
