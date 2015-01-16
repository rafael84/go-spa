/*
cfg package loads a .env file with the environment settings.

The .env file must have the following structure:

	# Server settings
	export SERVER_HOST="http://localhost"
	export SERVER_PORT="3001"
	export SERVER_FRONTEND_PATH="../../frontend"
	export SERVER_API_PREFIX="/api/v1"
	export SERVER_PRIVKEY="keys/app.rsa"    # openssl genrsa -out app.rsa 2048
	export SERVER_PUBKEY="keys/app.rsa.pub" # openssl rsa -in app.rsa -pubout > app.rsa.pub

	# Email settings
	export EMAIL_FROM="GO-SPA"
	export EMAIL_SUBJECT="Reset your password"
	export EMAIL_USERNAME="******@gmail.com"
	export EMAIL_PASSWORD="******"

	# Database
	export DB_USER=postgres
	export DB_NAME=postgres
	export DB_PASSWORD=1234
	export DB_HOST=127.0.0.1
	export DB_PORT=5432
	export DB_SSLMODE=disable
	export DB_CONN_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
*/
package cfg
