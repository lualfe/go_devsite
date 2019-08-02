DROP TABLE IF EXISTS dev CASCADE;
DROP TABLE IF EXISTS empresa CASCADE;
DROP TABLE IF EXISTS vaga CASCADE;
DROP TABLE IF EXISTS candidatura;
DROP TABLE IF EXISTS skills CASCADE;
DROP TABLE IF EXISTS vaga_skills;
DROP TABLE IF EXISTS dev_skills;

CREATE TABLE dev (
	id SERIAL PRIMARY KEY,
	nome VARCHAR(50),
	email VARCHAR(100) NOT NULL,
	senha VARCHAR(255),
	descricao TEXT,
	stack INTEGER UNIQUE,
	exp INTEGER UNIQUE,
	provider VARCHAR(30) UNIQUE,
	avatar VARCHAR(100),
	dev_uuid VARCHAR(200) UNIQUE,
	oauthuserid VARCHAR(100),
	username VARCHAR(100) UNIQUE
);

CREATE TABLE empresa (
	id SERIAL PRIMARY KEY,
	nome VARCHAR(50),
	email VARCHAR(60) NOT NULL,
	senha VARCHAR(255) NOT NULL
);

CREATE TABLE vaga (
	id SERIAL PRIMARY KEY,
	titulo VARCHAR(50),
	descricao TEXT,
	empresa_id INTEGER REFERENCES empresa (id),
	stack_id INTEGER REFERENCES dev (stack),
	exp_id INTEGER REFERENCES dev (exp)
);

CREATE TABLE candidatura (
	id SERIAL PRIMARY KEY,
	vaga_id INTEGER REFERENCES vaga (id),
	dev_id INTEGER REFERENCES dev (id)
);

CREATE TABLE skills (
	id SERIAL PRIMARY KEY,
	skill VARCHAR(40)
);

CREATE TABLE vaga_skills (
	id INTEGER,
	vaga_id INTEGER REFERENCES vaga (id),
	skills_id INTEGER REFERENCES skills (id)
);

CREATE TABLE dev_skills (
	id INTEGER,
	dev_id INTEGER REFERENCES dev (id),
	skills_id INTEGER REFERENCES skills (id)
);
