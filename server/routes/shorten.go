package routes

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/SazedWorldbringer/url-shortener/database"
	"github.com/SazedWorldbringer/url-shortener/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Request schema
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// Response schema
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

/* database r is used for storing url shorts and r1 is used for managing IP addresses and for rate limiting */

func Shorten(ctx *fiber.Ctx) error {
	body := &request{}

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	// Rate limiting

	// Instantiate database for rate limiting
	r1 := database.CreateClient(1)
	defer r1.Close()

	// retrieve ip address value
	val, err := r1.Get(database.Ctx, ctx.IP()).Result()
	// retrive time to live value for ip and discard any error
	limit, _ := r1.TTL(database.Ctx, ctx.IP()).Result()

	if err == redis.Nil {
		// if the IP doesn't exist in the database (this is the first time this user is using the service)

		// set a quote for the current IP address, which will be deleted after
		_ = r1.Set(database.Ctx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else if err == nil {
		// if the IP does exist in the database

		// convert value received from the database to an integer
		valInt, _ := strconv.Atoi(val)
		// if the quota has been exhausted
		// return status code of 503, and show how long till it resets the rate limit
		if valInt <= 0 {
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// TODO: Validate request

	// Shorten

	// Instantiate database for storing URLs
	r := database.CreateClient(0)
	defer r.Close()

	// Custom short id for URL
	var id string

	if body.CustomShort == "" {
		id = utils.GenerateRandomString(rand.Uint64())
	} else {
		id = body.CustomShort
	}

	// Check if given custom short already exists in the database, return status code 403 if it does
	val, _ = r.Get(database.Ctx, id).Result()

	if val != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL Custom short already in use.",
		})
	}

	// If the expiry field in request body isn't defined, set the url to be expired in 24 hours
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// Set custom short id to the given URL, with an expiration duration of 24 hours
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	// return 500 for unexpected internal error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to the server.",
		})
	}

	// Response
	remainingQuota, _ := r1.Decr(database.Ctx, ctx.IP()).Result()
	remainingRateLimitReset := int(limit / time.Nanosecond / time.Minute)
	customShortUrl := os.Getenv("DOMAIN") + "/" + id

	// Response struct
	resp := response{
		URL:             body.URL,
		CustomShort:     customShortUrl,
		Expiry:          body.Expiry,
		XRateRemaining:  int(remainingQuota),
		XRateLimitReset: time.Duration(remainingRateLimitReset),
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
