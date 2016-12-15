// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package diagnostics

import (
	"encoding/json"
	"runtime"
	"strings"

	"github.com/mattermost/platform/api"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
	"github.com/segmentio/analytics-go"
)

const (
	SEGMENT_KEY    = "fwb7VPbFeQ7SKp3wHm1RzFUuXZudqVok"

	TRACK_CONFIG_SERVICE      = "config_service"
	TRACK_CONFIG_TEAM         = "config_team"
	TRACK_CONFIG_SQL          = "config_sql"
	TRACK_CONFIG_LOG          = "config_log"
	TRACK_CONFIG_FILE         = "config_file"
	TRACK_CONFIG_RATE         = "config_rate"
	TRACK_CONFIG_EMAIL        = "config_email"
	TRACK_CONFIG_PRIVACY      = "config_privacy"
	TRACK_CONFIG_OAUTH        = "config_oauth"
	TRACK_CONFIG_LDAP         = "config_ldap"
	TRACK_CONFIG_COMPLIANCE   = "config_compliance"
	TRACK_CONFIG_LOCALIZATION = "config_localization"
	TRACK_CONFIG_SAML         = "config_saml"
	TRACK_CONFIG_PASSWORD     = "config_password"
	TRACK_CONFIG_CLUSTER      = "config_cluster"
	TRACK_CONFIG_METRICS      = "config_metrics"
	TRACK_CONFIG_WEBRTC       = "config_webrtc"
	TRACK_CONFIG_SUPPORT      = "config_support"
	TRACK_CONFIG_NATIVEAPP    = "config_nativeapp"

	TRACK_ACTIVITY = "activity"
	TRACK_CHANNEL  = "channel"
	TRACK_LICENSE  = "license"
	TRACK_SERVER   = "server"
	TRACK_SESSION  = "session"
	TRACK_TEAM     = "team"
	TRACK_USER     = "user"
)

var client *analytics.Client

