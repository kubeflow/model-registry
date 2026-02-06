package main

// generateGitignore generates the .gitignore file.
func generateGitignore() error {
	return executeTemplate(TmplMiscGitignore, ".gitignore", nil)
}

// generateOpenAPIGeneratorIgnore generates the .openapi-generator-ignore file.
func generateOpenAPIGeneratorIgnore() error {
	return executeTemplate(TmplMiscOpenAPIGeneratorIgnore, ".openapi-generator-ignore", nil)
}
