package constants

// ProjectEnvironment is the string value of the project environment.
type ProjectEnvironment string

const (
	// ProjectEnvStaging is the project environment var for staging.
	ProjectEnvStaging = ProjectEnvironment("staging")
	// ProjectEnvDevelopment is the project environment var for development.
	ProjectEnvDevelopment = ProjectEnvironment("development")
	// ProjectEnvProduction is the project environment var for production.
	ProjectEnvProduction = ProjectEnvironment("production")
	// ProjectEnvCI is the project environment var for CI.
	ProjectEnvCI = ProjectEnvironment("ci")
)
