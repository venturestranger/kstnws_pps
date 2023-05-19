create table posts (
	id BIGSERIAL PRIMARY KEY,
	id_author BIGINT NOT NULL,
	title VARCHAR(200) NOT NULL,
	lead VARCHAR(200) NOT NULL,
	picture_url VARCHAR(300) NOT NULL,
	content VARCHAR(100000) NOT NULL,
	date_publication TIMESTAMP NOT NULL,
	date_edit TIMESTAMP NOT NULL,
	category VARCHAR(50) NOT NULL,
	hashtags VARCHAR(1000) NOT NULL
);