func SendDailyDiagnostics() {
	if *utils.Cfg.LogSettings.EnableDiagnostics {
		initDiagnostics()
		trackActivity()
		trackChannels()
		trackConfig()
		trackLicense()
		trackSessions()
		trackServer()
		trackTeams()
		trackUsers()
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

func isSetString(property string) bool {
	if len(property) > 0 {
		return true
	}
	return false
}

func isSetInt(property int64) bool {
	if property > 0 {
		return true
	}
	return false
}

func getPref(name string, prefs model.Preferences) string {
	for _, pref := range prefs {
		if pref.Name == name {
			return pref.Value
		}
	}

	return ""
}

func isDefault(setting interface{}, defaultValue interface{}) bool {
	if setting == defaultValue {
		return true
	}
	return false
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
		"enable_multifactor_authentication":    utils.Cfg.ServiceSettings.EnableMultifactorAuthentication,
		"enable_oauth_service_provider":        utils.Cfg.ServiceSettings.EnableOAuthServiceProvider,
		"connection_security":                  *utils.Cfg.ServiceSettings.ConnectionSecurity,
		"uses_letsencrypt":                     *utils.Cfg.ServiceSettings.UseLetsEncrypt,
		"forward_80_to_443":                    *utils.Cfg.ServiceSettings.Forward80To443,
		"maximum_login_attempts":               utils.Cfg.ServiceSettings.MaximumLoginAttempts,
		"session_length_web_in_days":           *utils.Cfg.ServiceSettings.SessionLengthWebInDays,
		"session_length_mobile_in_days":        *utils.Cfg.ServiceSettings.SessionLengthMobileInDays,
		"session_length_sso_in_days":           *utils.Cfg.ServiceSettings.SessionLengthSSOInDays,
		"session_cache_in_minutes":             *utils.Cfg.ServiceSettings.SessionCacheInMinutes,
		"isdefault_site_url":                   isDefault(*utils.Cfg.ServiceSettings.SiteURL, model.SERVICE_SETTINGS_DEFAULT_SITE_URL),
		"isdefault_tls_cert_file":              isDefault(*utils.Cfg.ServiceSettings.TLSCertFile, model.SERVICE_SETTINGS_DEFAULT_TLS_CERT_FILE),
		"isdefault_tls_key_file":               isDefault(*utils.Cfg.ServiceSettings.TLSKeyFile, model.SERVICE_SETTINGS_DEFAULT_TLS_KEY_FILE),
		"isdefault_read_timeout":               isDefault(*utils.Cfg.ServiceSettings.ReadTimeout, model.SERVICE_SETTINGS_DEFAULT_READ_TIMEOUT),
		"isdefault_write_timeout":              isDefault(*utils.Cfg.ServiceSettings.WriteTimeout, model.SERVICE_SETTINGS_DEFAULT_WRITE_TIMEOUT),
		"isdefault_segment_developer_key":      isDefault(utils.Cfg.ServiceSettings.SegmentDeveloperKey, ""),
		"isdefault_google_developer_key":       isDefault(utils.Cfg.ServiceSettings.GoogleDeveloperKey, ""),
		"isdefault_allow_cors_from":            isDefault(*utils.Cfg.ServiceSettings.AllowCorsFrom, model.SERVICE_SETTINGS_DEFAULT_ALLOW_CORS_FROM),
	})

	SendDiagnostic(TRACK_CONFIG_TEAM, map[string]interface{}{
		"enable_user_creation":                utils.Cfg.TeamSettings.EnableUserCreation,
		"enable_team_creation":                utils.Cfg.TeamSettings.EnableTeamCreation,
		"restrict_team_invite":                *utils.Cfg.TeamSettings.RestrictTeamInvite,
		"restrict_public_channel_management":  *utils.Cfg.TeamSettings.RestrictPublicChannelManagement,
		"restrict_private_channel_management": *utils.Cfg.TeamSettings.RestrictPrivateChannelManagement,
		"enable_open_server":                  *utils.Cfg.TeamSettings.EnableOpenServer,
		"enable_custom_brand":                 *utils.Cfg.TeamSettings.EnableCustomBrand,
		"restrict_direct_message":             *utils.Cfg.TeamSettings.RestrictDirectMessage,
		"max_notifications_per_channel":       *utils.Cfg.TeamSettings.MaxNotificationsPerChannel,
		"max_users_per_team":                  utils.Cfg.TeamSettings.MaxUsersPerTeam,
		"max_channels_per_team":               *utils.Cfg.TeamSettings.MaxChannelsPerTeam,
		"isdefault_site_name":                 isDefault(utils.Cfg.TeamSettings.SiteName, "Mattermost"),
		"isdefault_custom_brand_text":         isDefault(*utils.Cfg.TeamSettings.CustomBrandText, model.TEAM_SETTINGS_DEFAULT_CUSTOM_BRAND_TEXT),
		"isdefault_custom_description_text":   isDefault(*utils.Cfg.TeamSettings.CustomDescriptionText, model.TEAM_SETTINGS_DEFAULT_CUSTOM_DESCRIPTION_TEXT),
		"isdefault_user_status_away_timeout":  isDefault(*utils.Cfg.TeamSettings.UserStatusAwayTimeout, model.TEAM_SETTINGS_DEFAULT_USER_STATUS_AWAY_TIMEOUT),
	})

	SendDiagnostic(TRACK_CONFIG_SQL, map[string]interface{}{
		"driver_name":    utils.Cfg.SqlSettings.DriverName,
		"trace":          utils.Cfg.SqlSettings.Trace,
		"max_idle_conns": utils.Cfg.SqlSettings.MaxIdleConns,
		"max_open_conns": utils.Cfg.SqlSettings.MaxOpenConns,
		"data_source_replicas": len(utils.Cfg.SqlSettings.DataSourceReplicas),
	})

	SendDiagnostic(TRACK_CONFIG_LOG, map[string]interface{}{
		"enable_console":           utils.Cfg.LogSettings.EnableConsole,
		"console_level":            utils.Cfg.LogSettings.ConsoleLevel,
		"enable_file":              utils.Cfg.LogSettings.EnableFile,
		"file_level":               utils.Cfg.LogSettings.FileLevel,
		"enable_webhook_debugging": utils.Cfg.LogSettings.EnableWebhookDebugging,
		"isdefault_file_format":    isDefault(utils.Cfg.LogSettings.FileFormat, ""),
		"isdefault_file_location":  isDefault(utils.Cfg.LogSettings.FileLocation, ""),
	})

	SendDiagnostic(TRACK_CONFIG_PASSWORD, map[string]interface{}{
		"minimum_length": *utils.Cfg.PasswordSettings.MinimumLength,
		"lowercase":      *utils.Cfg.PasswordSettings.Lowercase,
		"number":         *utils.Cfg.PasswordSettings.Number,
		"uppercase":      *utils.Cfg.PasswordSettings.Uppercase,
		"symbol":         *utils.Cfg.PasswordSettings.Symbol,
	})

	SendDiagnostic(TRACK_CONFIG_FILE, map[string]interface{}{
		"enable_public_links": utils.Cfg.FileSettings.EnablePublicLink,
		"driver_name":         utils.Cfg.FileSettings.DriverName,
		"amazon_s3_ssl":       *utils.Cfg.FileSettings.AmazonS3SSL,
		"thumbnail_width":     utils.Cfg.FileSettings.ThumbnailWidth,
		"thumbnail_height":    utils.Cfg.FileSettings.ThumbnailHeight,
		"preview_width":       utils.Cfg.FileSettings.PreviewWidth,
		"preview_height":      utils.Cfg.FileSettings.PreviewHeight,
		"profile_width":       utils.Cfg.FileSettings.ProfileWidth,
		"profile_height":      utils.Cfg.FileSettings.ProfileHeight,
		"max_file_size":       *utils.Cfg.FileSettings.MaxFileSize,
	})

	SendDiagnostic(TRACK_CONFIG_EMAIL, map[string]interface{}{
		"enable_sign_up_with_email":       utils.Cfg.EmailSettings.EnableSignUpWithEmail,
		"enable_sign_in_with_email":       *utils.Cfg.EmailSettings.EnableSignInWithEmail,
		"enable_sign_in_with_username":    *utils.Cfg.EmailSettings.EnableSignInWithUsername,
		"require_email_verification":      utils.Cfg.EmailSettings.RequireEmailVerification,
		"send_email_notifications":        utils.Cfg.EmailSettings.SendEmailNotifications,
		"connection_security":             utils.Cfg.EmailSettings.ConnectionSecurity,
		"send_push_notifications":         *utils.Cfg.EmailSettings.SendPushNotifications,
		"push_notification_contents":      *utils.Cfg.EmailSettings.PushNotificationContents,
		"enable_email_batching":           *utils.Cfg.EmailSettings.EnableEmailBatching,
		"email_batching_buffer_size":      *utils.Cfg.EmailSettings.EmailBatchingBufferSize,
		"email_batching_intervavel":       *utils.Cfg.EmailSettings.EmailBatchingInterval,
		"isdefault_feedback_name":         isDefault(utils.Cfg.EmailSettings.FeedbackName, ""),
		"isdefault_feedback_email":        isDefault(utils.Cfg.EmailSettings.FeedbackEmail, ""),
		"isdefault_feedback_organization": isDefault(*utils.Cfg.EmailSettings.FeedbackOrganization, model.EMAIL_SETTINGS_DEFAULT_FEEDBACK_ORGANIZATION),
	})

	SendDiagnostic(TRACK_CONFIG_RATE, map[string]interface{}{
		"enable_rate_limiter":      *utils.Cfg.RateLimitSettings.Enable,
		"vary_by_remote_address":   utils.Cfg.RateLimitSettings.VaryByRemoteAddr,
		"per_sec":                  utils.Cfg.RateLimitSettings.PerSec,
		"max_burst":                *utils.Cfg.RateLimitSettings.MaxBurst,
		"memory_store_size":        utils.Cfg.RateLimitSettings.MemoryStoreSize,
		"isdefault_vary_by_header": isDefault(utils.Cfg.RateLimitSettings.VaryByHeader, ""),
	})

	SendDiagnostic(TRACK_CONFIG_PRIVACY, map[string]interface{}{
		"show_email_address": utils.Cfg.PrivacySettings.ShowEmailAddress,
		"show_full_name":     utils.Cfg.PrivacySettings.ShowFullName,
	})

	SendDiagnostic(TRACK_CONFIG_OAUTH, map[string]interface{}{
		"enable_gitlab":    utils.Cfg.GitLabSettings.Enable,
		"enable_google":    utils.Cfg.GoogleSettings.Enable,
		"enable_office365": utils.Cfg.Office365Settings.Enable,
	})

	SendDiagnostic(TRACK_CONFIG_SUPPORT, map[string]interface{}{
		"isdefault_terms_of_service_link": isDefault(*utils.Cfg.SupportSettings.TermsOfServiceLink, model.SUPPORT_SETTINGS_DEFAULT_TERMS_OF_SERVICE_LINK),
		"isdefault_privacy_policy_link":   isDefault(*utils.Cfg.SupportSettings.PrivacyPolicyLink, model.SUPPORT_SETTINGS_DEFAULT_PRIVACY_POLICY_LINK),
		"isdefault_about_link":            isDefault(*utils.Cfg.SupportSettings.AboutLink, model.SUPPORT_SETTINGS_DEFAULT_ABOUT_LINK),
		"isdefault_help_link":             isDefault(*utils.Cfg.SupportSettings.HelpLink, model.SUPPORT_SETTINGS_DEFAULT_HELP_LINK),
		"isdefault_report_a_problem_link": isDefault(*utils.Cfg.SupportSettings.ReportAProblemLink, model.SUPPORT_SETTINGS_DEFAULT_REPORT_A_PROBLEM_LINK),
		"isdefault_support_email":         isDefault(*utils.Cfg.SupportSettings.SupportEmail, model.SUPPORT_SETTINGS_DEFAULT_SUPPORT_EMAIL),
	})

	SendDiagnostic(TRACK_CONFIG_LDAP, map[string]interface{}{
		"enable":                         *utils.Cfg.LdapSettings.Enable,
		"connection_security":            *utils.Cfg.LdapSettings.ConnectionSecurity,
		"skip_certificate_verification":  *utils.Cfg.LdapSettings.SkipCertificateVerification,
		"sync_interval_minutes":          *utils.Cfg.LdapSettings.SyncIntervalMinutes,
		"query_timeout":                  *utils.Cfg.LdapSettings.QueryTimeout,
		"max_page_size":                  *utils.Cfg.LdapSettings.MaxPageSize,
		"isdefault_first_name_attribute": isDefault(*utils.Cfg.LdapSettings.FirstNameAttribute, model.LDAP_SETTINGS_DEFAULT_FIRST_NAME_ATTRIBUTE),
		"isdefault_last_name_attribute":  isDefault(*utils.Cfg.LdapSettings.LastNameAttribute, model.LDAP_SETTINGS_DEFAULT_LAST_NAME_ATTRIBUTE),
		"isdefault_email_attribute":      isDefault(*utils.Cfg.LdapSettings.EmailAttribute, model.LDAP_SETTINGS_DEFAULT_EMAIL_ATTRIBUTE),
		"isdefault_username_attribute":   isDefault(*utils.Cfg.LdapSettings.UsernameAttribute, model.LDAP_SETTINGS_DEFAULT_USERNAME_ATTRIBUTE),
		"isdefault_nickname_attribute":   isDefault(*utils.Cfg.LdapSettings.NicknameAttribute, model.LDAP_SETTINGS_DEFAULT_NICKNAME_ATTRIBUTE),
		"isdefault_id_attribute":         isDefault(*utils.Cfg.LdapSettings.IdAttribute, model.LDAP_SETTINGS_DEFAULT_ID_ATTRIBUTE),
		"isdefault_login_field_name":     isDefault(*utils.Cfg.LdapSettings.LoginFieldName, model.LDAP_SETTINGS_DEFAULT_LOGIN_FIELD_NAME),
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
		"isdefault_first_name_attribute": isDefault(*utils.Cfg.SamlSettings.FirstNameAttribute, model.SAML_SETTINGS_DEFAULT_FIRST_NAME_ATTRIBUTE),
		"isdefault_last_name_attribute":  isDefault(*utils.Cfg.SamlSettings.LastNameAttribute, model.SAML_SETTINGS_DEFAULT_LAST_NAME_ATTRIBUTE),
		"isdefault_email_attribute":      isDefault(*utils.Cfg.SamlSettings.EmailAttribute, model.SAML_SETTINGS_DEFAULT_EMAIL_ATTRIBUTE),
		"isdefault_username_attribute":   isDefault(*utils.Cfg.SamlSettings.UsernameAttribute, model.SAML_SETTINGS_DEFAULT_USERNAME_ATTRIBUTE),
		"isdefault_nickname_attribute":   isDefault(*utils.Cfg.SamlSettings.NicknameAttribute, model.SAML_SETTINGS_DEFAULT_NICKNAME_ATTRIBUTE),
		"isdefault_locale_attribute":     isDefault(*utils.Cfg.SamlSettings.LocaleAttribute, model.SAML_SETTINGS_DEFAULT_LOCALE_ATTRIBUTE),
		"isdefault_login_button_text":    isDefault(*utils.Cfg.SamlSettings.LoginButtonText, model.USER_AUTH_SERVICE_SAML_TEXT),
	})

	SendDiagnostic(TRACK_CONFIG_CLUSTER, map[string]interface{}{
		"enable": *utils.Cfg.ClusterSettings.Enable,
	})

	SendDiagnostic(TRACK_CONFIG_METRICS, map[string]interface{}{
		"enable": *utils.Cfg.MetricsSettings.Enable,
	})

	SendDiagnostic(TRACK_CONFIG_NATIVEAPP, map[string]interface{}{
		"isdefault_app_download_link":         isDefault(*utils.Cfg.NativeAppSettings.AppDownloadLink, model.NATIVEAPP_SETTINGS_DEFAULT_APP_DOWNLOAD_LINK),
		"isdefault_android_app_download_link": isDefault(*utils.Cfg.NativeAppSettings.AndroidAppDownloadLink, model.NATIVEAPP_SETTINGS_DEFAULT_ANDROID_APP_DOWNLOAD_LINK),
		"isdefault_iosapp_download_link":      isDefault(*utils.Cfg.NativeAppSettings.IosAppDownloadLink, model.NATIVEAPP_SETTINGS_DEFAULT_IOS_APP_DOWNLOAD_LINK),
	})

	SendDiagnostic(TRACK_CONFIG_WEBRTC, map[string]interface{}{
		"enable":             *utils.Cfg.WebrtcSettings.Enable,
		"isdefault_stun_uri": isDefault(*utils.Cfg.WebrtcSettings.StunURI, model.WEBRTC_SETTINGS_DEFAULT_STUN_URI),
		"isdefault_turn_uri": isDefault(*utils.Cfg.WebrtcSettings.TurnURI, model.WEBRTC_SETTINGS_DEFAULT_TURN_URI),
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
		"public_channels_deleted":  deletedPublicChannelCount,
		"private_channels_deleted": deletedPrivateChannelCount,
		"posts":                    postsCount,
	})
}

