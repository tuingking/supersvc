package user

const (
	getUserQuery    = `SELECT id, name, phone, email, status, created_at FROM user`
	countUserQuery  = `SELECT COUNT(1) FROM user`
	createUserQuery = `INSERT INTO user(id, name, phone, email, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`
)
