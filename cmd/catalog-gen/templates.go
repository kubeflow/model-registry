package main

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

// Template path constants
const (
	// cmd templates
	TmplCmdMain = "templates/cmd/main.gotmpl"

	// models templates
	TmplModelsEntity   = "templates/models/entity.gotmpl"
	TmplModelsArtifact = "templates/models/artifact.gotmpl"
	TmplModelsBase     = "templates/models/base.gotmpl"

	// service templates
	TmplServiceRepository         = "templates/service/repository.gotmpl"
	TmplServiceArtifactRepository = "templates/service/artifact_repository.gotmpl"
	TmplServiceSpec               = "templates/service/spec.gotmpl"

	// server templates
	TmplServerOpenAPIServiceImpl = "templates/server/openapi_service_impl.gotmpl"

	// catalog templates
	TmplCatalogLoader = "templates/catalog/loader.gotmpl"

	// providers templates
	TmplProvidersYAML = "templates/providers/yaml.gotmpl"
	TmplProvidersHTTP = "templates/providers/http.gotmpl"

	// api templates
	TmplAPIOpenAPIMain       = "templates/api/openapi_main.gotmpl"
	TmplAPIOpenAPIComponents = "templates/api/openapi_components.gotmpl"

	// manifests templates
	TmplManifestsDeployment       = "templates/manifests/deployment.gotmpl"
	TmplManifestsService          = "templates/manifests/service.gotmpl"
	TmplManifestsSources          = "templates/manifests/sources.gotmpl"
	TmplManifestsSampleCatalog    = "templates/manifests/sample_catalog.gotmpl"
	TmplManifestsKustomization    = "templates/manifests/kustomization.gotmpl"
	TmplManifestsDevSources       = "templates/manifests/dev_sources.gotmpl"
	TmplManifestsDevSampleCatalog = "templates/manifests/dev_sample_catalog.gotmpl"
	TmplManifestsDevKustomization = "templates/manifests/dev_kustomization.gotmpl"

	// misc templates
	TmplMiscMakefile               = "templates/misc/makefile.gotmpl"
	TmplMiscReadme                 = "templates/misc/readme.gotmpl"
	TmplMiscGitignore              = "templates/misc/gitignore.gotmpl"
	TmplMiscOpenAPIGeneratorIgnore = "templates/misc/openapi_generator_ignore.gotmpl"

	// agent templates
	TmplAgentSeedDataSkill   = "templates/agent/seed_data_skill.gotmpl"
	TmplAgentSeedDataCmd     = "templates/agent/seed_data_cmd.gotmpl"
	TmplAgentRegenerateSkill = "templates/agent/regenerate_skill.gotmpl"
)

// executeTemplate reads a template from the embedded filesystem and executes it to a file.
func executeTemplate(templatePath, outputPath string, data any) error {
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New(templatePath).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer func() { _ = file.Close() }()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return nil
}

// executeTemplateToString executes a template and returns the result as a string.
func executeTemplateToString(templatePath string, data any) (string, error) {
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New(templatePath).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}
