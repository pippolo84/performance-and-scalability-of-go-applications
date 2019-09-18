package social

// endpoints for social data

// users
const usersURI = "https://jsonplaceholder.typicode.com/users"

// to simulate fetching data from more than one endpoints, we use
// fragments (#1, #2, #3, etc.). Doing this we pretend each entry
// of these array is a different data source
// Fragments have no effect for the test endpoint
// so we are fetching the same data set over and over

// posts
var postsURIs = [...]string{
	"https://jsonplaceholder.typicode.com/posts#1",
	"https://jsonplaceholder.typicode.com/posts#2",
	"https://jsonplaceholder.typicode.com/posts#3",
	"https://jsonplaceholder.typicode.com/posts#4",
	"https://jsonplaceholder.typicode.com/posts#5",
	"https://jsonplaceholder.typicode.com/posts#6",
	"https://jsonplaceholder.typicode.com/posts#7",
	"https://jsonplaceholder.typicode.com/posts#8",
}

// comments
var commentsURIs = [...]string{
	"https://jsonplaceholder.typicode.com/comments#1",
	"https://jsonplaceholder.typicode.com/comments#2",
	"https://jsonplaceholder.typicode.com/comments#3",
	"https://jsonplaceholder.typicode.com/comments#4",
	"https://jsonplaceholder.typicode.com/comments#5",
	"https://jsonplaceholder.typicode.com/comments#6",
	"https://jsonplaceholder.typicode.com/comments#7",
	"https://jsonplaceholder.typicode.com/comments#8",
}

// types for social data

type (
	Geo struct {
		Lat string `json:"lat"`
		Lng string `json:"lng"`
	}

	Address struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     Geo    `json:"Geo"`
	}

	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	}

	User struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Address  Address `json:"address"`
		Phone    string  `json:"phone"`
		Website  string  `json:"website"`
		Company  Company `json:"company"`
	}
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:id`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}
