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

	propertyDomain := domain.NewPropertyDomain()
	monitorDomain := domain.NewMonitorDomain()
	checkDomain := domain.NewCheckDomain()
	monitorStatusDomain := domain.NewMonitorStatusDomain()
	alarmPolicyDomain := domain.NewAlarmPolicyDomain()
	monitorRegionDomain := domain.NewMonitorRegionDomain()
	regionDomain := domain.NewRegionDomain()
	assertionDomain := domain.NewAssertionDomain()
	errorLogDomain := domain.NewErrorLogDomain()
	dailyReportDomain := domain.NewDailyReportDomain()
	alarmDomain := domain.NewAlarmDomain()
	productDomain := domain.NewProductDomain()
	planDomain := domain.NewPlanDomain()
	alarmChannelDomain := domain.NewAlarmChannelDomain()
	monitorIntegrationDomain := domain.NewMonitorIntegrationDomain()
	monitorNotificationDomain := domain.NewMonitorNotificationDomain()

	//  ========== Age of the services ==========
	authService := service.NewAuthService(userDomain)
	userService := service.NewUserService(userDomain)
	paymentService := service.NewPaymentService(userDomain, paymentDomain)
	organizationService := service.NewOrganizationService(organizationDomain, alarmPolicyDomain)
	monitorService := service.NewMonitorService(monitorDomain, monitorStatusDomain)
	assertionService := service.NewAssertionService(assertionDomain)
	monitorRegionService := service.NewMonitorRegionService(monitorRegionDomain)
	checkService := service.NewCheckService(checkDomain)
	errorLogService := service.NewErrorLogService(errorLogDomain)
	dailyReportService := service.NewDailyReportService(dailyReportDomain)
	alarmPolicyService := service.NewAlarmPolicyService(alarmPolicyDomain)
	productService := service.NewProductService(productDomain)
	propertyService := service.NewPropertyService(propertyDomain)
	monitorIntegrationService := service.NewMonitorIntegrationService(monitorIntegrationDomain)

	//  ========== Age of the modules ==========
	// Setup Watchdog
	dog := watchdog.NewWatchDog(
		checkDomain,
		regionDomain,
		monitorDomain,
		monitorRegionDomain,
		monitorStatusDomain,
		monitorIntegrationDomain,
		alarmDomain,
		alarmChannelDomain,
		monitorNotificationDomain,
		checkService,
		monitorService,
		monitorRegionService,
		errorLogService,
		dailyReportService,
		alarmPolicyService,
	)

	// Setup Tasks
	syncProductsTask := task.NewSyncProductsTask(planDomain, productService)
	runCheckTask := task.NewRunCheckTask(dog, monitorDomain, monitorRegionDomain)
	startMonitorTask := task.NewStartMonitorTask(dog, monitorDomain, regionDomain, assertionDomain)

	cogman := cron.NewCron(
		jobDomain,
		regionDomain,
		monitorDomain,
		monitorRegionDomain,
		propertyService,
		syncProductsTask,
	)
	wheel := worker.NewWorker(runCheckTask, startMonitorTask)

	//  ========== Age of the routers ==========

	// start - versioned route
	regionController := controller.NewRegionController(regionDomain)
	v1.Get("/regions", regionController.List)
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
		dog,
		monitorDomain,
		regionDomain,
		monitorRegionDomain,
		authService,
		monitorService,
		assertionService,
	)

	productRouter := v1.Group("/product")
	registerProductHandlers(productRouter, productDomain, userDomain, authService)

	integrationRouter := v1.Group("/integration")
	registerIntegrationHandlers(
		integrationRouter,
		alarmChannelDomain,
		monitorIntegrationDomain,
		monitorIntegrationService,
		authService,
	)

	alarmChannelRouter := v1.Group("/alarm_channel")
	registerAlarmChannelHandlers(alarmChannelRouter, alarmChannelDomain, authService)
	// end - versioned route

	webhookRouter := app.Group("/webhook")
	registerWebhookHandlers(webhookRouter, paymentService)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
	app.Hooks().OnListen(func() error {
		if err := cogman.Start(ctx); err != nil {
			panic(err)
		}
		return wheel.StartGue(ctx)
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
	dog *watchdog.WatchDog,
	monitorDomain *domain.MonitorDomain,
	regionDomain *domain.RegionDomain,
	monitorRegionDomain *domain.MonitorRegionDomain,
	authService *service.AuthService,
	monitorService *service.MonitorService,
	assertionService *service.AssertionService,
) {
	auth := middlelayer.Protected(authService)

	handler := controller.NewMonitorController(
		dog,
		monitorDomain,
		regionDomain,
		monitorRegionDomain,
		monitorService,
		assertionService,
	)

	router.Post("/", auth, handler.Create)
	router.Post("/start", auth, handler.Start)
	router.Post("/dry", auth, handler.DryRun)
	router.Get("/list", auth, handler.ListMonitors)
	router.Get("/:id", auth, handler.Get)
	router.Get("/:id/all", auth, handler.GetAll)
}

func registerProductHandlers(
	router fiber.Router,
	productDomain *domain.ProductDomain,
	userDomain *domain.UserDomain,
	authService *service.AuthService,
) {
	auth := middlelayer.Protected(authService)

	handler := controller.NewProductController(productDomain, userDomain)

	router.Get("/list/external", auth, handler.ListExternal)
	router.Get("/list/internal", handler.ListInternal)
	router.Get("/billing/customer", auth, handler.CreateBillingCustomer)
}

func registerWebhookHandlers(
	router fiber.Router,
	paymentService *service.PaymentService,
) {
	handler := controller.NewWebhookController(paymentService)

	router.Post("/stripe", handler.StripePayment)
}

func registerIntegrationHandlers(
	router fiber.Router,
	alarmChannelDomain *domain.AlarmChannelDomain,
	monitorIntegrationDomain *domain.MonitorIntegrationDomain,
	monitorIntegrationService *service.MonitorIntegrationService,
	authService *service.AuthService,
) {
	auth := middlelayer.Protected(authService)

	handler := controller.NewMonitorIntegrationController(
		alarmChannelDomain,
		monitorIntegrationDomain,
		monitorIntegrationService,
	)

	router.Post("/", auth, handler.Create)
	router.Get("/list", auth, handler.List)
}

func registerAlarmChannelHandlers(
	router fiber.Router,
	alarmChannelDomain *domain.AlarmChannelDomain,
	authService *service.AuthService,
) {
	auth := middlelayer.Protected(authService)

	handler := controller.NewAlarmChannelController(alarmChannelDomain)

	router.Get("/list", auth, handler.List)
}
