INSERT INTO users
    ("name", "nick", "email", "password")
VALUES
    ('Teste 1', 'teste1', 'teste1@teste.com.br', '$2a$10$ex66rvkqyOYVNxEhp9ilj.oK61DwCNsCv58pAT5uch9AgObr/5NAS'),
    ('Teste 2', 'teste2', 'teste2@teste.com.br', '$2a$10$ex66rvkqyOYVNxEhp9ilj.oK61DwCNsCv58pAT5uch9AgObr/5NAS'),
    ('Teste 3', 'teste3', 'teste3@teste.com.br', '$2a$10$ex66rvkqyOYVNxEhp9ilj.oK61DwCNsCv58pAT5uch9AgObr/5NAS'),
    ('Teste 4', 'teste4', 'teste4@teste.com.br', '$2a$10$ex66rvkqyOYVNxEhp9ilj.oK61DwCNsCv58pAT5uch9AgObr/5NAS')


INSERT INTO followers
    ("user_id", "follower_id")
VALUES
    (1, 2),
    (1, 3)