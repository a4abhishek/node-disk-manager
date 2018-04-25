# Integration tests
- Various tests with some examples.
-
-

## Setup
- Prepare an environment with Ubuntu-Xenial.
- Install `Go-1.9` and configure it properly.- Clone project under `$GOPATH/src/github.com/openebs/node-disk-manager`.
- Install `kubectl-1.9`
  ```
  curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.9.4/bin/darwin/amd64/kubectl # This is for kubectl 1.9.4 You can change it.
  sudo chmod +x ./kubectl
  sudo mv ./kubectl /usr/local/bin/kubectl
  ```
- Install `minikube-0.25`
  ```
  curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.25.0/minikube-linux-amd64 # This is for minikube 0.25.0 You can change it.
  sudo chmod +x minikube
  sudo mv minikube /usr/local/bin/
  ```

  ### Install test-script related dependencies:
  - `apt` should be present in system.
  - `python 2.7` should be present in system as we are currently using Python 2.7 for testing.
  - Install latest `python-pip`.
    ```
    apt install python-pip
    pip install --upgrade pip
    ```
  - Install `pyYAML`
    ```
    pip install pyYAML
    ```
  - Install `kubernetes python client`
    ```
    pip install kubernetes
    ```

## Framework to write test case
### Behavior-driven development (BDD)
- Ginkgo (Go).
- Pytest-bdd, behave, cucumber (Python)
-
### Table-driven tests (TDT)
- Some framework which support TDT.
-
-
