ALTER TABLE "ArtifactProperty" ADD CONSTRAINT "ArtifactProperty_artifact_id_fkey" FOREIGN KEY (artifact_id) REFERENCES "Artifact" (id) ON DELETE CASCADE;
ALTER TABLE "Association" ADD CONSTRAINT "Association_context_id_fkey" FOREIGN KEY (context_id) REFERENCES "Context" (id) ON DELETE CASCADE;
ALTER TABLE "Association" ADD CONSTRAINT "Association_execution_id_fkey" FOREIGN KEY (execution_id) REFERENCES "Execution" (id) ON DELETE CASCADE;
ALTER TABLE "Attribution" ADD CONSTRAINT "Attribution_context_id_fkey" FOREIGN KEY (context_id) REFERENCES "Context" (id) ON DELETE CASCADE;
ALTER TABLE "Attribution" ADD CONSTRAINT "Attribution_artifact_id_fkey" FOREIGN KEY (artifact_id) REFERENCES "Artifact" (id) ON DELETE CASCADE;
ALTER TABLE "ContextProperty" ADD CONSTRAINT "ContextProperty_context_id_fkey" FOREIGN KEY (context_id) REFERENCES "Context" (id) ON DELETE CASCADE;
ALTER TABLE "ExecutionProperty" ADD CONSTRAINT "ExecutionProperty_execution_id_fkey" FOREIGN KEY (execution_id) REFERENCES "Execution" (id) ON DELETE CASCADE;
ALTER TABLE "TypeProperty" ADD CONSTRAINT "TypeProperty_type_id_fkey" FOREIGN KEY (type_id) REFERENCES "Type" (id) ON DELETE CASCADE;
