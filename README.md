# PoC code to update bootimages in OpenShift 4

See https://github.com/openshift/enhancements/pull/201

## Status

Tested and working in GCP from a cluster installed as 4.5 and upgraded in place to 4.6.

Other platforms:

 - AWS: Code written but not tested yet

## Usage

Your cluster must be using an OpenShift release channel that ends in `-4.6`; for example, `stable-4.6`.

First, a good thing to do would be to create a backup of your machinesets:

```
$ oc -n openshift-machine-api get -o yaml machineset > machinesets.yaml
```

Let's create an admin service account for this, then run it as a pod:

```
$ oc project kube-system
$ oc create sa bootimage-updater
$ oc adm policy add-cluster-role-to-user cluster-admin -z bootimage-updater
$ oc run --restart=Never --serviceaccount=bootimage-updater --image=registry.svc.ci.openshift.org/cgwalters/openshift-update-bootimages update-bootimages 
```

Look at the logs:
```
$ oc logs pod/update-bootimages
```

If that succeeded, in order to test that this fully works, try [scaling up one of your machinesets](https://docs.openshift.com/container-platform/4.6/machine_management/manually-scaling-machineset.html).

## How it works

We first translate e.g. `worker-user-data` to Ignition spec 3x, creating a new `worker-user-data-ignv3`.

Then we loop over all the machinesets, and update them to use the new user-data secret plus the updated bootimage.

## Reverting/debugging

See [this FAQ entry](https://github.com/openshift/machine-api-operator/blob/master/FAQ.md#i-created-a-machine-but-it-never-joined-the-cluster) if the scaled up node doesn't join.

If you want to revert; at this time, we don't automatically back up the machinesets; hopefully you did that above.
