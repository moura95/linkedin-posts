-- Insert categories
INSERT INTO categories (name, parent_id) VALUES ('Hardware', NULL);
INSERT INTO categories (name, parent_id) VALUES ('Software', NULL);
INSERT INTO categories (name, parent_id) VALUES ('Services', NULL);

-- Insert SubCategories Hardware
INSERT INTO categories (name, parent_id) VALUES ('Laptops', 1);
INSERT INTO categories (name, parent_id) VALUES ('Desktops', 1);
INSERT INTO categories (name, parent_id) VALUES ('Mouse', 1);

-- Insert SubCategories Software
INSERT INTO categories (name, parent_id) VALUES ('Operating Systems', 2);
INSERT INTO categories (name, parent_id) VALUES ('Productivity', 2);
INSERT INTO categories (name, parent_id) VALUES ('Development Tools', 2);

-- Insert Subcategories Services
INSERT INTO categories (name, parent_id) VALUES ('Consulting', 3);
INSERT INTO categories (name, parent_id) VALUES ('Support', 3);
INSERT INTO categories (name, parent_id) VALUES ('Maintenance', 3);
