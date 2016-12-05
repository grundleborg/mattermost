// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package diagnostics

import (
	"runtime"

	"github.com/mattermost/platform/api"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
	"github.com/segmentio/analytics-go"
)

const (
	DIAGNOSTIC_URL = "https://d7zmvsa9e04kk.cloudfront.net"
	SEGMENT_KEY    = "ua1qQtmgOZWIM23YjD842tQAsN7Ydi5X"

	PROP_DIAGNOSTIC_ID                = "id"
	PROP_DIAGNOSTIC_CATEGORY          = "c"
	VAL_DIAGNOSTIC_CATEGORY_DEFAULT   = "d"
	PROP_DIAGNOSTIC_BUILD             = "b"
	PROP_DIAGNOSTIC_ENTERPRISE_READY  = "be"
	PROP_DIAGNOSTIC_DATABASE          = "db"
	PROP_DIAGNOSTIC_OS                = "os"
	PROP_DIAGNOSTIC_USER_COUNT        = "uc"
	PROP_DIAGNOSTIC_TEAM_COUNT        = "tc"
	PROP_DIAGNOSTIC_ACTIVE_USER_COUNT = "auc"
	PROP_DIAGNOSTIC_UNIT_TESTS        = "ut"

	TRACK_CONFIG_SERVICE      = "service"
	TRACK_CONFIG_TEAM         = "team"
	TRACK_CONFIG_SQL          = "sql"
	TRACK_CONFIG_LOG          = "log"
	TRACK_CONFIG_FILE         = "file"
	TRACK_CONFIG_RATE         = "rate"
	TRACK_CONFIG_EMAIL        = "email"
	TRACK_CONFIG_PRIVACY      = "privacy"
	TRACK_CONFIG_OAUTH        = "oauth"
	TRACK_CONFIG_LDAP         = "ldap"
	TRACK_CONFIG_COMPLIANCE   = "compliance"
	TRACK_CONFIG_LOCALIZATION = "localization"
	TRACK_CONFIG_SAML         = "saml"

	TRACK_LICENSE  = "license"
	TRACK_ACTIVITY = "activity"
	TRACK_VERSION  = "version"
)

var client *analytics.Client

func SendDailyDiagnostics() {
	if *utils.Cfg.LogSettings.EnableDiagnostics {
		initDiagnostics()
		trackActivity()
		trackConfig()
		trackLicense()
		trackVersion()
	}
}

func initDiagnostics() {
	if client == nil {
		client = analytics.New(SEGMENT_KEY)
		client.Identify(&analytics.Identify{
			UserId: utils.CfgDiagnosticId,
		})
	}
}

func SendDiagnostic(event string, properties map[string]interface{}) {
	client.Track(&analytics.Track{
		Event:      event,
		UserId:     utils.CfgDiagnosticId,
		Properties: properties,
	})
}

