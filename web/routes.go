package web

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/cron"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/watchdog"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/worker"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

func SetupRoutes(ctx context.Context, app *fiber.App) {
	// Default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// API V1
	v1 := app.Group("/v1")
	v1.Use(middlelayer.Header())
	v1.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	//  ========== Age of the domains ==========
	jobDomain := domain.NewJobDomain()
	userDomain := domain.NewUserDomain()
	paymentDomain := domain.NewPaymentDomain()
	organizationDomain := domain.NewOrganizationDomain()

	monitorDomain := domain.NewMonitorDomain()
	checkDomain := domain.NewCheckDomain()
	monitorStatusDomain := domain.NewMonitorStatusDomain()
	alarmPolicyDomain := domain.NewAlarmPolicyDomain()
	monitorRegionDomain := domain.NewMonitorRegionDomain()
	regionDomain := domain.NewRegionDomain()
	assertionDomain := domain.NewAssertionDomain()

	//  ========== Age of the services ==========
	authService := service.NewAuthService(userDomain)
	userService := service.NewUserService(userDomain)
	paymentService := service.NewPaymentService(paymentDomain)
	organizationService := service.NewOrganizationService(organizationDomain, alarmPolicyDomain)
	monitorService := service.NewMonitorService(monitorDomain, monitorStatusDomain)
	assertionService := service.NewAssertionService(assertionDomain)
	monitorRegionService := service.NewMonitorRegionService(monitorRegionDomain)
	checkService := service.NewCheckService(checkDomain)

	//  ========== Age of the modules ==========
	// Setup Watchdog
	dog := watchdog.NewWatchDog(
		checkDomain,
		regionDomain,
		assertionDomain,
		monitorRegionDomain,
		monitorStatusDomain,
		checkService,
		monitorService,
		monitorRegionService,
	)

	// Setup Tasks
	syncProductsTask := task.NewSyncProductsTask()
	runCheckTask := task.NewRunCheckTask()

	cogman := cron.NewCron(jobDomain, regionDomain, monitorDomain, monitorRegionDomain, syncProductsTask)
	wheel := worker.NewWorker(runCheckTask)

	//  ========== Age of the routers ==========
	// User router for auth and user account
	userRouter := v1.Group("/user")
	registerUserHandlers(userRouter, userDomain, authService, userService)

	// Organization router for managing organization
	orgRouter := v1.Group("/organization")
	registerOrganizationHandlers(
		orgRouter,
		userDomain,
		paymentDomain,
		organizationDomain,
		authService,
		paymentService,
		organizationService,
	)

	// Monitor router for managing monitor
	monitorRouter := v1.Group("/monitor")
	registerMonitorHandlers(
		monitorRouter,
		monitorDomain,
		authService,
		monitorService,
		assertionService,
		dog,
	)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
	app.Hooks().OnListen(func() error {
		if err := cogman.Start(ctx); err != nil {
			return err
		}
		return wheel.Start(ctx)
	})
}

func registerUserHandlers(
	router fiber.Router,
	userDomain *domain.UserDomain,
	authService *service.AuthService,
	userService *service.UserService,
) {
	auth := middlelayer.Protected(authService)

	handler := controller.NewUserController(userDomain, authService, userService)
	router.Post("/guest", handler.CreateGuestUser)
	router.Post("/guest/login", handler.GuestUserLogin)
	router.Post("/provider/login", handler.ProviderLogin)

	router.Get("/me", auth, handler.GetMe)
	router.Patch("/", auth, handler.Update)
}

func registerOrganizationHandlers(
	router fiber.Router,
	userDomain *domain.UserDomain,
	paymentDomain *domain.PaymentDomain,
	organizationDomain *domain.OrganizationDomain,
	authService *service.AuthService,
	paymentService *service.PaymentService,
	organizationService *service.OrganizationService,
) {
	auth := middlelayer.Protected(authService)

	handler := controller.NewOrganizationController(
		userDomain,
		paymentDomain,
		organizationDomain,
		paymentService,
		organizationService,
	)

	router.Post("/", auth, handler.CreateOrganization)
	router.Get("/list", auth, handler.ListOrganizationsOfUser)
}

func registerMonitorHandlers(
	router fiber.Router,
	monitorDomain *domain.MonitorDomain,
	authService *service.AuthService,
	monitorService *service.MonitorService,
	assertionService *service.AssertionService,
	dog *watchdog.WatchDog,
) {
	auth := middlelayer.Protected(authService)

	handler := controller.NewMonitorController(
		monitorDomain,
		monitorService,
		assertionService,
		dog,
	)

	router.Post("/", auth, handler.Create)
	router.Post("/dry", auth, handler.DryRun)
	router.Get("/list", auth, handler.ListMonitors)
}
