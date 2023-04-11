DROP TABLE IF EXISTS Movie;
USE myDB;
CREATE TABLE IF NOT EXISTS Movie (
	movie_id int NOT NULL,
  	movie_name varchar(200) NOT NULL,
  	movie_year varchar(4) NOT NULL DEFAULT '9999',
  	movie_genre varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'genre',
  	movie_duration varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'xx h yy m',
	movie_origin varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'earth',
	movie_director varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'human',
	movie_rating float NOT NULL DEFAULT '10',
	movie_rating_count bigint NOT NULL DEFAULT '0',
	movie_link varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'url'
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
ALTER TABLE Movie
ADD PRIMARY KEY (movie_id);
ALTER TABLE Movie
MODIFY movie_id int NOT NULL AUTO_INCREMENT