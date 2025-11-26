-- Seed data for development and testing
-- This script can be run independently after migrations

-- Insert admin user
INSERT INTO users (email, role, created_at)
VALUES ('admin@example.com', 'admin', NOW())
ON CONFLICT (email) DO NOTHING;

-- Insert regular users for testing
INSERT INTO users (email, role, created_at)
VALUES
    ('user1@example.com', 'user', NOW()),
    ('user2@example.com', 'user', NOW()),
    ('user3@example.com', 'user', NOW())
ON CONFLICT (email) DO NOTHING;

-- Insert notification settings for all users
INSERT INTO notification_settings (user_id, email_enabled, push_enabled, new_posts_enabled)
SELECT id, TRUE, TRUE, TRUE FROM users
ON CONFLICT (user_id) DO NOTHING;

-- Insert sample posts
INSERT INTO posts (title, slug, content, preview, author_id, published, published_at, created_at, updated_at)
SELECT
    'Getting Started with AI Development',
    'getting-started-with-ai-development',
    E'# Introduction\n\nArtificial Intelligence is transforming the way we build applications. In this post, we''ll explore the basics of AI development and how to get started with machine learning.\n\n## What is AI?\n\nArtificial Intelligence (AI) refers to the simulation of human intelligence in machines that are programmed to think like humans and mimic their actions.\n\n## Getting Started\n\nTo begin your AI journey, you''ll need:\n\n1. **Python** - The most popular language for AI development\n2. **Libraries** - NumPy, Pandas, TensorFlow, or PyTorch\n3. **Understanding of mathematics** - Linear algebra, calculus, and statistics\n\n## Conclusion\n\nAI development is an exciting field with endless possibilities. Start small, practice regularly, and never stop learning!',
    'Learn the fundamentals of AI development and how to get started with machine learning.',
    (SELECT id FROM users WHERE email = 'admin@example.com'),
    TRUE,
    NOW(),
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM posts WHERE slug = 'getting-started-with-ai-development');

INSERT INTO posts (title, slug, content, preview, author_id, published, published_at, created_at, updated_at)
SELECT
    'Understanding Neural Networks',
    'understanding-neural-networks',
    E'# Neural Networks Explained\n\nNeural networks are the backbone of modern AI systems. Let''s dive deep into how they work.\n\n## What are Neural Networks?\n\nA neural network is a series of algorithms that endeavors to recognize underlying relationships in a set of data through a process that mimics the way the human brain operates.\n\n## Architecture\n\n### Input Layer\nThe input layer receives the initial data for processing.\n\n### Hidden Layers\nHidden layers perform most of the computations required by the network.\n\n### Output Layer\nThe output layer produces the final result.\n\n## Training Process\n\n1. Forward propagation\n2. Calculate loss\n3. Backward propagation\n4. Update weights\n\n## Applications\n\n- Image recognition\n- Natural language processing\n- Predictive analytics\n- Autonomous vehicles',
    'Deep dive into neural networks, their architecture, and training process.',
    (SELECT id FROM users WHERE email = 'admin@example.com'),
    TRUE,
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '2 days'
WHERE NOT EXISTS (SELECT 1 FROM posts WHERE slug = 'understanding-neural-networks');

INSERT INTO posts (title, slug, content, preview, author_id, published, published_at, created_at, updated_at)
SELECT
    'The Future of Machine Learning',
    'future-of-machine-learning',
    E'# The Future is Here\n\nMachine learning continues to evolve at a rapid pace. What can we expect in the coming years?\n\n## Trends to Watch\n\n### AutoML\nAutomated machine learning will make AI accessible to non-experts.\n\n### Edge AI\nML models running directly on devices, reducing latency and improving privacy.\n\n### Explainable AI\nMaking AI decisions more transparent and understandable.\n\n### Quantum Machine Learning\nLeveraging quantum computing for unprecedented computational power.\n\n## Challenges Ahead\n\n- Ethical considerations\n- Data privacy concerns\n- Computational resources\n- Bias in AI systems\n\n## Opportunities\n\nThe future of ML opens doors to:\n- Better healthcare diagnostics\n- Climate change solutions\n- Enhanced cybersecurity\n- Personalized education',
    'Explore the emerging trends and future directions of machine learning technology.',
    (SELECT id FROM users WHERE email = 'admin@example.com'),
    TRUE,
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '5 days'
WHERE NOT EXISTS (SELECT 1 FROM posts WHERE slug = 'future-of-machine-learning');

