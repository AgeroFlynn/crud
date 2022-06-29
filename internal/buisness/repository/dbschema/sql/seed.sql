INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
                                                                                               ('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN,USER}', '$2a$10$pDrzO6UaEHJMb8nniy4QNOkZLOK09.HqTJrTQTBnEIoFNMwMvqn3a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
                                                                                               ('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$pDrzO6UaEHJMb8nniy4QNOkZLOK09.HqTJrTQTBnEIoFNMwMvqn3a', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO phone_dict (phone_dict_id, user_id, telegram, date_created, date_updated) VALUES
                                                                                                 ('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', '5cf37266-3473-4006-984f-9325122678b7', '@admin', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
                                                                                                 ('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '@user', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
ON CONFLICT DO NOTHING;