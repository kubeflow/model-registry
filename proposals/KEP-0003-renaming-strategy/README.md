# KEP-0003: Technical Implementation Strategy for Kubeflow Hub Rename

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
  - [Notes/Constraints/Caveats](#notesconstraintscaveats)
  - [Risks and Mitigations](#risks-and-mitigations)
- [Design Details](#design-details)
  - [Core Architectural Decision](#core-architectural-decision)
  - [Impact Analysis by Component](#impact-analysis-by-component)
  - [Test Plan](#test-plan)
  - [Graduation Criteria](#graduation-criteria)
- [Implementation History](#implementation-history)
- [Drawbacks](#drawbacks)
- [Alternatives](#alternatives)
- [References](#references)
<!-- /toc -->

## Summary

This KEP defines the technical implementation strategy for renaming "Kubeflow Model Registry" to "Kubeflow Hub" following community approval of [KEP-907](https://github.com/kubeflow/community/pull/907). The proposal establishes **Kubeflow Hub as an umbrella project name** that groups AI asset management components (Model Registry, Catalog, and future capabilities) while maintaining stable, accurate component names to minimize disruption.

**Key Changes:**
- Repository name: `kubeflow/model-registry` → `kubeflow/hub`
- Go module paths: Updated with major version bump (users encouraged to migrate to Kubeflow SDK)
- Container images: Hard cutover to `ghcr.io/kubeflow/hub/*` namespace in new version
- Documentation: Updated to explain Hub architecture
- API: Add `/catalog/` as simplified alias alongside `/model_catalog/`

**Zero Breaking Changes to:**
- API paths: `/api/model_registry/v1alpha3/` and `/model_catalog/` remain unchanged
- Python SDK: Package name stays `model-registry`
- Kubernetes resources: Service names like `model-registry-ui`, `catalog-service` unchanged
- Database names: No schema changes

**Important Migration Note:**
New version releases will publish container images exclusively to `ghcr.io/kubeflow/hub/*` namespace. Users must update their Kubernetes manifests before upgrading to the new version.

---

## Motivation

The Kubeflow Model Registry project has evolved beyond its original scope to support two distinct use cases:

1. **Model Registry (Tenant-Scoped)**: Tracks model evolution during development, focusing on per-team model iterations, experiments, training runs, parameters, and metrics.

2. **Catalog (Cluster-Scoped)**: Showcases organization-approved models, enables enterprise-wide model sharing, and supports GenAI/LLM model discovery. Also provides MCP (Model Context Protocol) server functionality.

The current name "Kubeflow Model Registry" inadequately represents this broader AI/ML asset management capability. The name "Kubeflow Hub" aligns with industry conventions (Docker Hub, Hugging Face Hub) and provides clearer project identity for future expansion to datasets, prompts, notebooks, and other AI assets.

### Goals

- Rename repository and related artifacts to reflect "Kubeflow Hub" identity
- Minimize breaking changes to existing users and integrations
- Provide clear migration paths with appropriate deprecation periods
- Maintain backwards compatibility where feasible
- Document architectural rationale for component naming decisions
- Coordinate with downstream projects (KServe, Pipelines, Manifests)
- Establish phased migration plan with community communication

### Non-Goals

- Renaming API paths (remain unchanged for backwards compatibility)
- Renaming Python SDK package (component name remains accurate)
- Renaming Kubernetes service names (component names remain accurate)
- Changing database schemas or table names (internal implementation details)
- Immediate removal of legacy artifacts (require deprecation periods)
- Creating new breaking changes beyond what's necessary for repository rename

---

## Proposal

### User Stories

#### Story 1: Existing Go Developer
As a developer currently using the Model Registry Go SDK, I want clear guidance on migrating to the Kubeflow SDK so that I can adopt the officially supported client library.

**Acceptance Criteria:**
- Migration documentation clearly explains transition to Kubeflow SDK
- Kubeflow SDK supports Model Registry functionality
- Migration examples provided for common use cases
- Deprecation timeline for direct Go SDK usage communicated

#### Story 2: Platform Engineer Running Kubernetes
As a platform engineer running Kubeflow Hub in production, I want service names and API paths to remain unchanged so that my existing deployments continue working without modification.

**Acceptance Criteria:**
- Kubernetes service DNS names remain stable
- API endpoints remain unchanged
- Existing manifests continue to work
- Optional labels added for Hub context without breaking selectors

#### Story 3: Data Scientist Using Python SDK
As a data scientist using the Python SDK, I want my existing notebooks and scripts to continue working without changes so that I can focus on model development.

**Acceptance Criteria:**
- `pip install model-registry` continues to work
- Import statements remain unchanged
- Existing tutorials and notebooks work without modification
- Clear documentation explains relationship to Kubeflow Hub

#### Story 4: CI/CD Pipeline Maintainer
As a CI/CD pipeline maintainer, I want clear advance notice of the container registry namespace change so that I can update my configurations before the new version releases.

**Acceptance Criteria:**
- Container registry namespace change announced well in advance
- New version uses `ghcr.io/kubeflow/hub/*` namespace exclusively
- Migration timeline clearly communicated in release notes
- Documentation updated with new image references

### Notes/Constraints/Caveats

1. **Repository Redirects**: GitHub automatically redirects `kubeflow/model-registry` to `kubeflow/hub`, but this doesn't update Go module paths in code.

2. **Go Module Versioning**: Go module path changes require a major version bump (e.g., v0.3.0 or v1.0.0) to follow semantic versioning properly.

3. **Container Registry Hard Cutover**: Cannot create automatic redirects for container images. New version releases will publish exclusively to `ghcr.io/kubeflow/hub/*` namespace, requiring users to update manifests before upgrading.

4. **Downstream Coordination**: Changes require coordination with:
   - kubeflow/manifests (deployment configurations)
   - kubeflow/pipelines (integration components)
   - KServe (storage initializer references)

5. **SDK Clarification**: The existing `model-registry` Python package is specific to Model Registry component; a separate `kubeflow/sdk` project provides broader Kubeflow Hub client under `kubeflow.hub` module.

### Risks and Mitigations

| Risk | Impact | Likelihood | Mitigation |
|:-----|:-------|:-----------|:-----------|
| Users upgrade without updating manifests | Critical | Medium | Clear advance communication; migration guide in release notes; coordinated release with kubeflow/manifests |
| Go import path confusion | Medium | Medium | Direct users to Kubeflow SDK; provide migration guide; GitHub redirects help |
| Confusion about component vs. project naming | Medium | Medium | Clear documentation explaining Hub umbrella vs. component names; FAQ document |
| External documentation links break | Low | High | GitHub automatically redirects URLs; monitor broken links; update major references |
| Air-gapped environments can't pull new images | High | Medium | Well-advertised timeline; image migration checklist; coordinate with enterprise users |
| KServe/Pipelines integrations break | High | Low | Coordinate releases; test integration points; maintain API path stability |

---

## Design Details

### Core Architectural Decision

**Kubeflow Hub is an umbrella project name** that groups AI asset management components:

- **Model Registry** (tenant-scoped model evolution tracking)
- **Catalog** (cluster-scoped approved model showcase and MCP servers)
- **Future capabilities** (datasets, prompts, notebooks, etc.)

**Critical Principle**: Component names (`model-registry`, `catalog`) accurately describe their specific functionality and **remain unchanged** even under the Hub umbrella. This minimizes breaking changes while allowing the project to evolve.

### Impact Analysis by Component

### 1\. Repository Name Change

**Current**: `github.com/kubeflow/model-registry` **Proposed**: `github.com/kubeflow/hub`

#### Pros

- Clean break with new identity  
- GitHub automatically redirects old URLs to new repository  
- Fresh namespace for issues, PRs, discussions

#### Cons

- **HIGH IMPACT**: All Go import paths break immediately  
- All existing forks reference old repository  
- All external documentation links need updating  
- Bookmarks and CI/CD configurations across the ecosystem break

#### Backwards Compatibility Impact

| Aspect | Impact Level | Mitigation |
| :---- | :---- | :---- |
| Go imports | **BREAKING** | Requires code changes in all consumers |
| Git clone URLs | Low | GitHub provides automatic redirects |
| Issue/PR links | Low | GitHub provides automatic redirects |
| Release artifacts | Medium | Old tags remain accessible via redirects |

#### Recommendation

**PHASE 1** \- Rename repository with redirect.

---

### 2\. Go Module Path Changes

**Current**: `github.com/kubeflow/model-registry` **Proposed**: `github.com/kubeflow/hub`

#### Affected Files (Primary)

- `/go.mod` (root module)  
- `/pkg/openapi/go.mod`  
- `/catalog/pkg/openapi/go.mod`  
- `/clients/ui/bff/go.mod`  
- `/gorm-gen/go.mod`

#### Affected Imports

**889+ occurrences across 325+ files** require updating.

#### Pros

- Consistent with new project identity  
- Go module versioning can restart cleanly

#### Cons

- **BREAKING CHANGE**: All downstream Go consumers must update imports  
- Requires coordinated release with consumers (kubeflow/manifests, kubeflow/pipelines integrations)  
- Existing `replace` directives in external projects break

#### Backwards Compatibility Impact

| Consumer Type | Impact | Mitigation Strategy |
| :---- | :---- | :---- |
| Direct Go imports | **BREAKING** | Repository redirects |
| `go get` users | **BREAKING** | Document migration path |
| Vendored dependencies | **BREAKING** | Announce in release notes |

#### Recommendation

**DO**: Rename Go module paths as part of a major version bump (e.g., v0.3.0 or v1.0.0). **CONSIDER**: Using Go module `replace` directives temporarily to support both import paths during transition.

---

### 3\. Container Image Names

**Current Registry**: `ghcr.io/kubeflow/model-registry/` **Images**:

- `ghcr.io/kubeflow/model-registry/server`  
- `ghcr.io/kubeflow/model-registry/ui`  
- `ghcr.io/kubeflow/model-registry/ui-standalone`  
- `ghcr.io/kubeflow/model-registry/ui-federated`  
- `ghcr.io/kubeflow/model-registry/storage-initializer`  
- `ghcr.io/kubeflow/model-registry/async-upload`

**Proposed Registry**: `ghcr.io/kubeflow/hub/`

#### Pros

- Consistent with new branding  
- Clear identity in container registries

#### Cons

- **HIGH IMPACT**: All Kubernetes deployments referencing old images break  
- Helm charts, kustomize overlays, and external deployments need updating  
- CI/CD pipelines across the ecosystem need updating

#### Backwards Compatibility Impact

| Aspect | Impact Level | Mitigation |
| :---- | :---- | :---- |
| Existing deployments | **BREAKING** | Image tags won't resolve |
| Kubeflow Manifests | **BREAKING** | Will need update to new images following a release |
| Helm values | **BREAKING** | Requires values.yaml changes |
| Air-gapped environments | **HIGH** | Must re-mirror images |

#### Recommendation

**DO**: Announce container registry namespace change well in advance of release
**DO**: Update all internal manifests and documentation to use new namespace
**DO**: Coordinate with kubeflow/manifests team for synchronized update
**DO**: Publish new version images exclusively to `ghcr.io/kubeflow/hub/*` namespace
**DO**: Provide clear migration guidance in release notes

**Note**: This is a hard cutover - new releases will only publish to the new registry namespace. Users must update their manifests to reference `ghcr.io/kubeflow/hub/*` before upgrading to the new version.

#### Implementation Strategy

```
# CI workflow change - publish to new registry only
- name: Push to new registry (kubeflow/hub)
  run: docker push ghcr.io/kubeflow/hub/server:${{ env.VERSION }}
```

#### Consumer Migration Path

Users must update their manifests before upgrading to the new version:

```
# Before (old versions use legacy registry)
image: ghcr.io/kubeflow/model-registry/server:v0.2.x

# After (new versions use new registry exclusively)
image: ghcr.io/kubeflow/hub/server:v0.3.0
```

**Important**: Old image tags remain available at the legacy registry path for previous versions, but new releases will only be published to the new namespace.

---

### 4\. API Path Changes

**Current Paths**:

- Model Registry: `/api/model_registry/v1alpha3/`
- Catalog: `/model_catalog/`

**Architectural Decision**: Kubeflow Hub is a *conceptual grouping* of AI asset management components. Internal component names ("model-registry", "catalog") remain unchanged as they accurately describe their specific functionality.

#### Model Registry API \- No Changes

**Decision**: **KEEP** `/api/model_registry/v1alpha3/` unchanged.

**Rationale**:

- Zero breaking changes for API consumers
- Model Registry accurately describes this component's purpose (tracking model evolution)
- Component name remains valid even under Hub umbrella
- Python clients, integrations continue working seamlessly

#### Catalog API \- Additive Changes Only

**Current**: `/model_catalog/` **Proposed Addition**: `/catalog/` (routes to same handler)

**Decision**: **ADD** `/catalog/` as an alias path alongside existing `/model_catalog/` path.

**Rationale**:

- **Simplified naming**: Cleaner path that aligns with component name
- **Backward compatibility**: Existing `/model_catalog/` continues working
- **Zero breaking changes**: Additive only, no removals
- **Future-proof**: Generic "catalog" name supports expansion to MCP servers and other AI assets

**Implementation**:

```go
// Both paths route to the same handler
router.Handle("/model_catalog/", catalogHandler) // Legacy path
router.Handle("/catalog/", catalogHandler)       // New simplified alias
```

#### Future Considerations

- Catalog already supports MCP servers and can be extended to other asset types
- Component-specific paths (`/model_catalog/`, `/model_registry/`) remain authoritative for backwards compatibility

#### Recommendation

**DO**: Add `/catalog/` path as alias for catalog functionality
**DO**: Keep all existing API paths unchanged
**DO NOT**: Rename existing API paths
**DOCUMENT**: Clearly explain that both `/model_catalog/` and `/catalog/` paths are supported and equivalent

---

### 5\. Python Client Package

**Current**:

- PyPI package: `model-registry`  
- Import path: `from model_registry import ModelRegistry`  
- OpenAPI module: `mr_openapi`

**Architectural Decision**: The Python SDK provides programmatic access to the **Model Registry component** specifically. Since "Model Registry" accurately describes this component's functionality (tracking model evolution, versioning, metadata), the package name remains correct even under the Kubeflow Hub umbrella.

#### Decision: Keep Package Name Unchanged

**Rationale**:

- **Component naming accuracy**: "model-registry" correctly describes what this SDK does \- it's the client for the Model Registry component  
- **Zero breaking changes**: All existing code continues working  
- **Kubeflow Hub is a grouping concept**: The Hub groups Model Registry \+ Catalog \+ future AI asset components  
- **SDK scope**: This SDK is specifically for Model Registry operations, not for the entire Hub  
- **Industry precedent**: Docker Hub contains many SDKs that keep component-specific names

#### Backwards Compatibility Impact

| Aspect | Impact Level | Mitigation |
| :---- | :---- | :---- |
| pip install | **NONE** | No changes required |
| Import statements | **NONE** | No changes required |
| Existing notebooks | **NONE** | No changes required |
| Tutorials | **NONE** | No changes required |

#### Future Considerations

- Component-specific SDKs remain more useful than monolithic Hub SDK for most use cases

#### Recommendation

**DO NOT** rename the Python package. Keep `model-registry` on PyPI. **DOCUMENT**: Clearly explain that `model-registry` is the SDK for the Model Registry component within Kubeflow Hub. **DOCUMENT**: Any new “user facing” client shall be developed directly in Kubeflow/SDK.\[**CONSIDER**: Add Kubeflow Hub branding to documentation while keeping functional package name.

---

### 6\. Kubernetes Manifests and CRDs

**Current References**:

- Service names: `model-registry-ui`, `model-registry`, `catalog-service`  
- Deployment names: `model-registry-deployment`, `catalog-deployment`  
- ConfigMaps: `model-registry-*`, `catalog-*`  
- Labels: `app: model-registry-ui`, `app: catalog`  
- RBAC: `model-registry-manager-role`, `catalog-manager-role`

#### Critical Analysis: Should Kubernetes Resources Be Renamed?

**Architectural Question**: If Kubeflow Hub is a *conceptual grouping* of AI asset management components (Model Registry \+ Catalog), and these component names remain valid and descriptive, do the Kubernetes resources need renaming?

**Answer: NO \- Component Names Should Remain Unchanged**

#### Rationale for Keeping Current K8s Resource Names

1. **Kubeflow Hub is organizational, not operational**  
     
   - "Hub" is the project/repository name and umbrella concept  
   - Actual deployed components are "Model Registry" and "Catalog"  
   - K8s resources should reflect deployed component functionality, not project branding

   

2. **Component names are accurate and specific**  
     
   - `model-registry-ui` accurately describes what this service does  
   - `catalog-service` accurately describes what this service does  
   - These names remain correct under the Hub umbrella

   

3. **Industry precedent**  
     
   - Docker Hub doesn't name its services `hub-*`  
   - Artifact repositories don't rename component services to match umbrella branding  
   - Component-level naming is standard practice

   

4. **Zero breaking changes**  
     
   - Service DNS names remain stable (`model-registry.kubeflow.svc.cluster.local`)  
   - Label selectors don't break  
   - RBAC bindings continue working  
   - External integrations (KServe, Pipelines) don't break

   

5. **Separation of concerns**  
     
   - Repository name (`kubeflow/hub`) \= where code lives  
   - Component names (`model-registry`, `catalog`) \= what gets deployed  
   - These don't need to match

#### Backwards Compatibility Impact

| Aspect | Impact Level | Mitigation |
| :---- | :---- | :---- |
| Service DNS names | **NONE** | No changes |
| Label selectors | **NONE** | No changes |
| RBAC bindings | **NONE** | No changes |
| kubeflow/manifests sync | **LOW** | Update comments/documentation only |

#### Recommendation

**DO NOT** rename Kubernetes service names, deployments, or CRDs
**DO** update documentation to explain component naming under Hub umbrella
**DO NOT** break existing service DNS names or label selectors

#### Files Requiring Updates

- Documentation comments in manifest files  
- README files explaining the deployment  
- **NO** functional resource name changes needed

---

### 7\. Documentation Updates

#### Documentation Strategy: Emphasize Architectural Clarity

**Key Message**: Kubeflow Hub is an **umbrella project** grouping AI asset management components (Model Registry, Catalog, and future capabilities). Component names remain unchanged as they accurately describe specific functionality.

#### Internal Documentation (This Repository)

| File/Path | Update Type | Key Changes |
| :---- | :---- | :---- |
| `README.md` | **MAJOR** | \- Introduce "Kubeflow Hub" as project name \- Explain Hub as grouping of Model Registry \+ Catalog \- Clarify component names remain unchanged \- Update repository URLs |
| `CONTRIBUTING.md` | **MINOR** | \- Update repository URL references \- Keep component-specific contribution guides |
| `RELEASE.md` | **MODERATE** | \- Update container registry publishing naming \- Note image naming in both registries \- Keep component-specific release notes |
| `docs/*.md` | **MODERATE** | \- Add "Kubeflow Hub" context to introductions \- Maintain component-specific documentation \- Clarify Model Registry vs. Catalog distinction |
| `clients/python/README.md` | **MINOR** | \- Add note: "Python SDK for Model Registry component within Kubeflow Hub" \- Keep all functional documentation unchanged |
| `CLAUDE.md` | **MINOR** | \- Update repository references \- Note: component names unchanged |
| `manifests/*/README.md` | **MODERATE** | \- Explain K8s resource naming strategy \- Clarify why service names remain unchanged |

#### External Documentation

| Location | Owner | Action Required |
| :---- | :---- | :---- |
| kubeflow.org | kubeflow/website | update icon by raising CNCF ticket (Matteo can do it) update landing page in the root website page |
| [kubeflow.org/docs](https://www.kubeflow.org/docs/components/model-registry/) | kubeflow/website | \- Update page title to "Kubeflow Hub" \- Add redirect from \`/model-registry/\` to \`/hub/\` or \`/kubeflow-hub/\` \- Document Model Registry and Catalog as Hub components \- Maintain separate guides for each component |
| Community Blog Posts | Community | \- Publish announcement post explaining rename \- Clarify what's changing vs. unchanged \- Historical posts remain as-is (searchable) |
| YouTube Tutorials | Community | \- Add pinned comment for new videos explaining rename for a sensible period \- Link to migration guide \- Add entry line in biweekly meeting notes \- Videos remain valid (component functionality unchanged) |

#### New Documentation Required

1. **Migration Guide** (`docs/migration-to-kubeflow-hub.md`)  
     
   - Repository URL changes  
   - What's NOT changing (critical section)  
   - Timeline and support calendar

   

2. **Architecture Overview** (`docs/architecture/kubeflow-hub-components.md`)  
     
   - Explain Hub as umbrella concept  
   - Detail Model Registry component  
   - Detail Catalog component  
   - Explain why component names remain unchanged  
   - Future expansion areas (datasets, prompts, etc.)

   

3. **FAQ Document** (`docs/faq-kubeflow-hub-rename.md`)  
     
   - "Why keep 'model-registry' in API paths?"  
   - "Why keep 'model-registry' Python package?"  
   - "Why keep K8s service names unchanged?"  
   - "What is Kubeflow Hub vs. Model Registry?"

#### Recommendation

**DO**: Create comprehensive "What's Changing vs. What's Not" documentation **DO**: Emphasize that functional component names remain valid and unchanged **DO**: Coordinate with kubeflow/website for same-day documentation update **DO**: Update component/project icon and update [kubeflow.org](http://kubeflow.org) landing page accordingly **DO**: Add "Kubeflow Hub" branding while maintaining component-specific docs **DO**: Create migration guide even though breaking changes are minimized **DOCUMENT**: Clearly explain architectural decision to keep component names

---

### 8\. CI/CD Pipeline Changes

#### GitHub Actions Workflows Affected

| Workflow | Changes Required |
| :---- | :---- |
| `build-and-push-image.yml` | Update IMG\_REPO, image tags |
| `build-and-push-ui-images.yml` | Update image names |
| `build-and-push-csi-image.yml` | Update image names |
| `trivy-image-scanning.yaml` | Update image references |
| All workflows | Repository name in URLs |

#### Container Registry Configuration

- GitHub Packages namespace change  
- GHCR permissions and tokens  
- Signing keys for images

#### Recommendation

**DO**: Update all workflows in the rename PR. **DO**: Test thoroughly in a fork before merge. **DO**: Maintain ability to push to both old and new registries during transition.

---

### 9\. Database and Schema Considerations

**Current Database Names**:

- `model_registry` (MySQL)  
- `model_catalog` (PostgreSQL)

**Table Names**: Internal to database, not exposed externally.

#### Recommendation

**DO NOT** rename database or table names. This provides:

- Zero impact on existing deployments  
- No migration scripts needed for data  
- Seamless upgrade path

Database names are internal implementation details and need not match the project name.

---

### 10\. Configuration Files and Environment Variables

**Current State**:

- Config file: `.model-registry.yaml` (currently unused in codebase)
- Environment prefix: `MODEL_REGISTRY_*`

#### Analysis

**.model-registry.yaml Config File**: Investigation shows this file is not actually used by the application. It can be removed without impact.

**Environment Variables**: The `MODEL_REGISTRY_*` prefix accurately describes the Model Registry component and remains valid under the Hub umbrella. No changes needed.

#### Backwards Compatibility Impact

| Item | Impact | Recommendation |
| :---- | :---- | :---- |
| Config file name | **NONE** | Remove unused `.model-registry.yaml` references |
| Env variables | **NONE** | Keep `MODEL_REGISTRY_*` prefix unchanged |

#### Recommendation

**DO**: Remove unused `.model-registry.yaml` file and references if present
**DO NOT**: Create new `.kubeflow-hub.yaml` config file (not needed)
**DO NOT**: Rename environment variables - `MODEL_REGISTRY_*` remains accurate for the Model Registry component
**DOCUMENT**: Note that config file was unused and removed for code cleanup

---

## Additional Areas Not Previously Covered

### 11\. Helm Charts (If Applicable)

**Current State**: Project uses Kustomize, not Helm.

**Future Consideration**: If Helm charts are created, use new naming from start.

### 12\. Metrics and Observability

**Current Prometheus Metrics**: May include `model_registry_*` prefixes.

**Recommendation**: Audit metrics endpoints for naming. Consider:

- Keeping old metric names for dashboard compatibility  
- Adding new metric names with old as aliases  
- Documenting metric name changes for monitoring teams

### 13\. Logging and Tracing

**Log Prefixes**: `[model-registry]` or similar.

**Recommendation**: Update log prefixes but note this may affect log parsing rules in production environments.

### 14\. Integration Points

#### KServe Integration

- InferenceService annotations  
- Custom Storage Initializer naming  
- Controller CRD references

#### Kubeflow Pipelines Integration

- Component references  
- SDK imports

**Recommendation**: Coordinate with KServe and Pipelines teams for synchronized updates.

### 15\. Security Considerations

#### Image Signing

- New images need new signatures  
- Cosign/Sigstore configuration updates

#### SBOM (Software Bill of Materials)

- Security scanning references

#### Vulnerability Databases

- CVE references may use old project name

**Recommendation**: Document security artifact migration in release notes.

---

### Test Plan

Testing for this KEP focuses on validating migration paths and backwards compatibility rather than new functional code.

**Configuration Cleanup Testing:**
- Verify removal of unused `.model-registry.yaml` file doesn't impact functionality
- Verify `MODEL_REGISTRY_*` environment variables continue working correctly

**API Path Compatibility Testing:**
- Verify both `/model_catalog/` and `/catalog/` paths route correctly
- Verify identical responses from both paths
- Verify existing API clients continue working

**Container Image Migration Testing:**
- Deploy new version using `ghcr.io/kubeflow/hub/*` images
- Verify deployment succeeds with new namespace
- Verify old version images remain accessible at `ghcr.io/kubeflow/model-registry/*` for rollback scenarios

**Integration Testing:**
- Test KServe storage initializer with catalog
- Test Kubeflow Pipelines model registration
- Test Python SDK operations
- Test Go SDK with updated import paths

### Graduation Criteria

This KEP follows a phased migration approach with clear gates:

**Alpha (Weeks 1-4: Preparation)**

*Entry Criteria:*
- Community approval of KEP-0003
- Migration documentation drafted
- CI/CD updated for new registry namespace

*Exit Criteria:*
- Community announcement published
- Downstream projects (KServe, Pipelines, Manifests) notified
- FAQ document created
- Migration timeline agreed
- Coordination with kubeflow/manifests team for synchronized updates

**Beta (Weeks 5-8: Repository Rename and Version Release)**

*Entry Criteria:*
- Alpha phase complete
- All stakeholders ready for change
- Internal manifests updated to new registry namespace
- Rollback plan documented

*Exit Criteria:*
- Repository renamed with GitHub redirects active
- Go module paths updated with version bump
- New version published with container images at `ghcr.io/kubeflow/hub/*`
- Internal documentation updated
- `/catalog/` alias implemented
- kubeflow.org documentation updated
- kubeflow/manifests updated with new image references
- No critical issues reported

**Stable (Post-Release: Ongoing Support)**

*Entry Criteria:*
- Beta phase complete
- New version successfully released
- Migration guide published

*Exit Criteria:*
- Community successfully upgraded to new version
- No critical migration blockers reported
- Community feedback positive
- Migration considered complete

---

## Implementation History

- **2025-01-15**: KEP-907 approved renaming Model Registry to Kubeflow Hub
- **2026-02-02**: KEP-0003 created for technical implementation strategy
- **2026-02-16**: KEP-0003 updated to follow proper KEP format and refined based on implementation decisions
- **TBD**: Community review and approval of KEP-0003
- **TBD**: Alpha phase begins (Preparation and coordination)
- **TBD**: Beta phase begins (Repository rename and version release)
- **TBD**: Stable phase (Post-release support)
- **TBD**: Implementation complete

---

## Drawbacks

While this proposal minimizes breaking changes through careful architectural decisions, several drawbacks remain:

1. **Go Ecosystem Disruption**: Existing Go consumers are encouraged to migrate to the Kubeflow SDK which provides broader functionality. Direct use of internal Go modules may require import path updates from `github.com/kubeflow/model-registry` to `github.com/kubeflow/hub`.

2. **Container Registry Hard Cutover**: Platform engineers must update Kubernetes manifests, Helm values, CI/CD pipelines before upgrading to the new version. New releases will only publish to `ghcr.io/kubeflow/hub/*` namespace. Risk of deployment failures if manifests aren't updated before version upgrade. Air-gapped environments must re-mirror images from new registry.

3. **Version Upgrade Coordination**: The hard cutover for container images means users must coordinate manifest updates with version upgrades. Cannot upgrade to new version without first updating image references.

4. **Naming Conceptual Complexity**: Introducing dual naming (Hub for project umbrella, component names for actual services) may confuse new users despite providing clear architectural benefits and backwards compatibility.

5. **Documentation Fragmentation**: Historical blog posts, videos, and tutorials become partially outdated. External documentation across the ecosystem requires updates, and search results may show mixed naming for extended periods.

6. **Coordination Burden**: Requires synchronization with multiple Kubeflow projects (manifests, KServe, Pipelines), potentially delaying their releases or creating version incompatibilities during transition.

7. **Migration Fatigue**: Community may experience "change fatigue" during the migration period, particularly for users who must update multiple systems and configurations.

---

## Alternatives

### Alternative 1: Keep "Model Registry" Name

**Approach**: Retain current name despite evolved functionality supporting both Registry and Catalog use cases.

**Pros:**
- Zero breaking changes
- No migration effort required
- No ecosystem disruption
- Immediate cost savings

**Cons:**
- Name doesn't accurately represent dual Registry + Catalog functionality
- Confusing for new users discovering catalog-only use cases
- Limits perception of future expansion to datasets, prompts, and other AI assets
- Inconsistent with industry naming patterns (Docker Hub, HuggingFace Hub)
- Missed opportunity for clearer project identity

**Decision**: Rejected - Community already approved rename via KEP-907 based on these limitations.

---

### Alternative 2: Rename Everything Including Component Names

**Approach**: Comprehensive rename where API paths become `/api/hub/v1alpha3/`, Python package becomes `kubeflow-hub`, Kubernetes services become `hub-ui`, `hub-catalog-service`, etc.

**Pros:**
- Complete naming consistency across all layers
- Clear "Kubeflow Hub" brand identity everywhere
- No dual naming complexity

**Cons:**
- **MASSIVE breaking changes** across entire ecosystem:
  - All API consumers must update code
  - All Python notebooks and tutorials break
  - All Kubernetes integrations break (KServe, Pipelines)
  - Service DNS names change, breaking external dependencies
  - Database migrations or permanent naming disconnect
- 12-24 month migration timeline with high fragmentation risk
- Component names like "model-registry" and "catalog" are accurate and descriptive
- Much higher disruption for marginal naming consistency benefit

**Decision**: Rejected - Breaking changes far outweigh benefits; component names remain functionally accurate.

---

### Alternative 3: Create New "Hub" Repository Alongside Existing

**Approach**: Create `kubeflow/hub` as new repository, gradually migrate Model Registry and Catalog components over time while maintaining `kubeflow/model-registry`.

**Pros:**
- No immediate breaking changes
- Users migrate at their own pace
- Legacy repository remains stable during transition
- Can experiment with new structure without disrupting existing users

**Cons:**
- Fragments community across two repositories
- Duplicates issues, PRs, discussions, and releases
- Confusing for new contributors - which repository to use?
- Unclear which repository owns what components
- No clear sunset path for legacy repository
- Double the maintenance burden
- Perpetuates confusion this rename aims to resolve

**Decision**: Rejected - Repository fragmentation creates worse long-term problems than coordinated migration.

---

### Alternative 4: Dual Container Registry Publishing with Gradual Transition

**Approach**: Publish container images to both `ghcr.io/kubeflow/model-registry/*` and `ghcr.io/kubeflow/hub/*` for an extended transition period (6-24 months) before deprecating the legacy registry.

**Pros:**
- More time for ecosystem to migrate
- Reduced pressure on downstream projects
- Lower risk of production disruptions
- Users can upgrade version without updating manifests

**Cons:**
- Extended period maintaining dual CI/CD publishing
- Prolonged community confusion about which registry to use
- Higher maintenance overhead and CI/CD complexity
- Doubles storage costs for container images
- Delays clear migration completion
- Users may not migrate unless forced

**Decision**: Rejected - Hard cutover provides clearer migration signal, simpler CI/CD, and forces necessary updates. Well-advertised timeline and coordination with kubeflow/manifests mitigates disruption risk.

---

### Alternative 5: Metadata-Only Rename (No Repository Changes)

**Approach**: Update documentation and branding to reference "Kubeflow Hub" but leave repository name as `kubeflow/model-registry`.

**Pros:**
- Zero breaking changes
- Immediate branding benefit in documentation
- No migration effort required
- Lowest risk approach

**Cons:**
- Permanent disconnect between project name and repository URL
- Confusing for contributors trying to find code ("Where is kubeflow/hub?")
- Container images remain under misleading `ghcr.io/kubeflow/model-registry` namespace
- Half-measure that doesn't fully realize renaming benefits
- Search results and external links remain misaligned
- Go import paths still say `model-registry`

**Decision**: Rejected - Incomplete solution that creates long-term confusion and doesn't achieve community's goals from KEP-907.

---

## Summary Matrix

| Area | Rename? | Breaking? | Deprecation Period | Priority |
|:-----|:--------|:----------|:-------------------|:---------|
| GitHub Repository | YES | Low (redirects) | N/A | P0 |
| Go Module Paths | YES | HIGH | N/A (use Kubeflow SDK) | P0 |
| Container Images | YES | HIGH | N/A (hard cutover) | P0 |
| API Paths | NO\* | - | - | - |
| Python Package | NO | - | - | - |
| K8s Manifests | NO\*\* | - | - | - |
| Documentation | YES | Low | Immediate | P1 |
| CI/CD | YES | Low | With repo rename | P0 |
| Database Names | NO | - | - | - |
| Config Files | REMOVE\*\*\* | NONE | N/A | P2 |
| Env Variables | NO | - | - | - |

**Notes:**
- \*API Paths: Keep existing, ADD `/catalog/` alias for catalog
- \*\*K8s Manifests: Keep component names unchanged
- \*\*\*Config Files: Remove unused `.model-registry.yaml` file (currently not used by application)

---

## References

- [KEP-907: Model Registry Renaming](https://github.com/kubeflow/community/pull/907) - Community proposal approving the rename
- [KEP-907 Proposal Directory](https://github.com/kubeflow/community/tree/master/proposals/907-model-registry-renaming) - Full proposal documentation
- [Semantic Versioning](https://semver.org/) - Version numbering guidelines
- [Go Module Versioning](https://go.dev/doc/modules/version-numbers) - Go module version best practices
- [Kubernetes Recommended Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/) - Standard label conventions