-- Insert draft post
INSERT INTO posts (title, slug, content, preview, author_id, published, created_at, updated_at)
SELECT
    'Work in Progress: Advanced NLP Techniques',
    'advanced-nlp-techniques-draft',
    'This is a draft post about advanced NLP techniques. Content coming soon...',
    'Draft: Exploring advanced natural language processing techniques.',
    (SELECT id FROM users WHERE email = 'admin@example.com'),
    FALSE,
    NOW(),
    NOW()
WHERE NOT EXISTS (SELECT 1 FROM posts WHERE slug = 'advanced-nlp-techniques-draft');

-- Insert sample comments
INSERT INTO comments (post_id, user_id, content, created_at, updated_at)
SELECT
    (SELECT id FROM posts WHERE slug = 'getting-started-with-ai-development'),
    (SELECT id FROM users WHERE email = 'user1@example.com'),
    'Great introduction! This really helped me understand the basics of AI development.',
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '1 day'
WHERE NOT EXISTS (
    SELECT 1 FROM comments
    WHERE post_id = (SELECT id FROM posts WHERE slug = 'getting-started-with-ai-development')
    AND user_id = (SELECT id FROM users WHERE email = 'user1@example.com')
);

INSERT INTO comments (post_id, user_id, content, created_at, updated_at)
SELECT
    (SELECT id FROM posts WHERE slug = 'getting-started-with-ai-development'),
    (SELECT id FROM users WHERE email = 'user2@example.com'),
    'Could you recommend some good courses for beginners?',
    NOW() - INTERVAL '12 hours',
    NOW() - INTERVAL '12 hours'
WHERE NOT EXISTS (
    SELECT 1 FROM comments
    WHERE post_id = (SELECT id FROM posts WHERE slug = 'getting-started-with-ai-development')
    AND user_id = (SELECT id FROM users WHERE email = 'user2@example.com')
);

-- Insert nested comment (reply)
INSERT INTO comments (post_id, user_id, content, parent_id, created_at, updated_at)
SELECT
    (SELECT id FROM posts WHERE slug = 'getting-started-with-ai-development'),
    (SELECT id FROM users WHERE email = 'admin@example.com'),
    'I recommend starting with Andrew Ng''s Machine Learning course on Coursera. It''s excellent for beginners!',
    (SELECT id FROM comments
     WHERE post_id = (SELECT id FROM posts WHERE slug = 'getting-started-with-ai-development')
     AND user_id = (SELECT id FROM users WHERE email = 'user2@example.com')),
    NOW() - INTERVAL '6 hours',
    NOW() - INTERVAL '6 hours'
WHERE EXISTS (
    SELECT 1 FROM comments
    WHERE post_id = (SELECT id FROM posts WHERE slug = 'getting-started-with-ai-development')
    AND user_id = (SELECT id FROM users WHERE email = 'user2@example.com')
)
AND NOT EXISTS (
    SELECT 1 FROM comments
    WHERE post_id = (SELECT id FROM posts WHERE slug = 'getting-started-with-ai-development')
    AND user_id = (SELECT id FROM users WHERE email = 'admin@example.com')
    AND parent_id IS NOT NULL
);

INSERT INTO comments (post_id, user_id, content, created_at, updated_at)
SELECT
    (SELECT id FROM posts WHERE slug = 'understanding-neural-networks'),
    (SELECT id FROM users WHERE email = 'user3@example.com'),
    'The explanation of backpropagation is very clear. Thanks for sharing!',
    NOW() - INTERVAL '1 day',
    NOW() - INTERVAL '1 day'
WHERE NOT EXISTS (
    SELECT 1 FROM comments
    WHERE post_id = (SELECT id FROM posts WHERE slug = 'understanding-neural-networks')
    AND user_id = (SELECT id FROM users WHERE email = 'user3@example.com')
);

-- Output summary
SELECT 'Seed data inserted successfully!' as message;
SELECT COUNT(*) as total_users FROM users;
SELECT COUNT(*) as total_posts FROM posts;
SELECT COUNT(*) as published_posts FROM posts WHERE published = TRUE;
SELECT COUNT(*) as total_comments FROM comments;
