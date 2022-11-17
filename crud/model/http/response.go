package http

type CreateJSON struct {
	Success bool
	Message string
}

var Cors = map[string]string{
	"Access-Control-Allow-Headers": "*",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "*",
}
