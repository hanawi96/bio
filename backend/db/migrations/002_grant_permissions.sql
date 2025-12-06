-- Grant all privileges on all tables to linkbio user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO linkbio;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO linkbio;

-- Grant privileges on future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO linkbio;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO linkbio;

-- Ensure linkbio user can create tables
GRANT CREATE ON SCHEMA public TO linkbio;
