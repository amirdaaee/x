package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-gost/x/config"
	"github.com/go-gost/x/config/parsing"
	"github.com/go-gost/x/registry"
)

// swagger:parameters createRateLimiterRequest
type createRateLimiterRequest struct {
	// in: body
	Data config.LimiterConfig `json:"data"`
}

// successful operation.
// swagger:response createRateLimiterResponse
type createRateLimiterResponse struct {
	Data Response
}

func createRateLimiter(ctx *gin.Context) {
	// swagger:route POST /config/rlimiters Limiter createRateLimiterRequest
	//
	// Create a new rate limiter, the name of limiter must be unique in limiter list.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: createRateLimiterResponse

	var req createRateLimiterRequest
	ctx.ShouldBindJSON(&req.Data)

	if req.Data.Name == "" {
		writeError(ctx, ErrInvalid)
		return
	}

	v := parsing.ParseRateLimiter(&req.Data)

	if err := registry.RateLimiterRegistry().Register(req.Data.Name, v); err != nil {
		writeError(ctx, ErrDup)
		return
	}

	cfg := config.Global()
	cfg.CLimiters = append(cfg.CLimiters, &req.Data)
	config.SetGlobal(cfg)

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}

// swagger:parameters updateRateLimiterRequest
type updateRateLimiterRequest struct {
	// in: path
	// required: true
	Limiter string `uri:"limiter" json:"limiter"`
	// in: body
	Data config.LimiterConfig `json:"data"`
}

// successful operation.
// swagger:response updateRateLimiterResponse
type updateRateLimiterResponse struct {
	Data Response
}

func updateRateLimiter(ctx *gin.Context) {
	// swagger:route PUT /config/rlimiters/{limiter} Limiter updateRateLimiterRequest
	//
	// Update rate limiter by name, the limiter must already exist.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: updateRateLimiterResponse

	var req updateRateLimiterRequest
	ctx.ShouldBindUri(&req)
	ctx.ShouldBindJSON(&req.Data)

	if !registry.RateLimiterRegistry().IsRegistered(req.Limiter) {
		writeError(ctx, ErrNotFound)
		return
	}

	req.Data.Name = req.Limiter

	v := parsing.ParseRateLimiter(&req.Data)

	registry.RateLimiterRegistry().Unregister(req.Limiter)

	if err := registry.RateLimiterRegistry().Register(req.Limiter, v); err != nil {
		writeError(ctx, ErrDup)
		return
	}

	cfg := config.Global()
	for i := range cfg.Limiters {
		if cfg.Limiters[i].Name == req.Limiter {
			cfg.Limiters[i] = &req.Data
			break
		}
	}
	config.SetGlobal(cfg)

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}

// swagger:parameters deleteRateLimiterRequest
type deleteRateLimiterRequest struct {
	// in: path
	// required: true
	Limiter string `uri:"limiter" json:"limiter"`
}

// successful operation.
// swagger:response deleteRateLimiterResponse
type deleteRateLimiterResponse struct {
	Data Response
}

func deleteRateLimiter(ctx *gin.Context) {
	// swagger:route DELETE /config/rlimiters/{limiter} Limiter deleteRateLimiterRequest
	//
	// Delete rate limiter by name.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: deleteRateLimiterResponse

	var req deleteRateLimiterRequest
	ctx.ShouldBindUri(&req)

	if !registry.RateLimiterRegistry().IsRegistered(req.Limiter) {
		writeError(ctx, ErrNotFound)
		return
	}
	registry.RateLimiterRegistry().Unregister(req.Limiter)

	cfg := config.Global()
	limiteres := cfg.Limiters
	cfg.Limiters = nil
	for _, s := range limiteres {
		if s.Name == req.Limiter {
			continue
		}
		cfg.Limiters = append(cfg.Limiters, s)
	}
	config.SetGlobal(cfg)

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}
