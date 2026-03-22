# /status - Estado del proyecto SpainCoin

Muestra el estado actual del proyecto: tests, build, nodo en producción.

```bash
echo "=== SpainCoin Status ===" && \
echo "" && \
echo "--- Tests ---" && \
CGO_ENABLED=0 go test ./... 2>&1 | grep -E "^(ok|FAIL)" && \
echo "" && \
echo "--- Build ---" && \
CGO_ENABLED=0 go build ./... 2>&1 && echo "Go: OK" && \
npm run build --prefix frontend 2>&1 | tail -3 && \
echo "" && \
echo "--- Nodo en producción ---" && \
curl -s http://204.168.176.40:8545/status 2>/dev/null | python3 -m json.tool || echo "Nodo no accesible" && \
echo "" && \
echo "--- Estructura ---" && \
find . -name "*.go" -not -path "*/vendor/*" | wc -l && echo "archivos Go"
```
