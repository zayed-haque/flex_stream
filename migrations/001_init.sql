CREATE TABLE processed_data (
    id SERIAL PRIMARY KEY,
    original_id TEXT,
    data_type TEXT,
    original_data JSONB,
    processed_result JSONB,
    timestamp TIMESTAMPTZ
);

CREATE INDEX idx_processed_data_original_id ON processed_data(original_id);
CREATE INDEX idx_processed_data_data_type ON processed_data(data_type);
CREATE INDEX idx_processed_data_timestamp ON processed_data(timestamp);