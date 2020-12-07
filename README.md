# PoC code to update bootimages in OpenShift 4

See https://github.com/openshift/enhancements/pull/201

## Status

Tested and working in GCP from a cluster installed as 4.5 and upgraded in place to 4.6.

Other platforms:

 - AWS: Code written but not tested yet

## Usage

Let's create an admin service account for this, then run it as a pod:

```
$ oc project kube-system
$ oc create sa bootimage-updater
$ oc adm policy add-cluster-role-to-user cluster-admin -z bootimage-updater
$ oc run --serviceaccount=bootimage-updater  --image=registry.svc.ci.openshift.org/cgwalters/openshift-update-bootimages update-bootimages
```
