package main

import (
	"embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

// Template path constants
const (
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

	// misc templates
	TmplMiscGitignore              = "templates/misc/gitignore.gotmpl"
	TmplMiscOpenAPIGeneratorIgnore = "templates/misc/openapi_generator_ignore.gotmpl"

	// plugin templates
	TmplPluginPlugin   = "templates/plugin/plugin.gotmpl"
	TmplPluginRegister = "templates/plugin/register.gotmpl"

	// agent templates
	TmplAgentClaudeMD             = "templates/agent/claude_md.gotmpl"
	TmplAgentCmdAddProperty       = "templates/agent/commands/add_property.gotmpl"
	TmplAgentCmdAddArtifact       = "templates/agent/commands/add_artifact.gotmpl"
	TmplAgentCmdAddArtifactProp   = "templates/agent/commands/add_artifact_property.gotmpl"
	TmplAgentCmdRegenerate        = "templates/agent/commands/regenerate.gotmpl"
	TmplAgentCmdFixBuild          = "templates/agent/commands/fix_build.gotmpl"
	TmplAgentSkillAddProperty     = "templates/agent/skills/add_property.gotmpl"
	TmplAgentSkillAddArtifact     = "templates/agent/skills/add_artifact.gotmpl"
	TmplAgentSkillAddArtifactProp = "templates/agent/skills/add_artifact_property.gotmpl"
	TmplAgentSkillRegenerate      = "templates/agent/skills/regenerate.gotmpl"
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
