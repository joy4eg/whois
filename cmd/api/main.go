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
	<style>
		.spinner {
			display: none;
			width: 40px;
			height: 40px;
			border: 4px solid #f3f3f3;
			border-top: 4px solid #3498db;
			border-radius: 50%;
			animation: spin 1s linear infinite;
			margin: 20px auto;
		}
		@keyframes spin {
			0% { transform: rotate(0deg); }
			100% { transform: rotate(360deg); }
		}
	</style>
</head>
<body>
	<h1>WHOIS Lookup</h1>
	<form id="whois-form">
		<label for="host">Host:</label>
		<input type="text" id="host" name="host" required>
		<button type="submit">whois</button>
	</form>
	<div id="spinner" class="spinner"></div>
	<pre id="result"></pre>
	<script>
		document.getElementById('whois-form').addEventListener('submit', async function(event) {
			event.preventDefault();
			const spinner = document.getElementById('spinner');
			const result = document.getElementById('result');

			spinner.style.display = 'block';
			result.textContent = '';

			const host = document.getElementById('host').value;
			try {
				const response = await fetch('/whois', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({ host })
				});
				result.textContent = await response.text();
			} catch (error) {
				result.textContent = 'Error: ' + error.message;
			} finally {
				spinner.style.display = 'none';
			}
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
		c.Response().Header.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
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
