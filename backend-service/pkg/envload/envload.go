package envload

import "github.com/joho/godotenv"

// Load reads env files so commands work from the monorepo root (gate-crown/) or from backend-service/.
// Tries ./.env then backend-service/.env so the service file wins when both exist under the repo root.
func Load() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("backend-service/.env")
}
