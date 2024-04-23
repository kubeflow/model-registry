#!/bin/bash

# Define variables for ODH deployment deployment
OPENDATAHUB_SUBSCRIPTION="openshift-ci/resources/opendatahub-subscription.yaml"
DSC_INITIALIZATION_MANIFEST="openshift-ci/resources/model-registry-DSCInitialization.yaml"
DATA_SCIENCE_CLUSTER_MANIFEST="openshift-ci/resources/opendatahub-data-science-cluster.yaml"
MODEL_REGISTRY_DB_MANIFEST="openshift-ci/resources/model-registry-operator/mysql-db.yaml"
MODEL_REGISTRY_SAMPLE_MANIFEST="openshift-ci/resources/model-registry-operator/modelregistry_v1alpha1_modelregistry.yaml"
OPENDATAHUB_CRDS="datascienceclusters.datasciencecluster.opendatahub.io,dscinitializations.dscinitialization.opendatahub.io,featuretrackers.features.opendatahub.io"
DATA_SCIENCE_CLUSTER_CRDS="acceleratorprofiles.dashboard.opendatahub.io,datasciencepipelinesapplications.datasciencepipelinesapplications.opendatahub.io,odhapplications.dashboard.opendatahub.io,odhdashboardconfigs.opendatahub.io,odhdocuments.dashboard.opendatahub.io"
MODEL_REGISTRY_CRDS="modelregistries.modelregistry.opendatahub.io"
source "openshift-ci/scripts/colour_text_variables.sh"

# Function to monitor CRDS creation and deployment.
# The function takes two arguments, reference to manifest and a wait time in seconds.
monitoring_crds_installation() {
    IFS=',' read -ra crds_array <<< "$1"
    local timeout=$2

    echo "Monitoring the installation of CRDs: ${crds_array[*]}"
    echo "Timeout set to ${timeout}s"

    local start_time=$(date +%s)

    # Loop until all specified CRDs are installed or timeout is reached
    while true; do
        local elapsed_time=$(($(date +%s) - start_time))

        # Check if timeout has been reached
        if [ "$elapsed_time" -ge "$timeout" ]; then
            echo -e "${RED}X Error:${NC} Timeout reached. Installation of CRDs failed."
            return 1
        fi

        # Get the list of installed CRDs
        local installed_crds=($(oc get crd -o=name | cut -d '/' -f2))

        # Check if all CRDs are installed
        local all_installed=true
        for crd in "${crds_array[@]}"; do
            if ! [[ " ${installed_crds[@]} " =~ " ${crd} " ]]; then
                all_installed=false
                break
            fi
        done

        # If all CRDs are installed, break out of the loop
        if [ "$all_installed" = true ]; then
            echo -e "${GREEN}✔ Success:${NC} All specified CRDs are installed."
            return 0
        fi

        # Print the status of each CRD
        for crd in "${crds_array[@]}"; do
            if [[ " ${installed_crds[@]} " =~ " ${crd} " ]]; then
                echo "CRD '$crd' is installed."
            else
                echo "CRD '$crd' is not installed."
            fi
        done

        # Wait for a few seconds before checking again
        sleep 5
    done
}

# Function to deploy and wait for deployment
# The function takes two arguments, reference to manifest and a wait time in seconds.
deploy_and_wait() {
    local namespace=$1
    local manifest=$2
    local resource_name=$(basename -s .yaml $manifest)
    local wait_time=$3
    
    sleep $wait_time
    
    echo "Deploying $resource_name from $manifest..."

    if oc apply -f $manifest $namespace --wait=true --timeout=300s; then
        echo -e "${GREEN}✔ Success:${NC} Deployment of $resource_name succeeded."
    else
        echo -e "${RED}X Error:${NC} Deployment of $resource_name failed or timed out." >&2
        return 1
    fi
}

check_deployment_availability() {
    local namespace="$1"
    local deployment="$2"
    local timeout=300  # Timeout in seconds
    local start_time=$(date +%s)

    # Loop until timeout
    while (( $(date +%s) - start_time < timeout )); do
        # Get the availability status of the deployment
        local deployment_status=$(oc get deployment "$deployment" -n "$namespace" --no-headers -o custom-columns=:.status.availableReplicas)

        # Check if the deployment is available
        if [[ $deployment_status != "" ]]; then
            echo -e "${GREEN}✔ Success:${NC} Deployment $deployment is available"
            return 0  # Success
        fi

        sleep 5  # Wait for 5 seconds before checking again
    done

    echo -e "${RED}X Error:${NC}  Timeout reached. Deployment $deployment did not become available within $timeout seconds"
    return 1  # Failure
}

