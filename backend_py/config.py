import os

keycloak_client_secret = "00000000–0000–0000–0000–000000000000"  # nosec
keycloak_client_id = "backend_py"

keycloak_base_url = os.getenv("KEYCLOAK_BASE_URL", "http://localhost:8080")
backend_go_url = os.getenv("BACKEND_GO_URL", "http://localhost:3333")
