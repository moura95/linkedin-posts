CREATE TABLE IF NOT EXISTS tickets (
                                       id SERIAL PRIMARY KEY,
                                       title VARCHAR(255) NOT NULL,
                                       description TEXT,
                                       status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
                                       severity_id INT NOT NULL,
                                       category_id INT NOT NULL,
                                       subcategory_id INT,
                                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                       updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                       completed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
                                          id SERIAL PRIMARY KEY,
                                          name VARCHAR(255) NOT NULL,
                                          parent_id INT,
                                          FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS severities (
                                          id SERIAL PRIMARY KEY,
                                          name VARCHAR(50) NOT NULL,
                                          description TEXT
);

INSERT INTO severities (id, name, description) VALUES
                                                   (1, 'HIGH', 'Problemas críticos que requerem atenção imediata'),
                                                   (2, 'MEDIUM', 'Problemas importantes que devem ser resolvidos em breve'),
                                                   (3, 'LOW', 'Problemas de menor importância que podem esperar'),
                                                   (4, 'FEATURE', 'Solicitação de nova funcionalidade')
ON CONFLICT (id) DO NOTHING;

ALTER TABLE tickets
    ADD CONSTRAINT fk_ticket_severity
        FOREIGN KEY (severity_id) REFERENCES severities(id);

ALTER TABLE tickets
    ADD CONSTRAINT fk_ticket_category
        FOREIGN KEY (category_id) REFERENCES categories(id);

ALTER TABLE tickets
    ADD CONSTRAINT fk_ticket_subcategory
        FOREIGN KEY (subcategory_id) REFERENCES categories(id);

CREATE INDEX IF NOT EXISTS idx_ticket_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_ticket_severity ON tickets(severity_id);
CREATE INDEX IF NOT EXISTS idx_ticket_category ON tickets(category_id);



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


INSERT INTO tickets (title, status, description, severity_id, category_id, subcategory_id) VALUES
                                                                                               ('Login Fails with Correct Credentials', 'OPEN', 'Users can''t log in despite correct credentials', 1, 2, 7),
                                                                                               ('Printer Not Responding', 'OPEN', 'Office printer is not responding to any commands', 2, 1, 4),
                                                                                               ('Software Update Required', 'IN_PROGRESS', 'Latest software version needs to be installed', 3, 2, 8),
                                                                                               ('System Crash on Boot', 'BLOCKED', 'System crashes during boot process', 1, 1, 5),
                                                                                               ('Request for New Workstation', 'DONE', 'New workstation required for the marketing department', 4, 1, 6),
                                                                                               ('Email Server Down', 'IN_PROGRESS', 'Email server is down and needs immediate attention', 1, 3, 10),
                                                                                               ('Slow Internet Connection', 'OPEN', 'Internet speed is significantly slower than usual', 3, 3, 9),
                                                                                               ('Bug in Accounting Software', 'IN_PROGRESS', 'Error in the accounting software affecting reports', 2, 2, 9),
                                                                                               ('Consultation for Network Upgrade', 'DONE', 'Need a consultation for upgrading the network infrastructure', 4, 3, 10),
                                                                                               ('Monitor Flickering', 'OPEN', 'Monitor flickers intermittently during use', 3, 1, 6);