func trackChannels() {
	if res := <-api.Srv.Store.Channel().AnalyticsGetAll(); res.Err == nil {
		for _, channel := range res.Data.([]*model.ChannelWithMemberCount) {
			SendDiagnostic(TRACK_CHANNEL, map[string]interface{}{
				"team_id":       channel.TeamId,
				"channel_id":    channel.Id,
				"posts_count":   channel.TotalMsgCount,
				"channel_type":  channel.Type,
				"members_count": channel.MemberCount,
			})
		}
	}
}

func trackLicense() {
	if utils.IsLicensed {
		data := map[string]interface{}{
			"name":    utils.License.Customer.Name,
			"company": utils.License.Customer.Company,
			"issued":  utils.License.IssuedAt,
			"start":   utils.License.StartsAt,
			"expire":  utils.License.ExpiresAt,
			"users":   *utils.License.Features.Users,
		}

		features := utils.License.Features.ToMap()
		for featureName, featureValue := range features {
			data["feature_"+featureName] = featureValue
		}

		SendDiagnostic(TRACK_LICENSE, data)
	}
}

func trackServer() {
	data := map[string]interface{}{
		"edition":          model.BuildEnterpriseReady,
		"version":          model.CurrentVersion,
		"database_type":    utils.Cfg.SqlSettings.DriverName,
		"operating_system": runtime.GOOS,
	}

	if scr := <-api.Srv.Store.User().AnalyticsGetSystemAdminCount(); scr.Err == nil {
		data["system_admins"] = scr.Data.(int64)
	}

	SendDiagnostic(TRACK_SERVER, data)
}

