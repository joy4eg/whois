package main

import (
	"flag"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/joy4eg/whois"
)

// appHTML is a simple HTML template for the app.
const appHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WHOIS Lookup</title>
</head>
<body>
    <h1>WHOIS Lookup</h1>
    <form id="whois-form">
        <label for="host">Host:</label>
        <input type="text" id="host" name="host" required>
        <button type="submit">whois</button>
    </form>
    <pre id="result"></pre>
    <script>
        document.getElementById('whois-form').addEventListener('submit', async function(event) {
            event.preventDefault();
            const host = document.getElementById('host').value;
            const response = await fetch('/whois', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ host })
            });
            const result = await response.text();
            document.getElementById('result').textContent = result;
        });
    </script>
</body>
</html>
`

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	client, err := whois.New(whois.WithCache(time.Hour))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		c.Response().Header.Set("Content-Type", "text/html")
		return c.SendString(appHTML)
	})
	app.Post("/whois", func(c fiber.Ctx) error {
		var data struct {
			Host string `json:"host"`
		}
		if err := c.Bind().Body(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
		}
		now := time.Now()
		slog.InfoContext(c.Context(), "new whois request", "host", data.Host)
		result, err := client.Whois(c.Context(), data.Host)
		slog.InfoContext(c.Context(), "whois request completed", "host", data.Host, "duration", time.Since(now), "err", err)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendString(result)
	})

	log.Fatal(app.Listen(":" + strconv.Itoa(*port)))
}
