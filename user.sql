INSERT INTO users (username, hashed_password, full_name, email) VALUES ($1,$2$3$4) RETURNING *

SELECT * FROM users WHERE username = $1 LIMIT 1
