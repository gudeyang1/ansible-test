#!/bin/bash
rm -fr /var/log/service-deploy-controller/*

if [ -z "${SERVICE_PROPERTY_FILE}" ]; then
    ./service-deploy-controller -installerURL=${INSTALLER_URL} -log_dir=/var/log/service-deploy-controller/
else
    if [ -e "${SERVICE_PROPERTY_FILE}" ]; then
        echo "The property file is ${SERVICE_PROPERTY_FILE}"
    fi
    ./service-deploy-controller -installerURL=${INSTALLER_URL} -propertyFile=${SERVICE_PROPERTY_FILE} -log_dir=/var/log/service-deploy-controller/
fi
