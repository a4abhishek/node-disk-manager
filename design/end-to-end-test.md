# Introduction

## Purpose
This test plan describes the testing approach and overall framework that will drive the testing of the **node-disk-manager**. The document introduces:
- **Test Strategy:** rules the test will be based on, including the givens of the project (e.g.: objectives, assumptions); description of the process to set up a valid test (e.g.: entry / exit criteria, creation of test cases, specific tasks to perform).
- **Execution Strategy:** describes how the test will be performed and process to identify and report defects, and to fix and implement fixes.
- **Test Management:** process to handle the logistics of the test and all the events that come up during execution (e.g.: communications, escalation procedures, risk and mitigation)

## Project Overview
**node-disk-manager** that will automate the management of the disks attached to the node. There will be new Custom Resources added to Kubernetes to represent the underlying storage infrastructure like - Disks and StoragePools and their associated Claims - DiskClaim, StoragePoolClaim.

The disks could be of any type ranging from local disks (ssds or spinning media), NVMe/PCI based or disks coming from external SAN/NAS. A new Kubernetes Custom Resource called - Disk Object will be added that will represent the disks attached to the node.

The **node-disk-manager** can also help in creating Storage Pools using one or more Disk Objects. The Storage Pools can be of varying types starting from a plain ext4 mounts, lvm, zfs to ceph, glusterfs, cstor pools.

**node-disk-manager** is intended to:

- be extensible in terms of adding support for different kinds of Disk Objects and Storage Pools.
- act as a bridge to augment the functionality already provided by the kubelet or other infrastructure containers running on the node. For example, node-bot will have to interface with node-problem-detector or heapster etc.,
- be considered as a Infrastructure Pod (or a daemonset) that runs on each of the Kubernetes node or the functionality could be embedded into other Infrastructure Pods, similar to node-problem-detector.

<!-- Audience: -->

## Test Strategy
### Test Objectives
The objective of the test is to verify that the functionality of **node-disk-manager** works according to the specifications.

The test will execute and verify the test scripts, identify and fix errors on the code or test cases.

The final product of the test is twofold: 
- A well-tested software.
- A set of stable test scripts that can be reused for further test execution.

### Test Assumptions
- Yaml for NDM daemonset: 'node-disk-manager.yaml' is present in the root directory of the project.
- Docker image name: 'openebs/node-disk-manager'.
- Docker image tag is output of `git rev-parse --abbrev-ref HEAD`.
- If 'Started controller' is in log, log is considered OK.
- **node-disk-manager** runs under 'dafault' namespace.
- Pod name of **node-disk-manager** starts with string 'node-disk-manager'.
- This pod has only one container.
- Environment variables `USER` and `HOME` is well defined.
- Additional test scripts related assumptions.

### Test Principles
- Testing will be focused on meeting the well written, high quality, bug free and efficient code.
- There will be common, consistent procedures for all teams supporting testing activities.
- Testing processes will be well defined, yet flexible, with the ability to change as needed.
- Testing environment and data will emulate a production environment as much as possible.
- Testing will be a repeatable, quantifiable, and measurable activity.
- Testing will be divided into distinct modules, each with clearly defined objectives.
<!-- -There will be entrance and exit criteria -->

## Execution Strategy
### Entry and Exit Criteria

### Creation of test cases

### Specific tasks to perform
All tasks are automated by a python script for now. The script performs following tasks:
- Runs minikube with this command: `sudo minikube start --vm-driver=none --feature-gates=MountPropagation=true`.
- In YAML changes value of `image` to match with newly built image and `imagePullPolicy` to `IfNotPresent`.
- Apply the YAML.
- Check pod description and log.
- Runs `lsblk` and `ndm` in both inside and outside the pod and match the results.

## Test Management
### Test Design Process
- The tester will understand each requirement and prepare corresponding test case to ensure all requirements are covered.
- Each of the Test cases will undergo review an architect or Project manager. The testers will rework on the review defects and finally obtain approval.
- Tester should automate the test by updating the aforementioned test-script.
- Test cases / aforementioned test-script changes should be signed-off.

### Test Execution Process
- A ci environment with necessary resources and accesses should be prepared.
- Aforementioned test-script should be run on a ci tool and checked for any failures.
- In case of any failures, testers should correct the code and/or test case to match 