func trackSessions() {
	if res := <-api.Srv.Store.Session().AnalyticsGetAllSessions(); res.Err == nil {
		for _, session := range res.Data.([]*model.Session) {
			data := map[string]interface{}{
				"user_id": session.UserId,
				"session_id": session.Id,
				"first_active_time": session.CreateAt,
				"last_active_time": session.LastActivityAt,
				"os": session.Props["os"],
				"browser": session.Props["browser"],
				"platform": session.Props["platform"],
			}

			SendDiagnostic(TRACK_SESSION, data)
		}
	}
}

func trackTeams() {
	if res := <-api.Srv.Store.Team().GetAll(); res.Err == nil {
		for _, team := range res.Data.([]*model.Team) {
			data := map[string]interface{}{
				"team_id":           team.Id,
				"team_admins_count": 0,
				"open_team":         team.AllowOpenInvite,
			}

			if cor := <-api.Srv.Store.Channel().AnalyticsTypeCount(team.Id, "O"); cor.Err == nil {
				data["channel_count_public"] = cor.Data.(int64)
			}

			if cpr := <-api.Srv.Store.Channel().AnalyticsTypeCount(team.Id, "P"); cpr.Err == nil {
				data["channel_count_private"] = cpr.Data.(int64)
			}

			if cdr := <-api.Srv.Store.Channel().AnalyticsTypeCount(team.Id, "D"); cdr.Err == nil {
				data["channel_count_direct"] = cdr.Data.(int64)
			}

			if pcr := <-api.Srv.Store.Post().AnalyticsPostCount(team.Id, false, false); pcr.Err == nil {
				data["posts_count"] = pcr.Data.(int64)
			}

			if tcr := <-api.Srv.Store.User().AnalyticsUniqueUserCount(team.Id); tcr.Err == nil {
				data["users_count"] = tcr.Data.(int64)
			}

			if tar := <-api.Srv.Store.Team().AnalyticsGetTeamAdminCount(team.Id); tar.Err == nil {
				data["team_admins_count"] = tar.Data.(int64)
			}

			SendDiagnostic(TRACK_TEAM, data)
		}
	}
}

