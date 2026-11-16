/*
Reads environment variables (DB URL, port, JWT secret, etc.) 
	into a typed Config struct so the rest of the app doesn’t call os.Getenv directly.
*/