package lang

var En = map[string]map[string]string{
	"common": {
		"GenericError":              "An error occurred, please try again later.",
		"PermissionDenied":          "You don't have permission to perform this operation.",
		"RequestParameterIncorrect": "The request parameters are incorrect.",
		"ModificationFailed":        "Modification failed, please try again later.",
		"DeletionFailed":            "Deletion failed, please try again later.",
		"DoItToSelf":                "You cannot operate on yourself.",
		"OperationFailed":           "Operation failed, please try again later.",
		"LinkExpiredTitle":          "The link has expired",
		"LinkExpired":               "The link has expired.",
		"SuccessfulDesc":            "The system will automatically jump to the home page after 3 seconds. If not, please click here.",
		"ImageUploadFailed":         "Image upload failed, please try again later.",
		"ImageTooLarge":             "The file is too large. Please select a smaller image.",
		"EmailSendFailed":           "Email sending failed, please try again later.",
		"TooManyOperations":         "You have tried too many times, please try again later.",
	},
	"category": {
		"DoesNotExist":     "Category does not exist.",
		"CreationFailed":   "Category creation failed, please try again later.",
		"DeleteNonEmpty":   "Cannot delete a non-empty category.",
		"IsNotCategory":    "This is not a category.",
		"CannotBeCopied":   "Category cannot be copied.",
		"CannotBeExported": "Category cannot be exported.",
		"FailedToGet":      "Failed to get category, please try again later.",
		"FailedToDelete":   "Failed to delete category, please try again later.",
	},
	"share": {
		"FailedToGetStatus":           "Failed to get share status, please try again later.",
		"SharingFailed":               "Sharing failed, please try again later.",
		"FailedToDisable":             "Failed to disable sharing, please try again later.",
		"SharedKeyResetFailed":        "Shared key reset failed, please try again later.",
		"SharedKeyVerificationFailed": "Shared key verification failed, please try again later.",
		"SharedKeyError":              "Shared key error.",
		"PublicProjectShare":          "This is a public project that anyone can access.",
	},
	"user": {
		"LoginRequired":                    "Please log in to your account.",
		"OriginalPasswordWrong":            "The original password is wrong.",
		"PasswordUpdateFailed":             "Password modification failed, please try again later.",
		"PasswordResetFailed":              "Password reset failed, please try again later.",
		"PasswordResetSuccessfulTitle":     "Password reset successful",
		"EmailNotChanged":                  "This is the email you are currently using.",
		"EmailHasBeenUsed":                 "This email has already been used.",
		"EmailHasBeenVerified":             "This email has already been verified.",
		"EmailUpdateFailed":                "Email modification failed, please try again later.",
		"EmailUpdateSuccessfulTitle":       "Email modified successfully",
		"EmailUpdateSuccessfulDesc":        "You will use your new email to log in.",
		"EmailVerificationSuccessfulTitle": "Verification successful",
		"EmailVerificationFailedTitle":     "Email verification failed",
		"EmailHasVerifiedTitle":            "Email has been verified",
		"EmailHasVerifiedDesc":             "Your email has been verified.",
		"ResendEmail":                      "Please resend the email.",
		"EmailHasRegistered":               "This email has been registered.",
		"EmailVerificationFailed":          "Email verification failed.",
		"EmailDoesNotExist":                "Email does not exist.",
		"LoginFailed":                      "Login failed, please try again later.",
		"IncorrectEmailOrPassword":         "Incorrect email or password.",
		"RegisterFailed":                   "Registration failed, please try again later.",
		"OauthLoginFailed":                 "Oauth login failed, please try again later.",
		"LoginStatusExpired":               "Login status has expired, please log in again.",
		"InactiveEmail":                    "Inactive email, please activate before trying.",
		"NotSupportOauth":                  "%s is not supported.",
		"OauthConnectFailed":               "Unable to connect to %s, please try again later.",
		"OauthConnectRepeat":               "This %s account has already been linked.",
		"OauthDisconnectFailed":            "Unable to disconnect from %s, please try again later.",
		"FailedToGetList":                  "Failed to get user list, please try again later.",
		"DoesNotExist":                     "User does not exist.",
		"FailedToDelete":                   "Failed to delete user, please try again later.",
	},
	"team": {
		"CreationFailed":             "Team creation failed, please try again later.",
		"FailedToGetList":            "Failed to get team list, please try again later.",
		"FailedToGetCurrentTeam":     "Failed to get current team, please try again later.",
		"CurrentTeamDoesNotExist":    "The current team does not exist.",
		"NoTeam":                     "You currently don't have any teams. Please create a new team to start collaborating immediately.",
		"DoesNotExist":               "Team does not exist.",
		"FailedToSwitch":             "Failed to switch team, please try again later.",
		"InvitationTokenResetFailed": "Invitation token reset failed, please try again later.",
		"InvitationTokenNotFound":    "Invitation token not found.",
		"InvalidInvitationToken":     "Invalid invitation token.",
		"FailedToJoinTeam":           "Failed to join the team, please try again later.",
	},
	"teamMember": {
		"DoesNotExist":              "Team member does not exist.",
		"TeamTransferFailed":        "Team transfer failed, please try again later.",
		"TeamTransferInvalidMember": "You can only transfer the team to an administrator.",
		"FailedToGetList":           "Failed to get team member list, please try again later.",
		"RemoveFailed":              "Team member remove failed, please try again later.",
		"CanNotQuitOwnTeam":         "You cannot quit the team you own.",
		"TeamQuitFailed":            "Quit failed, please try again later.",
		"JoinTeamRepeat":            "You are already a member of the team.",
		"NotTeamMember":             "User is not a member of the team.",
		"Deactivated":               "Team member has been deactivated.",
		"NotInTheTeam":              "This user is not in the team.",
	},
	"project": {
		"DoesNotExist":             "Project does not exist.",
		"FileParseFailed":          "File parsing failed, please try again later.",
		"CreationFailed":           "Project creation failed, please try again later.",
		"FailedToGetList":          "Failed to get project list, please try again later.",
		"FailedToDelete":           "Failed to delete project, please try again later.",
		"TransferFailed":           "Project transfer failed, please try again later.",
		"TransferToErrMember":      "Project can only be transferred to members with write permissions.",
		"TransferToDisabledMember": "Project cannot be handed over to a disabled user.",
		"CanNotQuitOwnProject":     "You cannot quit the project you own.",
		"QuitFailed":               "Quit failed, please try again later.",
		"FailedToUnfollowProject":  "Failed to unfollow project, please try again later.",
		"ExportFailed":             "Project export failed, please try again later.",
		"NotSupportFileType":       "%s is not supported.",
	},
	"projectGroup": {
		"DoesNotExist":          "Project group does not exist.",
		"GroupingFailed":        "Failed to change group, please try again later.",
		"FailedToFollowProject": "Failed to follow project, please try again later.",
		"CreationFailed":        "Project group creation failed, please try again later.",
		"NameHasBeenUsed":       "This group name has already been used.",
		"FailedToGetList":       "Failed to get project group list, please try again later.",
		"FailedToDelete":        "Failed to delete project group, please try again later.",
		"RenameFailed":          "Project group rename failed, please try again later.",
		"SortingFailed":         "Project group sorting failed, please try again later.",
	},
	"projectMember": {
		"DoesNotExist":             "Project member does not exist.",
		"CanNotAddProjectManager":  "A project can only have one manager.",
		"FailedToAddProjectMember": "Failed to add project members, please try again later.",
		"FailedToGetList":          "Failed to get project member list, please try again later.",
		"RemoveFailed":             "Project member remove failed, please try again later.",
		"NotInTheProject":          "The member is not in the project.",
	},
	"projectServer": {
		"CreationFailed":  "Project server URL creation failed, please try again later.",
		"HasBeenUsed":     "This URL has already been used.",
		"FailedToGetList": "Failed to get project server URL list, please try again later.",
		"FailedToDelete":  "Failed to delete project server URL, please try again later.",
		"SortingFailed":   "Project server URL sorting failed, please try again later.",
	},
	"globalParameter": {
		"CreationFailed":  "Global parameter creation failed, please try again later.",
		"HasBeenUsed":     "This parameter has already been used.",
		"FailedToGetList": "Failed to get parameter list, please try again later.",
		"DoesNotExist":    "Global parameter does not exist.",
		"FailedToDelete":  "Failed to delete parameter, please try again later.",
		"FailedToSort":    "Failed to sort global parameter, please try again later.",
	},
	"definitionSchema": {
		"CreationFailed":   "Schema creation failed, please try again later.",
		"GenerationFailed": "Schema generation failed, please try again later.",
		"FailedToGetList":  "Failed to get schema list, please try again later.",
		"FailedToGet":      "Failed to get schema, please try again later.",
		"DoesNotExist":     "Schema does not exist.",
		"FailedToDelete":   "Failed to delete schema, please try again later.",
		"FailedToMove":     "Failed to move schema, please try again later.",
		"CopyFailed":       "Schema copy failed, please try again later.",
	},
	"definitionSchemaHistory": {
		"FailedToGetList": "Failed to get schema history list, please try again later.",
		"FailedToGet":     "Failed to get schema history, please try again later.",
		"DoesNotExist":    "Schema history does not exist.",
		"RestoreFailed":   "Schema restore failed, please try again later.",
		"DiffFailed":      "Schema comparison failed, please try again later.",
	},
	"definitionResponse": {
		"CreationFailed":  "Response creation failed, please try again later.",
		"DoesNotExist":    "Response does not exist.",
		"FailedToGetList": "Failed to get response list, please try again later.",
		"FailedToGet":     "Failed to get response, please try again later.",
		"FailedToDelete":  "Failed to delete response, please try again later.",
		"FailedToMove":    "Failed to move response, please try again later.",
		"CopyFailed":      "Response copy failed, please try again later.",
	},
	"iteration": {
		"DoesNotExist":    "Iteration does not exist.",
		"CreationFailed":  "Iteration creation failed, please try again later.",
		"FailedToGetList": "Failed to get iteration list, please try again later.",
		"FailedToGet":     "Failed to get iteration, please try again later.",
		"FailedToDelete":  "Failed to delete iteration, please try again later.",
	},
	"collection": {
		"FailedToGetList":  "Failed to get API list, please try again later.",
		"FailedToGet":      "Failed to get API, please try again later.",
		"CreationFailed":   "API creation failed, please try again later.",
		"GenerationFailed": "API generation failed, please try again later.",
		"DoesNotExist":     "API does not exist.",
		"FailedToDelete":   "Failed to delete API, please try again later.",
		"FailedToMove":     "Failed to move API, please try again later.",
		"CopyFailed":       "API copy failed, please try again later.",
		"ExportFailed":     "API export failed, please try again later.",
	},
	"collectionHistory": {
		"FailedToGetList": "Failed to get API history list, please try again later.",
		"FailedToGet":     "Failed to get API history, please try again later.",
		"DoesNotExist":    "API history does not exist.",
		"RestoreFailed":   "API restore failed, please try again later.",
		"DiffFailed":      "API comparison failed, please try again later.",
	},
	"testCase": {
		"GenerationFailed":   "Test cases generation failed, please try again later.",
		"FailedToGet":        "Failed to get test case, please try again later.",
		"DoesNotExist":       "Test case does not exist.",
		"RegenerationFailed": "Test case regeneration failed, please try again later.",
		"FailedToDelete":     "Failed to delete test case, please try again later.",
	},
	"mock": {
		"FailedToMock": "Failed to mock, please try again later.",
	},
	"sysConfig": {
		"ServiceUpdateFailed":      "Service setting failed, please try again later.",
		"ServiceBindFailed":        "The IP and port you want to bind are incorrect.",
		"ServiceBindPortSame":      "Your App service port and mock app service port cannot be the same.",
		"OauthUpdateFailed":        "Oauth setting failed, please try again later.",
		"FailedToGetStorageList":   "Failed to get storage config, please try again later.",
		"StorageUpdateFailed":      "Storage setting failed, please try again later.",
		"LocalPathInvalid":         "The local path is invalid.",
		"CloudflareConfigInvalid":  "Invalid cloudflare configuration.",
		"QiniuConfigInvalid":       "Invalid qiniu configuration.",
		"FailedToGetCacheList":     "Failed to get cache config, please try again later.",
		"CacheUpdateFailed":        "Cache setting failed, please try again later.",
		"RedisConfigInvalid":       "Invalid redis configuration.",
		"FailedToGetEmailList":     "Failed to get email config, please try again later.",
		"EmailUpdateFailed":        "Email setting failed, please try again later.",
		"SMTPConfigInvalid":        "Invalid SMTP configuration.",
		"FailedToGetModelList":     "Failed to get model config, please try again later.",
		"ModelUpdateFailed":        "Model setting failed, please try again later.",
		"OpenAIConfigInvalid":      "Invalid OpenAI configuration.",
		"AzureOpenAIConfigInvalid": "Invalid Azure OpenAI configuration.",
	},
	"jsonschema": {
		"JsonSchemaIncorrect": "Parsing failed, please check if the JSON Schema is correct.",
		"FailedToParse":       "Parsing failed, please try again later.",
	},
}