func trackUsers() {
	if res := <-api.Srv.Store.User().AnalyticsGetUsersWithTeamCount(); res.Err == nil {
		for _, user := range res.Data.([]*model.UserWithTeamCount) {
			data := map[string]interface{}{
				"user_id":                                user.Id,
				"teams_joined":                           user.TeamCount,
				"first_name_set":                         isSetString(user.FirstName),
				"last_name_set":                          isSetString(user.LastName),
				"nickname_set":                           isSetString(user.Nickname),
				"profile_picture_set":                    isSetInt(user.LastPictureUpdate),
				"mfa_activated":                          user.MfaActive,
				"signin_method":                          user.AuthService,
				"language":                               user.Locale,
				"send_desktop_notifications":             user.NotifyProps["desktop"],
				"desktop_notifications_sound":            user.NotifyProps["desktop_sound"],
				"desktop_notifications_duration":         user.NotifyProps["desktop_duration"],
				"email_notifications":                    user.NotifyProps["email"],
				"push_notifications_activity":            user.NotifyProps["push"],
				"push_notifications_status":              user.NotifyProps["push_status"],
				"notifications_trigger_first_name":       user.NotifyProps["first_name"],
				"notifications_trigger_channel_mentions": user.NotifyProps["channel"],
				"reply_notifications":                    user.NotifyProps["comments"],
			}

			mentionTriggers := strings.Split(user.NotifyProps["mention_keys"], ",")
			for _, trigger := range mentionTriggers {
				if trigger == user.Username {
					data["notifications_trigger_username"] = "true"
				} else if trigger == "@"+user.Username {
					data["notifications_trigger_at_username"] = "true"
				} else {
					data["notifications_trigger_other"] = "true"
				}
			}

			if pur := <-api.Srv.Store.Channel().AnalyticsTypeCountForUser(user.Id, "O"); pur.Err == nil {
				data["public_channels_joined"] = pur.Data.(int64)
			}

			if prr := <-api.Srv.Store.Channel().AnalyticsTypeCountForUser(user.Id, "P"); prr.Err == nil {
				data["private_channels_joined"] = prr.Data.(int64)
			}

			if dmr := <-api.Srv.Store.Channel().AnalyticsTypeCountForUser(user.Id, "O"); dmr.Err == nil {
				data["direct_channels_joined"] = dmr.Data.(int64)
			}

			if tpr := <-api.Srv.Store.Preference().GetCategory(user.Id, model.PREFERENCE_CATEGORY_THEME); tpr.Err == nil {
				themeData := getPref("", tpr.Data.(model.Preferences))

				decoder := json.NewDecoder(strings.NewReader(themeData))
				var theme map[string]string
				if err := decoder.Decode(&theme); err == nil {
					data["theme"] = theme["type"]
				}
			}

			if dpr := <-api.Srv.Store.Preference().GetCategory(user.Id, model.PREFERENCE_CATEGORY_DISPLAY_SETTINGS); dpr.Err == nil {
				prefs := dpr.Data.(model.Preferences)
				data["display_font"] = getPref("selected_font", prefs)
				data["24_hour_clock"] = getPref("use_military_time", prefs)
				data["teammate_name_display"] = getPref("name_format", prefs)
				data["collapse_link_previews"] = getPref("collapse_previews", prefs)
				data["message_display"] = getPref("message_display", prefs)
				data["channel_display_mode"] = getPref("channel_display_mode", prefs)
			}

			if opr := <-api.Srv.Store.Preference().GetCategory(user.Id, model.PREFERENCE_CATEGORY_AUTHORIZED_OAUTH_APP); opr.Err == nil {
				data["oauth_authorized_apps_count"] = len(opr.Data.(model.Preferences))
			}

			if apr := <-api.Srv.Store.Preference().GetCategory(user.Id, model.PREFERENCE_CATEGORY_ADVANCED_SETTINGS); apr.Err == nil {
				prefs := apr.Data.(model.Preferences)
				data["advanced_send_message_ctrl_enter"] = getPref("send_on_ctrl_enter", prefs)
				data["advanced_enable_post_formatting"] = getPref("formatting", prefs)
				data["advanced_enable_join_leave_messages"] = getPref("join_leave", prefs)
				data["feature_enabled_embed_preview"] = getPref("feature_enabled_embed_preview", prefs)
				data["feature_enabled_markdown_preview"] = getPref("feature_enabled_markdown_preview", prefs)
				data["feature_enabled_webrtc_preview"] = getPref("feature_enabled_webrtc_preview", prefs)
			}

			SendDiagnostic(TRACK_USER, data)
		}
	}
}
