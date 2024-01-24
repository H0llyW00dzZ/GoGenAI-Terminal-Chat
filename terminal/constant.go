// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import "time"

// Defined List of GitHub API
const (
	// GitHubAPIURL is the endpoint for the latest release information of the application.
	GitHubAPIURL      = "https://api.github.com/repos/H0llyW00dzZ/GoGenAI-Terminal-Chat/releases/latest"
	GitHubReleaseFUll = "https://api.github.com/repos/H0llyW00dzZ/GoGenAI-Terminal-Chat/releases/tags/%s"
	// CurrentVersion represents the current version of the application.
	CurrentVersion = "v0.5.1"
)

// Defined constants for the terminal package
const (
	SignalMessage                    = " Received an interrupt, shutting down gracefully..." // fix formatting ^C in linux/unix
	RecoverGopher                    = "%sRecovered from panic:%s %s%v%s"
	ObjectHighLevelString            = "%s %s"   // Catch High level string
	ObjectHighLevelStringWithNewLine = "%s %s\n" // Catch High level string With NewLine
	ObjectTripleHighLevelString      = "%%%s%%"  // Catch High level triple string
	ObjectHighLevelContextString     = "%s\n%s"  // Catch High level context string
	// TimeFormat is tailored for AI responses, providing a layout conducive to formatting chat transcripts.
	TimeFormat      = "2006/01/02 15:04:05"
	OtherTimeFormat = "January 2, 2006 at 15:04:05"
	StripChars      = "---"
	NewLineChars    = '\n'
	// this animated chars is magic, it used to show the user that the AI is typing just like human would type
	AnimatedChars = "%c"
	// this model is subject to changed in future
	ModelAi = "gemini-pro"
	// this may subject to changed in future for example can customize the delay
	TypingDelay = 60 * time.Millisecond
	// this clearing chat history in secret storage
	ChatHistoryClear = ColorHex95b806 + "All Chat history cleared." + ColorReset
)

// Defined constants for language
const (
	YouNerd = "🤓 You:"
	AiNerd  = "🤖 AI:"
	// Pass Context to LLM's Google AI
	youNerd                = "🤓"
	aiNerd                 = "🤖"
	TokenEmoji             = "🪙  Token count:"
	StatisticsEmoji        = "📈 Total Token:"
	ShieldEmoji            = "☠️  Safety:"
	ContextPrompt          = "Hello! How can I assist you today?"
	ShutdownMessage        = "Shutting down gracefully..."
	ContextCancel          = "Context canceled, shutting down..." // sending a messages to gopher officer
	ANewVersionIsAvailable = "A newer version is available: %s\n\n"
	ReleaseName            = "- %s\n\n"
	FullChangeLog          = "**%s**\n"
	DummyMessages          = "Hello, AI! from @H0llyW00dzZ"
	// Better prompt instead of typing manually hahaha
	//
	// Note: These prompts are not persisted in the chat history retrieved by the ChatHistory.GetHistory() method.
	// Therefore, if you continue interacting with the AI after using these command prompts,
	// the conversation will resume from the point prior to the invocation of these commands.
	ApplicationName = "GoGenAI Terminal Chat"
	// Check Version Prompt commands
	YouAreusingLatest = "a User attempted a command: **%s**\n" +
		"The user is using Version **%s** of **%s**\n" +
		"This is the latest version.\n" +
		"Tell the user, No need to update." // Better Response for AI
	ReleaseNotesPrompt = "a user attempted a command: **%s**\n" +
		"The user is using Version **%s** of **%s**\n" +
		"A newer version is available: **%s**\n" +
		"Can you tell\n" +
		"Release Name: **%s**\n" +
		"Published Date: **%s**\n\n%s" // Better Response for AI
	// Quit Prompt commands
	ContextPromptShutdown = "a user attempted an command: **%s** of **%s**\n" +
		"Please provide a shutdown message as you are AI."
	// Help Prompt commands
	HelpCommandPrompt = "**This a System messages**:**%s**\n\n" +
		"The user attempted an command: **%s**\n" +
		"Can you provide help information for the available commands?\n" +
		// Better Response for AI instead of "Hard Coded" hahaha
		"List Command Available:\n**%s** or **%s**\n**%s** or **%s**\n" +
		"**%s** or **%s**\n**%s** - **%s**, **%s**, **%s**\n" +
		"**%s** <text> **%s** <targetlanguage>\n" +
		"**%s** **%s** <number>\n\n**%s %s**\n\n**%s %s**\n\n" +
		"**Additional Note**: There are no **additional commands** or **HTML Markdown** available" +
		" because it is in a terminal and is limited.\n"
	// TranslateCommandPrompt commands
	AITranslateCommandPrompt = "**This a System messages**:**%s**\n\n" +
		"The user attempted an command: **%s**\n" +
		"Can you translate requested by user?\n" +
		"Text:\n**%s**\n" +
		"Translate To:\n **%s**"
)

