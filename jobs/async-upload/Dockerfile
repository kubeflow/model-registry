###############################################################################
#  Asynchronous Model-Sync Job image for Kubeflow Model Registry (https://github.com/kubeflow/model-registry/issues/1108)
###############################################################################

FROM registry.access.redhat.com/ubi9/python-312-minimal AS base

USER 0

# Install skopeo for Push/Pull of container images to integrate with Model Registry py client and Olot.
RUN microdnf update -y && \
    microdnf install -y skopeo && \
    microdnf clean all

# security/env hardening
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=off \
    PIP_DISABLE_PIP_VERSION_CHECK=on

WORKDIR /app
COPY ./poetry.lock ./pyproject.toml ./requirements.txt ./

RUN python -m pip install -r requirements.txt \
    && echo "Installation completed" \
    && echo "Python path:" && python -c "import sys; print('\\n'.join(sys.path))" \
    && echo "Installed packages:" \
    && python -m pip list 

# Copy application source
COPY . .
RUN ls -l /app

# Switch to USER 1000
RUN chown -R 1000:1000 /app
USER 1000
ENV PATH="/home/1000/.local/bin:$PATH"
RUN ls -l /app

# Sanity checks model_registry module import
RUN echo "Testing model_registry import..." \
&& python -c "import model_registry; print('model_registry package found')"

# TODO: unavailble in OCI format: Add an explicit health-check (K8s will surface this in Pod.status)
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD python -m job.healthcheck

# OCI-recommended labels
ARG VCS_REF=unknown
ARG BUILD_DATE
LABEL org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.title="Kubeflow Model Registry Async-Upload Job" \
    org.opencontainers.image.description="K8s Job image that copies a model between storage back-ends and registers it in Kubeflow Model Registry." \
    org.opencontainers.image.url="https://github.com/kubeflow/model-registry" \
    org.opencontainers.image.source="https://github.com/kubeflow/model-registry/tree/main/jobs/async-upload" \
    org.opencontainers.image.version="0.1.0" \
    org.opencontainers.image.licenses="Apache-2.0" \
    org.opencontainers.image.authors="Kubeflow Model Registry maintainers"

# The Job controller will pass CLI flags/ENV-vars described in issue #1108.
ENTRYPOINT ["python", "-m", "job.entrypoint"]
