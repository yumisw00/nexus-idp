CREATE EXTENSION IF NOT EXISTS"uuid-ossp";
CREATE TABLE organizations(id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),name TEXT NOT NULL,subscription_tier TEXT,api_key TEXT,is_active BOOLEAN,created_at TIMESTAMPTZ DEFAULT NOW());
CREATE TABLE users(id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),org_id UUID REFERENCES organizations(id),email TEXT UNIQUE NOT NULL,password_hash TEXT NOT NULL,role TEXT,created_at TIMESTAMPTZ DEFAULT NOW());
CREATE TABLE org_wallets(org_id UUID PRIMARY KEY REFERENCES organizations(id),balance_credits NUMERIC,total_spent NUMERIC,updated_at TIMESTAMPTZ DEFAULT NOW());
CREATE TABLE documents(id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),org_id UUID REFERENCES organizations(id),uploader_id UUID REFERENCES users(id),file_name TEXT,file_size_bytes BIGINT,mime_type TEXT,storage_url TEXT,created_at TIMESTAMPTZ DEFAULT NOW());
CREATE TABLE analysis_jobs(id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),document_id UUID REFERENCES documents(id),org_id UUID REFERENCES organizations(id),job_type TEXT,status TEXT,prompt_config JSONB,raw_ai_response JSONB,structured_result JSONB,error_log TEXT,started_at TIMESTAMPTZ,completed_at TIMESTAMPTZ);
CREATE TABLE token_usage_logs(id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),job_id UUID REFERENCES analysis_jobs(id),org_id UUID REFERENCES organizations(id),model_used TEXT,prompt_tokens INT,completion_tokens INT,cost_credits_deducted NUMERIC,created_at TIMESTAMPTZ DEFAULT NOW());
CREATE INDEX idx_users_org ON users(org_id);
CREATE INDEX idx_analysis_jobs_org_status ON analysis_jobs(org_id,status);
