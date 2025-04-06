package main

type DatabaseBuilder interface {
	SetHost(host string) DatabaseBuilder
	SetPort(port int) DatabaseBuilder
	SetCredentials(user, password string) DatabaseBuilder
	SetDatabase(dbName string) DatabaseBuilder
	Build() (Database, error)
}

type PostgresBuilder struct {
	config PostgresConfig
}

func NewPostgresBuilder() *PostgresBuilder {
	return &PostgresBuilder{
		config: PostgresConfig{
			Port:    5432,
			SSLMode: "disable",
		},
	}
}

func (b *PostgresBuilder) SetHost(host string) DatabaseBuilder {
	b.config.Host = host
	return b
}

func (b *PostgresBuilder) SetPort(port int) DatabaseBuilder {
	b.config.Port = port
	return b
}

func (b *PostgresBuilder) SetCredentials(user, password string) DatabaseBuilder {
	b.config.User = user
	b.config.Password = password
	return b
}

func (b *PostgresBuilder) SetDatabase(dbName string) DatabaseBuilder {
	b.config.Database = dbName
	return b
}

func (b *PostgresBuilder) SetSSLMode(sslMode string) *PostgresBuilder {
	b.config.SSLMode = sslMode
	return b
}

func (b *PostgresBuilder) Build() (Database, error) {
	if b.config.Host == "" || b.config.Database == "" {
		return nil, ConfigurationErr
	}

	return &PostgresDatabase{
		BaseDatabase: BaseDatabase{
			Name: "PostgreSQL",
		},
		Config: b.config,
	}, nil
}

type MongoBuilder struct {
	config MongoConfig
}

func NewMongoBuilder() *MongoBuilder {
	return &MongoBuilder{
		config: MongoConfig{
			Port: 27017,
		},
	}
}

func (b *MongoBuilder) SetHost(host string) DatabaseBuilder {
	b.config.Host = host
	return b
}

func (b *MongoBuilder) SetPort(port int) DatabaseBuilder {
	b.config.Port = port
	return b
}

func (b *MongoBuilder) SetCredentials(user, password string) DatabaseBuilder {
	b.config.User = user
	b.config.Password = password
	return b
}

func (b *MongoBuilder) SetDatabase(dbName string) DatabaseBuilder {
	b.config.Database = dbName
	return b
}

func (b *MongoBuilder) Build() (Database, error) {
	if b.config.Host == "" || b.config.Database == "" {
		return nil, ConfigurationErr
	}

	return &MongoDatabase{
		BaseDatabase: BaseDatabase{
			Name: "MongoDB",
		},
		Config: b.config,
	}, nil
}
