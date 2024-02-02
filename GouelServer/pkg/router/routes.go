package routes

import (
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/config"
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/handlers"
	middlewares "github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/middleware"
	"github.com/Gouel/GouelMonoRepo/tree/main/GouelServer/pkg/token"
	"github.com/gin-gonic/gin"
)

/*-> == à vérifier*/

func Routes(router *gin.Engine, cfg config.Config) {
	router.Use(middlewares.CORSMiddleware(), middlewares.JWTMiddleware(cfg.JWTSecretKey))

	// Routes TOKEN JWT
	router.GET("/token/view", token.ViewTokenRoute)
	router.POST("/token/auth", token.AuthRoute)
	router.POST("/token/refresh", token.RefreshRoute)
	/*->*/ router.POST("/token/auth/ticket", token.AuthRouteTicket)

	// Route de test de bon fonctionnement
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "serveur actif"})
	})

	//Routes USERS
	router.GET("/users/search/email/:email", middlewares.RoleAuthorizationMiddleware("API"), handlers.FindUsersByEmailStartsWithHandler)
	router.GET("/users/:user_id", handlers.GetUserByIdHandler)
	router.POST("/users/", middlewares.RoleAuthorizationMiddleware("API"), handlers.CreateUserHandler)
	router.POST("/users/event/:event_id", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("caisse"), handlers.CreateUserHandler)
	router.POST("/users/transaction/:event_id/:user_id", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("buvette", "caisse"), handlers.AddUserTransactionHandler)
	router.POST("/users/pay/:event_id/:ticket_id", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("buvette"), handlers.UserPayHandler)
	router.PUT("/users/:user_id", middlewares.RoleAuthorizationMiddleware("API"), handlers.UpdateUserHandler)

	//Routes TICKETS
	router.GET("/tickets/:event_id", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("entree"), handlers.GetAllTicketsFromEventHandler)
	router.GET("/tickets/:event_id/page/:page", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("entree"), handlers.GetPaginatedTicketsFromEventHandler)
	router.GET("/tickets/:event_id/ticket/:ticket_id", middlewares.EventAccessMiddleware(), handlers.GetTicketInfoHandler)
	router.POST("/tickets/:event_id/:ticket_code", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("caisse"), handlers.CreateTicketHandler)
	router.PUT("/tickets/:event_id/sam", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("entree", ""), handlers.SetSamHandler)
	router.POST("/tickets/:event_id/validate", middlewares.EventAuthorizationMiddleware("entree"), handlers.ValidateTicketHandler)
	router.DELETE("/tickets/:ticket_id", middlewares.RoleAuthorizationMiddleware("API"), handlers.DeleteTicketHandler)

	//Routes EVENTS
	router.GET("/events/:event_id/smtp", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("caisse"), handlers.GetSMTPConfiguration)
	router.GET("/events", handlers.GetAccessibleEventsHandler)
	router.GET("/events/:event_id", middlewares.EventAccessMiddleware(), handlers.GetEventByIdHandler)
	router.GET("/events/:event_id/simple", handlers.GetSimpleEventHandler)
	router.GET("/events/:event_id/products", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("buvette"), handlers.GetEventProductsHandler)
	router.GET("/events/:event_id/products/:product_code", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("buvette"), handlers.GetEventProductsCodeHandler)
	router.GET("/events/:event_id/admins", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN", "ADMIN"), handlers.GetEventAdminsHandler)
	router.GET("/events/:event_id/volunteers", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN", "ADMIN"), handlers.GetEventVolunteersHandler)
	router.GET("/events/:event_id/lockers", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("vestiaire"), handlers.GetEventLockersHandler)
	router.GET("/events/:event_id/tickets", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN", "ADMIN"), handlers.GetEventTicketsHandler)
	router.POST("/events/", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN"), handlers.AddEventHandler)
	router.POST("/events/:event_id/volunteers", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN", "ADMIN"), handlers.AddVolunteerToEventHandler)
	router.POST("/events/:event_id/tickets", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN", "ADMIN"), handlers.AddTicketToEventHandler)
	router.POST("/events/:event_id/lockers", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN", "ADMIN"), handlers.AddLockerToEventHandler)
	router.POST("/events/:event_id/lockers/range", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN", "ADMIN"), handlers.AddRangeOfLockersHandler)
	router.POST("/events/:event_id/products", middlewares.RoleAuthorizationMiddleware("API", "SUPERADMIN", "ADMIN"), handlers.AddProductToEventHandler)
	router.PUT("/events/:event_id", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.UpdateEventHandler)
	router.PUT("/events/:event_id/lockers", middlewares.EventAccessMiddleware(), middlewares.EventAuthorizationMiddleware("vestiaire"), handlers.UpdateLockerHandler)
	router.PUT("/events/:event_id/products/:product_code", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.UpdateProductHandler)
	router.PUT("/events/:event_id/volunteers", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.UpdateVolunteerHandler)
	router.PUT("/events/:event_id/tickets/:ticket_code", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.UpdateTicketHandler)
	router.DELETE("/events/:event_id/tickets/:ticket_code", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.DeleteEventTicketHandler)
	router.DELETE("/events/:event_id/lockers/:locker_code", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.DeleteLockerHandler)
	router.DELETE("/events/:event_id/lockers", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.DeleteAllLockersHandler)
	router.DELETE("/events/:event_id/volunteers/:user_id", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.DeleteVolunteerHandler)
	router.DELETE("/events/:event_id/products/:product_code", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API", "ADMIN"), handlers.DeleteProductHandler)
	router.DELETE("/events/:event_id", middlewares.RoleAuthorizationMiddleware("SUPERADMIN", "API"), handlers.DeleteEventHandler)
}
