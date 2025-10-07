ALTER TABLE "ArtifactProperty" DROP CONSTRAINT "ArtifactProperty_artifact_id_fkey";
ALTER TABLE "Association" DROP CONSTRAINT "Association_context_id_fkey";
ALTER TABLE "Association" DROP CONSTRAINT "Association_execution_id_fkey";
ALTER TABLE "Attribution" DROP CONSTRAINT "Attribution_context_id_fkey";
ALTER TABLE "Attribution" DROP CONSTRAINT "Attribution_artifact_id_fkey";
ALTER TABLE "ContextProperty" DROP CONSTRAINT "ContextProperty_context_id_fkey";
ALTER TABLE "ExecutionProperty" DROP CONSTRAINT "ExecutionProperty_execution_id_fkey";
ALTER TABLE "TypeProperty" DROP CONSTRAINT "TypeProperty_type_id_fkey";
