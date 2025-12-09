-- Create profile_info table
CREATE TABLE IF NOT EXISTS profile_info (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    photo_url VARCHAR(500),
    activity TEXT NOT NULL,
    contacts JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Insert default profile data
INSERT INTO profile_info (name, description, photo_url, activity, contacts)
VALUES (
    'Кривоносов Константин',
    'Интересуюсь AI-разработкой',
    '/base-avatar.jpg',
    'Воркаю.',
    '{"email": "koskriv2006@gmail.com", "github": "https://github.com/fUS1ONd", "linkedin": "https://linkedin.com/in/alexandr-petrov"}'::jsonb
) ON CONFLICT (id) DO NOTHING;

-- Create index on updated_at for efficient queries
CREATE INDEX idx_profile_info_updated_at ON profile_info(updated_at DESC);

-- Add comment to table
COMMENT ON TABLE profile_info IS 'Stores CV/profile information for the website owner';
