package docs

var IndexHtml = []byte(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>API Documentation</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>

  <body>
    <div id="app"></div>

    <!-- Load the Script -->
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>

    <!-- Initialize the Scalar API Reference -->
    <script>
      Scalar.createApiReference('#app', {
        // The URL of the OpenAPI/Swagger document
        url: '/api/v1/docs/swagger.yaml',
      })
    </script>
  </body>
</html>`)
