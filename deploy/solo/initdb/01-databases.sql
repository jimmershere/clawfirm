-- ClawFirm Solo - per-service database bootstrapping.
-- Each engine gets its own database in the shared Postgres instance.

CREATE DATABASE authentik;
CREATE DATABASE dify;
CREATE DATABASE dify_vec;
CREATE DATABASE n8n;
CREATE DATABASE langgraph;
CREATE DATABASE langfuse;
CREATE DATABASE clawrails;
CREATE DATABASE memory;

\c memory
CREATE EXTENSION IF NOT EXISTS vector;

\c dify_vec
CREATE EXTENSION IF NOT EXISTS vector;