# Function to check the status of deploying pods
# The function takes three arguments, namespace, descriptor to identify the component and number of containers expected.
check_pod_status() {
    local namespace="$1"
    local pod_selector="$2"
    local expected_ready_containers="$3"
    local timeout=300  # Timeout in seconds
    local start_time=$(date +%s)
    
    # Loop until timeout
    while (( $(date +%s) - start_time < timeout )); do
        # Get the list of pods in the specified namespace matching the provided partial names
        local pod_list=$(oc get pods -n $namespace $pod_selector --no-headers -o custom-columns=NAME:.metadata.name)
        
        # Iterate over each pod in the list
        while IFS= read -r pod_name; do
            # Get the pod info
            local pod_info=$(oc get pod "$pod_name" -n "$namespace" --no-headers)

            # Extract pod status and ready status from the info
            local pod_name=$(echo "$pod_info" | awk '{print $1}')
            local pod_status=$(echo "$pod_info" | awk '{print $3}')
            local pod_ready=$(echo "$pod_info" | awk '{print $2}')
            local ready_containers=$(echo "$pod_ready" | cut -d'/' -f1)

            # Check if the pod is Running and all containers are ready
            if [[ $pod_status == "Running" ]] && [[ $ready_containers -eq $expected_ready_containers ]]; then
                echo -e "${GREEN}✔ Success:${NC} Pod $pod_name is running and $ready_containers out of $expected_ready_containers containers are ready"
                return 0  # Success
            else
                echo -e "${YELLOW}! INFO:${NC}  Pod $pod_name is not running or does not have $expected_ready_containers containers ready"
            fi
        done <<< "$pod_list"

        sleep 5  # Wait for 5 seconds before checking again
    done

    echo -e "${RED}X Failure:${NC} Timeout reached. No pod matching '$pod_name_partial' became ready within $timeout seconds"
    return 1  # Failure
}

# Function to check the status of a route
# The function takes two arguments, namespace and route name.
check_route_status() {
    local namespace="$1"
    local route_name="$2"
    local key="items"
    local interval=5
    local timeout=300
    local start_time=$(date +%s)

    while (( $(date +%s) - start_time < timeout )); do
        # Get the route URL
        local route=$(oc get route -n "$namespace" "$route_name" -o jsonpath='{.spec.host}')
        local route_url="http://$route"

        if [[ -z "$route_url" ]]; then
             echo -e "${RED}X Error:${NC}  Route '$route_name' does not exist in namespace '$namespace'"
            return 1
        else 
            echo -e "${GREEN}✔ Success:${NC} Route '$route_name' exists in namespace '$namespace'"
        fi

        # Test if the route is live
        local response=$(curl -s -o /dev/null -w "%{http_code}" "$route_url/api/model_registry/v1alpha3/registered_models")

        # Check if the response status code is 200 OK or 404 Not Found
        if [[ "$response" == "200" ]]; then
            echo -e "${GREEN}✔ Success:${NC} Route server is reachable. Status code: 200 OK"
            return 0
        elif [[ "$response" == "404" ]]; then
            echo -e "${YELLOW}! WARNING:${NC} Route server is reachable, but service is not. Status code: 404 Not Found"
            return 0
        else
            echo -e "${RED}X Error:${NC}  Route server is unreachable. Status code: $response"
        fi

        sleep "$interval"
    done

    echo -e "${RED}X Error:${NC}  Timeout reached. Route '$route_name' did not become live within $timeout seconds."
    return 1
}

# Function to source the rest api tests and run them.
run_api_tests() {
    ./test/scripts/rest.sh "-n opendatahub"
}

# Run the deployment tests.
run_deployment_tests() {
    check_deployment_availability "opendatahub" model-registry-db
    check_deployment_availability "opendatahub" modelregistry-sample
    check_pod_status "opendatahub" "-l name=model-registry-db" 1
    check_pod_status "opendatahub" "-l app=modelregistry-sample" 2
    check_route_status "opendatahub" "modelregistry-sample-http"
}

# Main function for orchestrating deployments
main() {   
    deploy_and_wait "" $OPENDATAHUB_SUBSCRIPTION 0
    monitoring_crds_installation $OPENDATAHUB_CRDS 120
    deploy_and_wait "" $DSC_INITIALIZATION_MANIFEST 20 
    deploy_and_wait "" $DATA_SCIENCE_CLUSTER_MANIFEST 10 
    monitoring_crds_installation $DATA_SCIENCE_CLUSTER_CRDS 120
    check_pod_status "opendatahub" "-l component.opendatahub.io/name=model-registry-operator" 2 
    deploy_and_wait "-n opendatahub" $MODEL_REGISTRY_SAMPLE_MANIFEST 20
    monitoring_crds_installation $MODEL_REGISTRY_CRDS 120
    deploy_and_wait "-n opendatahub" $MODEL_REGISTRY_DB_MANIFEST 20
    run_deployment_tests
    run_api_tests "-n opendatahub"
}

# Execute main function
main