func trackConfig() {
	SendDiagnostic(TRACK_CONFIG_SERVICE, map[string]interface{}{
		"web_server_mode":                      *utils.Cfg.ServiceSettings.WebserverMode,
		"enable_security_fix_alert":            *utils.Cfg.ServiceSettings.EnableSecurityFixAlert,
		"enable_insecure_outgoing_connections": *utils.Cfg.ServiceSettings.EnableInsecureOutgoingConnections,
		"enable_incoming_webhooks":             utils.Cfg.ServiceSettings.EnableIncomingWebhooks,
		"enable_outgoing_webhooks":             utils.Cfg.ServiceSettings.EnableOutgoingWebhooks,
		"enable_commands":                      *utils.Cfg.ServiceSettings.EnableCommands,
		"enable_only_admin_integrations":       *utils.Cfg.ServiceSettings.EnableOnlyAdminIntegrations,
		"enable_post_username_override":        utils.Cfg.ServiceSettings.EnablePostUsernameOverride,
		"enable_post_icon_override":            utils.Cfg.ServiceSettings.EnablePostIconOverride,
		"enable_custom_emoji":                  *utils.Cfg.ServiceSettings.EnableCustomEmoji,
		"restrict_custom_emoji_creation":       *utils.Cfg.ServiceSettings.RestrictCustomEmojiCreation,
		"enable_testing":                       utils.Cfg.ServiceSettings.EnableTesting,
		"enable_developer":                     *utils.Cfg.ServiceSettings.EnableDeveloper,
	})

	SendDiagnostic(TRACK_CONFIG_TEAM, map[string]interface{}{
		"enable_user_creation":                utils.Cfg.TeamSettings.EnableUserCreation,
		"enable_team_creation":                utils.Cfg.TeamSettings.EnableTeamCreation,
		"restrict_team_invite":                *utils.Cfg.TeamSettings.RestrictTeamInvite,
		"restrict_public_channel_management":  *utils.Cfg.TeamSettings.RestrictPublicChannelManagement,
		"restrict_private_channel_management": *utils.Cfg.TeamSettings.RestrictPrivateChannelManagement,
		"enable_open_server":                  *utils.Cfg.TeamSettings.EnableOpenServer,
		"enable_custom_brand":                 *utils.Cfg.TeamSettings.EnableCustomBrand,
	})

	SendDiagnostic(TRACK_CONFIG_SQL, map[string]interface{}{
		"driver_name": utils.Cfg.SqlSettings.DriverName,
	})

	SendDiagnostic(TRACK_CONFIG_LOG, map[string]interface{}{
		"enable_console":           utils.Cfg.LogSettings.EnableConsole,
		"console_level":            utils.Cfg.LogSettings.ConsoleLevel,
		"enable_file":              utils.Cfg.LogSettings.EnableFile,
		"file_level":               utils.Cfg.LogSettings.FileLevel,
		"enable_webhook_debugging": utils.Cfg.LogSettings.EnableWebhookDebugging,
	})

	SendDiagnostic(TRACK_CONFIG_FILE, map[string]interface{}{
		"enable_public_links": utils.Cfg.FileSettings.EnablePublicLink,
	})

	SendDiagnostic(TRACK_CONFIG_RATE, map[string]interface{}{
		"enable_rate_limiter":    *utils.Cfg.RateLimitSettings.Enable,
		"vary_by_remote_address": utils.Cfg.RateLimitSettings.VaryByRemoteAddr,
	})

	SendDiagnostic(TRACK_CONFIG_EMAIL, map[string]interface{}{
		"enable_sign_up_with_email":    utils.Cfg.EmailSettings.EnableSignUpWithEmail,
		"enable_sign_in_with_email":    *utils.Cfg.EmailSettings.EnableSignInWithEmail,
		"enable_sign_in_with_username": *utils.Cfg.EmailSettings.EnableSignInWithUsername,
		"require_email_verification":   utils.Cfg.EmailSettings.RequireEmailVerification,
		"send_email_notifications":     utils.Cfg.EmailSettings.SendEmailNotifications,
		"connection_security":          utils.Cfg.EmailSettings.ConnectionSecurity,
		"send_push_notifications":      *utils.Cfg.EmailSettings.SendPushNotifications,
		"push_notification_contents":   *utils.Cfg.EmailSettings.PushNotificationContents,
	})

	SendDiagnostic(TRACK_CONFIG_PRIVACY, map[string]interface{}{
		"show_email_address": utils.Cfg.PrivacySettings.ShowEmailAddress,
		"show_full_name":     utils.Cfg.PrivacySettings.ShowFullName,
	})

	SendDiagnostic(TRACK_CONFIG_OAUTH, map[string]interface{}{
		"gitlab":    utils.Cfg.GitLabSettings.Enable,
		"google":    utils.Cfg.GoogleSettings.Enable,
		"office365": utils.Cfg.Office365Settings.Enable,
	})

	SendDiagnostic(TRACK_CONFIG_LDAP, map[string]interface{}{
		"enable":                        *utils.Cfg.LdapSettings.Enable,
		"connection_security":           *utils.Cfg.LdapSettings.ConnectionSecurity,
		"skip_certificate_verification": *utils.Cfg.LdapSettings.SkipCertificateVerification,
	})

	SendDiagnostic(TRACK_CONFIG_COMPLIANCE, map[string]interface{}{
		"enable":       *utils.Cfg.ComplianceSettings.Enable,
		"enable_daily": *utils.Cfg.ComplianceSettings.EnableDaily,
	})

	SendDiagnostic(TRACK_CONFIG_LOCALIZATION, map[string]interface{}{
		"default_server_locale": *utils.Cfg.LocalizationSettings.DefaultServerLocale,
		"default_client_locale": *utils.Cfg.LocalizationSettings.DefaultClientLocale,
		"available_locales":     *utils.Cfg.LocalizationSettings.AvailableLocales,
	})

	SendDiagnostic(TRACK_CONFIG_SAML, map[string]interface{}{
		"enable": *utils.Cfg.SamlSettings.Enable,
	})
}

