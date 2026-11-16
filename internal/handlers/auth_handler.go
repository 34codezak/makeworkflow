/*
This is the HTTP layer. 
Handler parses the request, validates input, calls a service, and writes the JSON response and status code. 
	It knows about HTTP but not about SQL.
*/

// register, login, refresh, etc.