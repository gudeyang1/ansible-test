FROM shc-harbor-dev.hpeswlab.net/itsma/itom-itsma-opensuse-base:dev

RUN mkdir -p /pv/itsma/services/yaml_templates
COPY manifest.yaml /pv/itsma/services/
COPY yaml_templates/*.yaml /pv/itsma/services/yaml_templates/

Add service-deploy-controller start.sh service_property.json /opt/

RUN chmod -R +x /opt/service-deploy-controller /opt/start.sh

WORKDIR /opt
CMD ["./start.sh"]
