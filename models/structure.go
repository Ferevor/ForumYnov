package models

type Categories struct {
	ID        int
	Name      string
	Parent_id int
}

type Users struct {
	ID            int
	Email         string
	Username      string
	Password_hash string
	Created_at    string
}

type Posts struct {
	ID          int
	User_id     int
	Category_id int
	Title       string
	Content     string
	Created_at  string
}

type Comments struct {
	ID         int
	Post_id    int
	User_id    int
	Content    string
	Created_at string
}

type LikesDislikes struct {
	ID         int
	User_id    int
	Post_id    int
	Comment_id int
	TypeValue  int
}