// Defined constants for commands
//
// Note: will add more in future based on the need,
// for example, to change the model, or to change the delay, another thing is syncing ai with goroutine (known as gopher)
const (
	QuitCommand        = ":quit"
	ShortQuitCommand   = ":q" // Short quit command
	VersionCommand     = ":checkversion"
	HelpCommand        = ":help"
	ShortHelpCommand   = ":h" // Short help command
	SafetyCommand      = ":safety"
	AITranslateCommand = ":aitranslate"
	LangArgs           = ":lang"
	CryptoRandCommand  = ":cryptorand"
	LengthArgs         = ":length"
	ShowCommands       = ":show"
	ChatArgs           = ":chat"
	PingCommand        = ":ping" // Currently marked as TODO
	ClearCommand       = ":clear"
	PrefixChar         = ":"
	// List args
	ChatHistoryArgs = "chat history"
)

// Defined List error message
const (
	ErrorGettingShutdownMessage                     = "Error getting shutdown message from AI: %v"
	ErrorHandlingCommand                            = "Error handling command: %v"
	ErrorCountingTokens                             = "Error counting tokens: %v\n"
	ErrorSendingMessage                             = "Error sending message to AI: %v"
	ErrorReadingUserInput                           = "Error reading user input: %v"
	ErrorFailedToFetchReleaseInfo                   = "Failed to fetch the latest release info: %v"
	ErrorReceivedNon200StatusCode                   = "[Github] [Check Version] Received non-200 status code: %v Skip Retrying" // Github non 500 lmao
	ErrorFailedToReadTheResponseBody                = "Failed to read the response body: %v"
	ErrorFaileduUnmarshalTheReleaseData             = "Failed to unmarshal the release data: %v"
	ErrorFailedTagToFetchReleaseInfo                = "Failed to fetch release info for tag '%s': %v"
	ErrorFailedTagUnmarshalTheReleaseData           = "Failed to unmarshal release data for tag '%s': %v"
	ErrorFailedTosendmessagesToAI                   = "Failed to send messages to AI: %v"
	ErrorFailedToCreateNewAiClient                  = "Failed to create new AI client: %v"
	ErrorFailedToStartAIChatSessionAttempt          = "Failed to start AI chat session, attempt %d/%d"
	ErrorFailedtoStartAiChatSessionAfter            = "Failed to start AI chat session after %d attempts"
	ErrorChatSessionisnill                          = "chat session is nil"
	ErrorFailedtoStartAiChatSession                 = "failed to start AI chat session"
	ErrorFailedToRenewSession                       = "Failed to renew session: %v"
	ErrorAiChatSessionStillNill                     = "AI chat session is still nil after renewal attempt"
	ErrorLowLevelFailedtoStartAiChatSession         = "failed to start a new AI chat session: %w"
	ErrorUserAttemptUnrecognizedCommandPrompt       = "**From System**:**%s**\n\nThe user attempted an unrecognized command: **%s**" // Better Response for AI
	ErrorFailedtoSendUnrecognizedCommandToAI        = "Failed to send unrecognized command to AI: %v"
	HumanErrorWhileTypingCommandArgs                = "Invalid Command Arguments: %v"
	ErrorPingFailed                                 = "Ping failed: %v"
	ErrorUnrecognizedCommand                        = "Unrecognized command: %s"
	ErrorLowLevelCommand                            = "command cannot be empty"
	ErrorUnknown                                    = "An error occurred: %v"
	ErrorUnknownSafetyLevel                         = "Unknown safety level: %s"
	ErrorInvalidApiKey                              = "Invalid API key: %v"
	ErrorLowLevelNoResponse                         = "no response from AI service"
	ErrorLowLevelMaximumRetries                     = "maximum retries reached without success" // low level
	ErrorLowLevelFailedToCountTokensAfterRetries    = "failed to count tokens after retries"    // low level
	ErrorNonretryableerror                          = "Failed to send messages after %d retries due to a non-retryable error: %v"
	ErrorFailedToSendHelpMessage                    = "Failed to send help message: %v"
	ErrorFailedToSendHelpMessagesAfterRetries       = "Failed to send help message after retries" // low level
	ErrorFailedToSendShutdownMessage                = "Failed to send shutdown message: %v"
	ErrorFailedToSendVersionCheckMessage            = "Failed to send version check message: %v"
	ErrorFailedToSendVersionCheckMessageAfterReties = "Failed to send version check message after retries" // low level
	ErrorFailedToSendTranslationMessage             = "Failed to send translation message: %v"
	ErrorFailedToSendTranslationMessageAfterRetries = "Failed to send translation message after retries" // low level
	// List Error not because of this go codes, it literally google apis issue
	// that so bad can't handle this a powerful terminal
	Error500GoogleApi   = "googleapi: Error 500:"
	ErrorGoogleInternal = "Google Internal Error: %s"
	// List Error Figlet include high and low level error
	ErrorStyleIsEmpty             = "style is empty"                   // low level
	ErrorCharacterNotFoundinStyle = "character %q not found in style"  // low level
	ErrorToASCIIArtbuildOutput    = "ToASCIIArt buildOutput error: %v" // High Level
	ErrorToASCIIArtcheckstyle     = "ToASCIIArt checkStyle error: %v"  // High Level
	// List Error Tools
	ErrorInvalidLengthArgs            = "Invalid length argument: %v"          // high level
	errorinvalidlengthArgs            = "invalid length argument: %v"          // low level
	ErrorFailedtoGenerateRandomString = "Failed to generate random string: %v" // high level
	errorfailedtogeneraterandomstring = "failed to generate random string: %v" // low level
	// List Other Error not because of this go codes
	// ErrorOtherAPI represents an error received from an external API server.
	// It indicates non-client-related issues, such as server-side errors (e.g., HTTP 500 errors) indicate that so fucking bad hahaha.
	ErrorOtherAPI = "Error: %s API server error: %v"
)

