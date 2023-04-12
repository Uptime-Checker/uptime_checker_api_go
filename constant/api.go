package constant

const (
	GuestUserRateLimitInMinutes            = 5
	GuestUserCodeExpiryInMinutes           = 10
	BearerTokenExpirationInDays            = 180
	FreeSubscriptionDurationInDays         = 120
	TrialSubscriptionDurationInDays        = 14
	MaxMonitorBodySizeInBytes              = 1024
	MinMonitorIntervalInSeconds            = 10
	MaxMonitorIntervalInSeconds            = 86_400
	MinAlarmReminderIntervalInMinutes      = 5
	MaxMonitorTimeoutInSeconds             = 30
	DefaultOrganizationAlarmErrorThreshold = 2
)

// DateCompare type
type DateCompare int

// List of DateCompares
const (
	Date1AfterDate2  DateCompare = 1
	Date1BeforeDate2 DateCompare = -1
	Date1EqualDate2  DateCompare = 0
)

// List of Stripe Events
const (
	StripeCustomerSubscriptionCreated = "customer.subscription.created"
	StripeCustomerSubscriptionUpdated = "customer.subscription.updated"
	StripeCustomerSubscriptionDeleted = "customer.subscription.deleted"
	StripeInvoiceCreated              = "invoice.created"
	StripeInvoicePaid                 = "invoice.paid"
	StripeInvoicePaymentFailed        = "invoice.payment_failed"
)
