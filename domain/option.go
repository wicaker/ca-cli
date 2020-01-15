package domain

// Option struct, collection for option input
type Option struct {
	Title       string
	Description string
}

var (
	// GoPg database handler option
	GoPg = "gopg"
	// Gorm database handler option
	Gorm = "gorm"
	// Sqlx dabase handler option
	Sqlx = "sqlx"
	// SQL standard database handle option
	SQL = "sql"
	// Echo server handler option
	Echo = "echo"
	// Gin server handler option
	Gin = "gin"
	// GorillaMux server handler option
	GorillaMux = "gorilla mux"
	// NetHTTP standard server handler option
	NetHTTP = "net/http"
	// Graphql transport option
	Graphql = "graphql"
	// Grpc transport option
	Grpc = "grpc"
)
