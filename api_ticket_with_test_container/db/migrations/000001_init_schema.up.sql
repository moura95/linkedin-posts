SET TIME ZONE 'America/Sao_Paulo';


CREATE TABLE severities_levels (
                            id SERIAL PRIMARY KEY,
                            description VARCHAR(255) NOT NULL
                        );

INSERT INTO severities_levels (id,description) VALUES
                                                (1, 'issue high'),
                                                (2, 'high'),
                                                (3, 'medium'),
                                                (4, 'low');

CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL,
                            parent_id INT,
                            FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE tickets (
                         id SERIAL PRIMARY KEY,
                         title VARCHAR(60) NOT NULL,
                         status VARCHAR(20) DEFAULT 'OPEN' NOT NULL CHECK (status IN ('OPEN', 'IN_PROGRESS','BLOCKED','DONE', 'CLOSED')),
                         description TEXT NOT NULL,
                         severity_id INT NOT NULL ,
                         category_id INT NOT NULL,
                         subcategory_id INT,
                         created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
                         completed_at TIMESTAMP,
                         FOREIGN KEY (severity_id) REFERENCES severities_levels(id),
                         FOREIGN KEY (category_id) REFERENCES categories(id),
                         FOREIGN KEY (subcategory_id) REFERENCES categories(id));




CREATE INDEX idx_status ON tickets(status);
