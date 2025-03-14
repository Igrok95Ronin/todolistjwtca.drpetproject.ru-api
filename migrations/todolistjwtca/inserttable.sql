-- Вставляем записи в таблицу users
INSERT INTO users (user_name, email, password_hash) VALUES
                                                        ('User1', 'user1@example.com', 'hash1'),
                                                        ('User2', 'user2@example.com', 'hash2'),
                                                        ('User3', 'user3@example.com', 'hash3');

-- Вставляем записи в таблицу all_notes
INSERT INTO all_notes (note, completed, user_id, created_at) VALUES
                                                                 ('Note 1', false, 1, CURRENT_TIMESTAMP),
                                                                 ('Note 2', true, 2, CURRENT_TIMESTAMP),
                                                                 ('Note 3', false, 3, CURRENT_TIMESTAMP),
                                                                 ('Note 4', true, 1, CURRENT_TIMESTAMP),
                                                                 ('Note 5', false, 2, CURRENT_TIMESTAMP),
                                                                 ('Note 6', true, 3, CURRENT_TIMESTAMP),
                                                                 ('Note 7', false, 1, CURRENT_TIMESTAMP),
                                                                 ('Note 8', true, 2, CURRENT_TIMESTAMP),
                                                                 ('Note 9', false, 3, CURRENT_TIMESTAMP),
                                                                 ('Note 10', true, 1, CURRENT_TIMESTAMP),
                                                                 ('Note 11', false, 2, CURRENT_TIMESTAMP),
                                                                 ('Note 12', true, 3, CURRENT_TIMESTAMP),
                                                                 ('Note 13', false, 1, CURRENT_TIMESTAMP),
                                                                 ('Note 14', true, 2, CURRENT_TIMESTAMP),
                                                                 ('Note 15', false, 3, CURRENT_TIMESTAMP),
                                                                 ('Note 16', true, 1, CURRENT_TIMESTAMP),
                                                                 ('Note 17', false, 2, CURRENT_TIMESTAMP),
                                                                 ('Note 18', true, 3, CURRENT_TIMESTAMP),
                                                                 ('Note 19', false, 1, CURRENT_TIMESTAMP),
                                                                 ('Note 20', true, 2, CURRENT_TIMESTAMP);
