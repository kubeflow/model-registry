# Model Registry Deployment and Deployment Test

This deployment of model-registry is deployed via Opendatahub and used the ODH nightly images deployed to an openshfit cluster.

The script will do the following:
* Create a catalogue source pointing to the latest successful nightly ODH image.
* Deploy the opendatahub operator using the catalogue source.
* Deploy a Data Science Cluster.
* Test the model-registry-operator-contoller-manager pods are running.
* Clone the model-registry-operator repository.
* Deploy model-registry using config/samples/mysql configuration files.
* Test the model-registry-db mysql database pod is running
* Test the modelregistry-sample pods are running

## Pre-requisites:

You will need to have an openshift cluster deployed and be oc logged into you cluster as admin.

## Runing the script:

From the root of the repository
```
./openshift-ci/scripts/oc-model-registry.-deploy.sh
```

## Runing the openshift-ci

You can start the openshift-ci job to test changes in your Pull Request. To do so put the following command into a comment in your Pull Request
```
/test e2e-odh-mro-optional
```

Previous jobs can be seen [here](https://prow.ci.openshift.org/job-history/gs/test-platform-results/pr-logs/directory/rehearse-49999-pull-ci-opendatahub-io-model-registry-main-e2e-odh-mro-optional)

