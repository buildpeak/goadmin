-- Populate the "user" table
INSERT INTO "user" (username, email, password, first_name, last_name)
VALUES ('user1', 'user1@gmail.com', 'password1', 'User', 'One'),
       ('user2', 'user2@gmail.com', 'password2', 'User', 'Two');

-- Populate the "role" table
INSERT INTO role (name)
VALUES ('Admin'),
       ('Editor'),
       ('Viewer');

-- Get two user IDs
DO $$
DECLARE user1_id bigint;
DECLARE user2_id bigint;
BEGIN
SELECT id INTO user1_id FROM "user" WHERE username = 'user1';
SELECT id INTO user2_id FROM "user" WHERE username = 'user2';

-- Populate the "user_role" table with the retrieved user IDs
INSERT INTO user_role (user_id, role_id)
VALUES (user1_id, 1),
       (user2_id, 2);
END$$;

-- Populate the "permission" table
INSERT INTO permission (name, rule_type, rule)
VALUES ('Permission1', 'type1', 'rule1'),
       ('Permission2', 'type2', 'rule2'),
       ('Permission3', 'type3', 'rule3');

-- Get permission IDs
DO $$
DECLARE user1_id bigint;
DECLARE user2_id bigint;
DECLARE perm1_id bigint;
DECLARE perm2_id bigint;
DECLARE perm3_id bigint;
BEGIN
SELECT id INTO user1_id FROM "user" WHERE username = 'user1';
SELECT id INTO user2_id FROM "user" WHERE username = 'user2';
SELECT id INTO perm1_id FROM permission WHERE name = 'Permission1';
SELECT id INTO perm2_id FROM permission WHERE name = 'Permission2';
SELECT id INTO perm3_id FROM permission WHERE name = 'Permission3';

-- Populate the "role_permission" table with the retrieved permission IDs
INSERT INTO role_permission (role_id, permission_id)
VALUES (1, perm1_id),
       (2, perm2_id),
       (3, perm3_id);

-- Populate the "user_permission" table with the retrieved user and permission IDs
INSERT INTO user_permission (user_id, permission_id)
VALUES (user1_id, perm1_id),
       (user2_id, perm2_id);
END$$;