func trackActivity() {
	var userCount int64
	var activeUserCount int64
	var inactiveUserCount int64
	var teamCount int64
	var publicChannelCount int64
	var privateChannelCount int64
	var directChannelCount int64
	var deletedPublicChannelCount int64
	var deletedPrivateChannelCount int64
	var postsCount int64

	if ucr := <-api.Srv.Store.User().GetTotalUsersCount(); ucr.Err == nil {
		userCount = ucr.Data.(int64)
	}

	if ucr := <-api.Srv.Store.Status().GetTotalActiveUsersCount(); ucr.Err == nil {
		activeUserCount = ucr.Data.(int64)
	}

	if iucr := <-api.Srv.Store.Status().GetTotalActiveUsersCount(); iucr.Err == nil {
		inactiveUserCount = iucr.Data.(int64)
	}

	if tcr := <-api.Srv.Store.Team().AnalyticsTeamCount(); tcr.Err == nil {
		teamCount = tcr.Data.(int64)
	}

	if ucc := <-api.Srv.Store.Channel().AnalyticsTypeCount("", "O"); ucc.Err == nil {
		publicChannelCount = ucc.Data.(int64)
	}

	if pcc := <-api.Srv.Store.Channel().AnalyticsTypeCount("", "P"); pcc.Err == nil {
		privateChannelCount = pcc.Data.(int64)
	}

	if dcc := <-api.Srv.Store.Channel().AnalyticsTypeCount("", "D"); dcc.Err == nil {
		directChannelCount = dcc.Data.(int64)
	}

	if duccr := <-api.Srv.Store.Channel().AnalyticsTypeCount("", "O"); duccr.Err == nil {
		deletedPublicChannelCount = duccr.Data.(int64)
	}

	if dpccr := <-api.Srv.Store.Channel().AnalyticsTypeCount("", "P"); dpccr.Err == nil {
		deletedPrivateChannelCount = dpccr.Data.(int64)
	}

	if pcr := <-api.Srv.Store.Post().AnalyticsPostCount("", false, false); pcr.Err == nil {
		postsCount = pcr.Data.(int64)
	}

	SendDiagnostic(TRACK_ACTIVITY, map[string]interface{}{
		"registered_users":         userCount,
		"active_users":             activeUserCount,
		"inactive_users":           inactiveUserCount,
		"teams":                    teamCount,
		"public_channels":          publicChannelCount,
		"private_channels":         privateChannelCount,
		"direct_message_channels":  directChannelCount,
		"deleted_public_channels":  deletedPublicChannelCount,
		"deleted_private_channels": deletedPrivateChannelCount,
		"posts":                    postsCount,
	})

	// TODO: registered inactive accounts.
}

func trackLicense() {
	if utils.IsLicensed {
		SendDiagnostic(TRACK_LICENSE, map[string]interface{}{
			"name":     utils.License.Customer.Name,
			"company":  utils.License.Customer.Company,
			"issued":   utils.License.IssuedAt,
			"start":    utils.License.StartsAt,
			"expire":   utils.License.ExpiresAt,
			"users":    *utils.License.Features.Users,
			"features": utils.License.Features.ToMap(),
		})
	}
}

func trackVersion() {
	edition := model.BuildEnterpriseReady
	version := model.CurrentVersion
	database := utils.Cfg.SqlSettings.DriverName
	operatingSystem := runtime.GOOS

	SendDiagnostic(TRACK_VERSION, map[string]interface{}{
		"edition":          edition,
		"version":          version,
		"database":         database,
		"operating_system": operatingSystem,
	})
}
