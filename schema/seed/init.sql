-- Insert Regions
INSERT INTO region (name, key, "default")
VALUES('Sunnyvale, California (US)', 'sjc', true);
INSERT INTO region (name, key, "default")
VALUES('Frankfurt, Germany', 'fra', true);
INSERT INTO region (name, key, "default")
VALUES('Singapore', 'sin', false);
INSERT INTO region (name, key, "default")
VALUES('Sydney, Australia', 'syd', false);
INSERT INTO region (name, key, "default")
VALUES('Toronto, Canada', 'yyz', false);
-- Insert Roles
INSERT INTO role (name, type)
VALUES('Super Admin', 1);
INSERT INTO role (name, type)
VALUES('Admin', 2);
INSERT INTO role (name, type)
VALUES('Editor', 3);
INSERT INTO role (name, type)
VALUES('Member', 4);
-- Insert Products
INSERT INTO product (name, description, tier)
VALUES('Free', 'Free for lifetime', 1);
-- Insert Features
INSERT INTO feature (name, type)
VALUES('API_CHECK_COUNT', 1);
INSERT INTO feature (name, type)
VALUES('API_CHECK_INTERVAL', 1);
INSERT INTO feature (name, type)
VALUES('USER_COUNT', 10);
-- Insert Plans
INSERT INTO plan (price, type, product_id)
VALUES(0, 1, 1);
-- Monthly/Free
-- Insert Product Features
INSERT INTO product_feature (count, feature_id, product_id)
VALUES(5, 1, 1);
-- Free/API_CHECK_COUNT/5
INSERT INTO product_feature (count, feature_id, product_id)
VALUES(300, 2, 1);
-- Free/API_CHECK_INTERVAL/300
INSERT INTO product_feature (count, feature_id, product_id)
VALUES(1, 3, 1);
-- Free/USER_COUNT/1
-- Insert Role Claims
INSERT INTO role_claim (claim, role_id)
VALUES('CREATE_RESOURCE', 3);
-- Editor
INSERT INTO role_claim (claim, role_id)
VALUES('UPDATE_RESOURCE', 3);
INSERT INTO role_claim (claim, role_id)
VALUES('DELETE_RESOURCE', 3);
INSERT INTO role_claim (claim, role_id)
VALUES('CREATE_RESOURCE', 2);
-- Admin
INSERT INTO role_claim (claim, role_id)
VALUES('UPDATE_RESOURCE', 2);
INSERT INTO role_claim (claim, role_id)
VALUES('DELETE_RESOURCE', 2);
INSERT INTO role_claim (claim, role_id)
VALUES('INVITE_USER', 2);
INSERT INTO role_claim (claim, role_id)
VALUES('CREATE_RESOURCE', 1);
-- Super Admin
INSERT INTO role_claim (claim, role_id)
VALUES('UPDATE_RESOURCE', 1);
INSERT INTO role_claim (claim, role_id)
VALUES('DELETE_RESOURCE', 1);
INSERT INTO role_claim (claim, role_id)
VALUES('INVITE_USER', 1);
INSERT INTO role_claim (claim, role_id)
VALUES('BILLING', 1);
INSERT INTO role_claim (claim, role_id)
VALUES('DESTROY_ORGANIZATION', 1);
-- Job
INSERT INTO job (status, "on", name, interval, recurring)
VALUES(1, true, 'SYNC_STRIPE_PRODUCTS', 60, true);
INSERT INTO job (status, "on", name, interval, recurring)
VALUES(1, true, 'CHECK_WATCHDOG', 35, true);
-- Property
INSERT INTO property (key, value)
VALUES('WATCHDOG', 'true');