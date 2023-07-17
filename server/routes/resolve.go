package routes

import (
	"github.com/SazedWorldbringer/url-shortener/database"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

/* database r is used for storing url shorts and r1 is used for managing IP addresses and for rate limiting */

func Resolve(ctx *fiber.Ctx) error {
	url := ctx.Params("url")

	// Instantiate database to retrieve short url
	r := database.CreateClient(0)
	defer r.Close()

	// Retrieve short url
	value, err := r.Get(database.Ctx, url).Result()

	if err == redis.Nil {
		//s Send status 404 if the URL doesn't exist in the database
		ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found in the database",
		})
	} else if err != nil {
		// Send status 500 if there's an internal error
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal error",
		})
	}

	// Instantiate databse for incrementing the count of
	//	how many times this url has been resolved
	r1 := database.CreateClient(1)
	defer r1.Close()

	// Increment counter
	_ = r1.Incr(database.Ctx, "counter")

	// Redirect the user to original URL
	return ctx.Redirect(value, 301)
}
