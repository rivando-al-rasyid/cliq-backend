package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func CliqRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	// txRepo := repository.NewCliqRepo(db)
	// txServ := service.NewCliqService(txRepo, rdb)
	// txCont := controller.NewAuthController(txServ)

	// cliq := router.Group("/cliq", middleware.VerifyTokenWithDB(db))

}