// Defined List of characters
const (
	SingleAsterisk          = "*"
	DoubleAsterisk          = "**"
	SingleBacktick          = "`"
	TripleBacktick          = "```"
	SingleUnderscore        = "_"
	StringNewLine           = "\n"
	BinaryAnsiChar          = '\x1b'
	BinaryLeftSquareBracket = '['
	BinaryAnsiSquenseChar   = 'm'
	BinaryAnsiSquenseString = "m"
	BinaryRegexAnsi         = `\x1b\[[0-9;]*m`
	CodeBlockRegex          = "```\\w+"
)

// Defined List of Environment variables
const (
	DEBUG_MODE  = "DEBUG_MODE"
	DEBUGPREFIX = "🔎 DEBUG:"
	// Note: Currently only executing CMD,RetryPolicy, will add more later
	DEBUGEXECUTINGCMD = "Executing " +
		// Better Readability use Custom HEX color
		ColorHex95b806 + "%s" + ColorReset +
		" command with parts: " +
		// Better Readability use Custom HEX color
		ColorHex95b806 + "%#v" + ColorReset
	DEBUGRETRYPOLICY     = "Retry Policy Attempt %d: error occurred - %v"
	SHOW_PROMPT_FEEDBACK = "SHOW_PROMPT_FEEDBACK"
	PROMPTFEEDBACK       = "Rating for category " + ColorHex95b806 + "%s" + ColorReset + ": " +
		ColorHex95b806 + "%s" + ColorReset
	SHOW_TOKEN_COUNT = "SHOW_TOKEN_COUNT"
	TokenCount       = ColorHex95b806 + "%d" + ColorReset + " tokens\n"
	TotalTokenCount  = "usage of this Session " + ColorHex95b806 + "%d" + ColorReset + " tokens"
	// Note: This is separate from the main package and is used for the token counter. The token counter is external and not a part of the Gemini session.
	API_KEY = "API_KEY"
)

// Defined Prefix System
const (
	// Note: This is a prefix for the system
	SYSTEMPREFIX     = "⚙️  SYSTEM:"
	SystemSafety     = "Safety level set to " + ColorHex95b806 + "%s" + ColorReset + "."
	Low              = "low"
	Default          = "default"
	High             = "high"
	MonitoringSignal = "Received signal: %v.\n"
	ShowChatHistory  = "Chat History:\n\n%s"
)

// Defined Tools
const (
	CryptoRandLength    = "Length: %s"
	CryptoRandStringRes = "Random String: %s"
	CryptoRandRes       = "Length: %s\n\nRandom String: %s"
)

// ASCII Art
const (
	// NOTE: ' is rune not a string
	G = 'G'
	V = 'V'
	N = 'N'
	// ASCII slant font
	_G   = "   ______      ______           ___    ____  "
	_O   = "  / ____/___  / ____/__  ____  /   |  /  _/  "
	_GEN = " / / __/ __ \\/ / __/ _ \\/ __ \\/ /| |  / /    "
	A_   = "/ /_/ / /_/ / /_/ /  __/ / / / ___ |_/ /     "
	I_   = "\\____/\\____/\\____/\\___/_/ /_/_/  |_/___/     "
	// Blank Art
	BLANK_ = "                                      "
	eMpty  = ""
)

// Text
const (
	Current_Version = "Current Version: " + ColorHex95b806 + CurrentVersion + ColorReset
	// Acknowledgment of the original author is appreciated as this project is developed in an open-source environment.
	Copyright = "Copyright (©️) 2024 @H0llyW00dzZ All rights reserved."
	TIP       = "* " + ColorHex95b806 + "Use the commands " + ColorReset +
		BoldText + ColorYellow + ShortHelpCommand + ColorYellow +
		BoldText + ColorHex95b806 + " or " + ColorReset + BoldText + ColorYellow + HelpCommand + ColorReset +
		BoldText + ColorHex95b806 + " to display a list of available commands." + ColorReset
)

// Context RAM's labyrinth
const (
	ContextUserInvokeTranslateCommands = "Translating to %s: %s"
)

// List RestfulAPI Error
const (
	Code500 = "500" // indicate that server so bad hahaha
)
