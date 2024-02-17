// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import "time"

// Defined List of GitHub API
const (
	// GitHubAPIURL is the endpoint for the latest release information of the application.
	GitHubAPIURL      = "https://api.github.com/repos/H0llyW00dzZ/GoGenAI-Terminal-Chat/releases/latest"
	GitHubReleaseFUll = "https://api.github.com/repos/H0llyW00dzZ/GoGenAI-Terminal-Chat/releases/tags/%s"
	// CurrentVersion represents the current version of the application.
	CurrentVersion = "v0.8.10"
)

// Defined constants for the terminal package
const (
	SignalMessage                    = " Received an interrupt, shutting down gracefully..." // fix formatting ^C in linux/unix
	RecoverGopher                    = "%s - %s - %sRecovered from panic:%s %s%v%s"
	StackTracePanic                  = "\n%sStack Trace:\n%s%s"
	StackPossiblyTruncated           = "...stack trace possibly truncated...\n"
	ObjectHighLevelString            = "%s %s"   // Catch High level string
	ObjectHighLevelStringWithSpace   = "%s %s "  // Catch High level string with space
	ObjectHighLevelStringWithNewLine = "%s %s\n" // Catch High level string With NewLine
	ObjectTripleHighLevelString      = "%%%s%%"  // Catch High level triple string
	ObjectHighLevelContextString     = "%s\n%s"  // Catch High level context string
	ObjectHighLevelFMT               = "%s: %s"
	ObjectHighLevelTripleString      = "%s %s %s"
	// TimeFormat is tailored for AI responses, providing a layout conducive to formatting chat transcripts.
	TimeFormat      = "2006/01/02 15:04:05"
	OtherTimeFormat = "January 2, 2006 at 15:04:05"
	StripChars      = "---"
	NewLineChars    = '\n'
	// this animated chars is magic, it used to show the user that the AI is typing just like human would type
	AnimatedChars = "%c"
	// this model is subject to changed in future
	GeminiPro       = "gemini-1.0-pro"
	GeminiProVision = "gemini-pro-vision"
	// this may subject to changed in future for example can customize the delay
	TypingDelay = 60 * time.Millisecond
	// this clearing chat history in secret storage
	ChatHistoryClear = ColorHex95b806 + "All Chat history cleared." + ColorReset
	// reset total token usage
	ResetTotalTokenUsage = ColorHex95b806 + "Total token usage has been reset." + ColorReset
	// clear sys summary messagess
	ChatSysSummaryMessages = ColorHex95b806 + "All System Summary Messages have been cleared." + ColorReset
)

// Defined constants for language
const (
	YouNerd = "🤓 You:"
	AiNerd  = "🤖 AI:"
	// Pass Context to LLM's Google AI
	youNerd                = "🤓"
	aiNerd                 = "🤖"
	sysEmoji               = "⚙️"
	statsEmoji             = "📈"
	TokenEmoji             = "🪙  Token count:"
	StatisticsEmoji        = "📈 Total Token:"
	ShieldEmoji            = "☠️  Safety:"
	ContextPrompt          = "Hello! How can I assist you today?"
	ShutdownMessage        = "Shutting down gracefully..."
	ContextCancel          = "Context canceled, shutting down..." // sending a messages to gopher officer
	ANewVersionIsAvailable = StripChars + "\nA newer version is available: %s\n\n"
	ReleaseName            = "- %s\n\n"
	FullChangeLog          = DoubleAsterisk + "%s" + DoubleAsterisk + "\n"
	DummyMessages          = "Hello, AI! from @H0llyW00dzZ"
	// Better prompt instead of typing manually hahaha
	ApplicationName = "GoGenAI Terminal Chat"
	// Check Version Prompt commands
	YouAreusingLatest = StripChars + "\nThe user invoked the command: " + DoubleAsterisk + "%s" + DoubleAsterisk + "\n" +
		"The current version of " + DoubleAsterisk + "%s" + DoubleAsterisk + " is: " + DoubleAsterisk + "%s" + DoubleAsterisk + ".\n" +
		"This is the latest version available.\n" +
		"Please inform the user that no update is necessary at this time." // Better Response for AI
	ReleaseNotesPrompt = StripChars + "\nThe user invoked the command: " + DoubleAsterisk + "%s" + DoubleAsterisk + "\n" +
		"The current version of the application " + DoubleAsterisk + "%s" + DoubleAsterisk + " is: " + DoubleAsterisk + "%s" + DoubleAsterisk + ".\n" +
		"There is a newer version available: " + DoubleAsterisk + "%s" + DoubleAsterisk + ".\n\n" +
		"Details of the latest release:\n" +
		"- Release Name: " + DoubleAsterisk + "%s" + DoubleAsterisk + "\n" +
		"- Published Date: " + DoubleAsterisk + "%s" + DoubleAsterisk + "\n\n" +
		"Release Notes:\n%s\n" // Better Response for AI
	// Quit Prompt commands
	ContextPromptShutdown = StripChars + "\nThe user has attempted the command: " + DoubleAsterisk + "%s" + DoubleAsterisk + " in " + DoubleAsterisk + "%s" + DoubleAsterisk + ".\n" +
		"AI, please provide an appropriate shutdown message."
	// Help Prompt commands
	HelpCommandPrompt = StripChars + "\n" + DoubleAsterisk + "This a System messages" + DoubleAsterisk + ":" + DoubleAsterisk + "%s" + DoubleAsterisk + "\n\n" +
		"The user attempted an command: " + DoubleAsterisk + "%s" + DoubleAsterisk + "\n" +
		"Can you provide help information for the available commands?\n\n" +
		"List of Available Commands:\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " or " + DoubleAsterisk + "%s" + DoubleAsterisk + ": Quit the application.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " or " + DoubleAsterisk + "%s" + DoubleAsterisk + ": Show this help information.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + ": Check the application version.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + ": Set the safety level - " + DoubleAsterisk + "%s" + DoubleAsterisk + " (low), " + DoubleAsterisk + "%s" + DoubleAsterisk + " (default), " +
		DoubleAsterisk + "%s" + DoubleAsterisk + " (high), " + DoubleAsterisk + "%s" + DoubleAsterisk + " (unspecified), " + DoubleAsterisk + "%s" + DoubleAsterisk + " (none).\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " <text> " + DoubleAsterisk + "%s" + DoubleAsterisk + " <target language>: Translate text to the specified language.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " " + DoubleAsterisk + "%s" + DoubleAsterisk + " <number>: Generate a random string of the specified length.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + ": Summarize a current conversation\n\n" +
		DoubleAsterisk + "Note" + DoubleAsterisk + ": When you summarize a current conversation, it will be displayed at the top of the chat history.\n\n" +
		DoubleAsterisk + "%s %s" + DoubleAsterisk + " " + DoubleAsterisk + "%s" + DoubleAsterisk + ": Show the chat history.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " " + DoubleAsterisk + "%s" + DoubleAsterisk + ": Show the chat statistic.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " " + DoubleAsterisk + "%s" + DoubleAsterisk + ": Clear all system summary messages from the chat history.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " " + DoubleAsterisk + "%s" + DoubleAsterisk + ": Clear all chat history and reset the total token usage count if enabled.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " <" + DoubleAsterisk + "model-name" + DoubleAsterisk + ">: Check the details of a specific AI model.\n" +
		DoubleAsterisk + "%s" + DoubleAsterisk + " " + DoubleAsterisk + "%s" + DoubleAsterisk + " <" + DoubleAsterisk + "path/file/data.txt" + DoubleAsterisk + "> or <" +
		DoubleAsterisk + "data.txt" + DoubleAsterisk + ">: Counts a token from the specified file.\n\n" +
		DoubleAsterisk + "Note" + DoubleAsterisk + ": The token count file feature supports multiple files simultaneously with the following extensions: " +
		dotMD + dotStringComma + dotTxt + dotStringComma + dotPng + dotStringComma +
		dotJpg + dotStringComma + dotJpeg + dotStringComma + dotWebp + dotStringComma +
		dotHeic + dotStringComma + dotHeif + ".\n\n" +
		DoubleAsterisk + "Additional Note" + DoubleAsterisk + ": There are no additional commands or HTML Markdown available " +
		"because this is a terminal application and is limited.\n"
	// TranslateCommandPrompt commands
	AITranslateCommandPrompt = DoubleAsterisk + "This a System messages" + DoubleAsterisk + ":" + DoubleAsterisk + "%s" + DoubleAsterisk + "\n\n" +
		"The user attempted an command: " + DoubleAsterisk + "%s" + DoubleAsterisk + "\n" +
		"Can you translate requested by user?\n" +
		"Text:\n" + DoubleAsterisk + "%s" + DoubleAsterisk + "\n" +
		"Translate To:\n " + DoubleAsterisk + "%s" + DoubleAsterisk
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
	ChatCommands       = ":chat"
	SummarizeCommands  = ":summarize"
	ClearCommand       = ":clear"
	StatsCommand       = ":stats"
	TokenCountCommands = ":tokencount"
	FileCommands       = ":file"
	CheckModelCommands = ":checkmodel"
	PingCommand        = ":ping" // Currently marked as TODO
	PrefixChar         = ":"
	// List args
	ChatHistoryArgs = "history"
)

// Defined List error message
const (
	ErrorGettingShutdownMessage                     = "Error getting shutdown message from AI: %v"
	ErrorHandlingCommand                            = "Error handling command: %v"
	ErrorCountingTokens                             = "Error counting tokens: %v\n"
	ErrorSendingMessage                             = "Error sending message to AI: %v"
	ErrorReadingUserInput                           = "Error reading user input: %v"
	ErrorFailedToFetchReleaseInfo                   = "Failed to fetch the latest release %s info: %v"
	ErrorReceivedNon200StatusCode                   = "[Retry Policy] [Github] [Check Version] Received non-200 status code: %v Skip Retrying" // Github non 500 lmao
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
	ErrorWhileTypingCommandArgs                     = "Invalid %s Command Arguments: %v"
	ErrorPingFailed                                 = "Ping failed: %v"
	ErrorUnrecognizedCommand                        = "Unrecognized command: %s"
	ErrorUnrecognizedSubCommand                     = "Unrecognized %s command sub commands/args : %s"
	ErrorLowLevelCommand                            = "command cannot be empty"
	ErrorUnknown                                    = "An error occurred: %v"
	ErrorUnknownSafetyLevel                         = "Unknown safety level: %s"
	ErrorInvalidAPIKey                              = "Invalid API key: %v"
	ErrorFailedToStartSession                       = "Failed To Start Session: %v"
	ErrorLowLevelNoResponse                         = "no response from AI service"
	ErrorLowLevelMaximumRetries                     = "[Retry Policy] maximum retries reached without success - %v" // low level
	ErrorLowLevelFailedToCountTokensAfterRetries    = "failed to count tokens after retries"                        // low level
	ErrorNonretryableerror                          = "[Retry Policy] Failed to Retrying after %d (SKIPPED) retries due to a non-retryable error: %v"
	ErrorFailedToSendHelpMessage                    = "Failed to send help message: %v"
	ErrorFailedToSendHelpMessagesAfterRetries       = "Failed to send help message after retries" // low level
	ErrorFailedToSendShutdownMessage                = "Failed to send shutdown message: %v"
	ErrorFailedToSendVersionCheckMessage            = "Failed to send version check message: %v"
	ErrorFailedToSendVersionCheckMessageAfterReties = "Failed to send version check message after retries" // low level
	ErrorFailedToSendTranslationMessage             = "Failed to send translation message: %v"
	ErrorFailedToSendTranslationMessageAfterRetries = "Failed to send translation message after retries" // low level
	ErrorFailedToApplyModelConfiguration            = "failed to apply model configuration"              // low level
	ErrorMaxOutputTokenMustbe                       = "maxOutputTokens must be %d or higher, got %d"     // low level
	ErrorFailedToSendSummarizeMessage               = "Failed To Send Summarize Message: %v"
	ErrorFailedToSendSummarizeMessageAfterRetries   = "failed to send summarize message after retries" // low level
	ErrorFailedToReadFile                           = "Failed to read the file at %s: %v"
	ErrorFailedToCountTokens                        = "Failed to count tokens in the file at %s: %v"
	ErrorUnrecognizedSubcommandForTokenCount        = "Unrecognized subcommand for token count: %s"
	ErrorInvalidFileExtension                       = "Invalid file extension: %v"
	ErrorFileTypeNotSupported                       = "file type not supported: only %s files are allowed." // Low level error
	ErrorFailedToSendCommandToAI                    = "Failed to send command to AI: %v"
	ErrorVariableImageFileTypeNotSupported          = "image file type not supported: only %s files are allowed." // Low level error
	ErrorNoInputProvideForTokenCounting             = "no input provided for token counting"                      // low level error
	ErrorGopherEncounteredAnError                   = "Goroutine %d encountered an error: %w"
	ErrorFailedToRetriveModelInfo                   = "Failed to retrieve model info: %v"

	// List Error not because of this go codes, it literally google apis issue
	// that so bad can't handle this a powerful terminal
	Error500GoogleAPI    = "googleapi: Error 500:"
	ErrorGoogleInternal  = "Google Internal Error: %s"
	ErrorGenAiReceiveNil = "received a nil option function" // low level
	ErrorGenAI           = "GenAI Error: %v"
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
	SingleMinusSign         = "-"
	StringNewLine           = "\n"
	BinaryAnsiChar          = '\x1b'
	BinaryLeftSquareBracket = '['
	BinaryAnsiSquenseChar   = 'm'
	BinaryAnsiSquenseString = "m"
	BinaryRegexAnsi         = `\x1b\[[0-9;]*m`
	CodeBlockRegex          = "```\\w+"
	SanitizeTextAIResponse  = "\n---\n"
	// This regex pattern will match:
	//
	// 1. An asterisk followed by a space (e.g., "* ")
	//
	// 2. Italic text surrounded by single asterisks (e.g., "*italic*")
	//
	// 3. A leading asterisk followed by space and italic text (e.g., "* *italic*")
	//
	// Ref: https://go.dev/play/p/3APEsqfgUV-
	//
	// It uses a non-greedy match for the italic text and optional groups.
	ItalicTextRegex = `\*(\s\*\S.*?\S\*|\s|\S.*?\S\*)`
	// TODO
	StandaloneAsteriskAnsiRegexPattern = `(?m)(^|\s)\*(\s|$)`
)

// Defined List of Environment variables
const (
	DebugMode   = "DEBUG_MODE"
	DEBUGPREFIX = "🔎 DEBUG:"
	// Note: Currently only executing CMD,RetryPolicy, will add more later
	DEBUGEXECUTINGCMD = "Executing " +
		// Better Readability use Custom HEX color
		ColorHex95b806 + "%s" + ColorReset +
		" command with parts: " +
		// Better Readability use Custom HEX color
		ColorHex95b806 + "%#v" + ColorReset
	DEBUGRETRYPOLICY   = "Retry Policy Attempt %d: error occurred - %v"
	ShowPromptFeedBack = "SHOW_PROMPT_FEEDBACK"
	PROMPTFEEDBACK     = "Rating for category " + ColorHex95b806 + "%s" + ColorReset + ": " +
		ColorHex95b806 + "%s" + ColorReset
	ShowTokenCount  = "SHOW_TOKEN_COUNT"
	TokenCount      = ColorHex95b806 + "%d" + ColorReset + " tokens\n"
	TotalTokenCount = "usage of this Session " + ColorHex95b806 + "%d" + ColorReset + " tokens"
	// Note: This is separate from the main package and is used for the token counter. The token counter is external and not a part of the Gemini session.
	APIKey = "API_KEY"
)

// Defined Prefix System
const (
	// Note: This is a prefix for the system
	SYSTEMPREFIX     = "⚙️ SYSTEM:"
	SystemSafety     = "Safety level set to " + ColorHex95b806 + "%s" + ColorReset + "."
	Low              = "low"
	Default          = "default"
	High             = "high"
	Unspecified      = "unspecified"
	None             = "none"
	MonitoringSignal = "Received signal: %v.\n"
	ShowChatHistory  = "Chat History:\n\n%s"
	SummarizePrompt  = StripChars + "\nIn 200 words or less, provide a brief summary of the ongoing discussion.\n" +
		"This summary will serve as a prompt for contextual reference in future interactions:\n\n"

	ListChatStats = statsEmoji + " List of Chat Statistics for This Session:\n\n" +
		youNerd + " User messages: " + ColorHex95b806 + BoldText + "%d" + ResetBoldText + ColorReset + "\n" +
		aiNerd + " AI messages: " + ColorHex95b806 + BoldText + "%d" + ResetBoldText + ColorReset + "\n" +
		sysEmoji + " System messages: " + ColorHex95b806 + BoldText + "%d" + ResetBoldText + ColorReset
	InfoTokenCountFile = "The file " + ColorHex95b806 + BoldText + "%s" + ResetBoldText + ColorReset +
		" contains " + ColorHex95b806 + BoldText + "%d" + ResetBoldText + ColorReset + " tokens."
	RetryingStupid500Error = "[Retry Policy] Retrying (" + ColorRed + "last error: %v" + ColorReset + ")" +
		" attempt number " + ColorHex95b806 + BoldText + "%d" + ResetBoldText + ColorReset
	// ModelFormat defines a template for displaying model information with color and bold formatting for placeholders.
	ModelFormat = "Model Name: " + ColorHex95b806 + BoldText + "%s" + ResetBoldText + ColorReset + "\n" +
		"Model ID: " + ColorHex95b806 + BoldText + "%s" + ResetBoldText + ColorReset + "\n" +
		"Version: " + ColorHex95b806 + BoldText + "%s" + ResetBoldText + ColorReset + "\n" +
		"Description: " + ColorHex95b806 + BoldText + "%s" + ResetBoldText + ColorReset + "\n" +
		"Supported Generation Methods: " + ColorHex95b806 + BoldText + "%s" + ResetBoldText + ColorReset + "\n" +
		"Input Token Limit: " + ColorHex95b806 + BoldText + "%d" + ResetBoldText + ColorReset + "\n" +
		"Output Token Limit: " + ColorHex95b806 + BoldText + "%d" + ResetBoldText + ColorReset
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
	P = 'P'
	D = 'D'
	// ASCII slant font
	_G   = "   ______      ______           ___    ____  "
	_O   = "  / ____/___  / ____/__  ____  /   |  /  _/  "
	_GEN = " / / __/ __ \\/ / __/ _ \\/ __ \\/ /| |  / /    "
	A    = "/ /_/ / /_/ / /_/ /  __/ / / / ___ |_/ /     "
	I    = "\\____/\\____/\\____/\\___/_/ /_/_/  |_/___/     "
	// Blank Art
	Blank   = "                                      "
	eMpty   = ""
	slantp1 = "    ____              _     "
	slanta2 = "   / __ \\____ _____  (_)____"
	slantn3 = "  / /_/ / __ `/ __ \\/ / ___/"
	slanti4 = " / ____/ /_/ / / / / / /__  "
	slantc5 = "/_/    \\__,_/_/ /_/_/\\___/  "
	slantA1 = "    ____       __            __           __"
	slantA2 = "   / __ \\___  / /____  _____/ /____  ____/ /"
	slantA3 = "  / / / / _ \\/ __/ _ \\/ ___/ __/ _ \\/ __  / "
	slantA4 = " / /_/ /  __/ /_/  __/ /__/ /_/  __/ /_/ /  "
	slantA5 = "/_____/\\___/\\__/\\___/\\___/\\__/\\___/\\__,_/   "
)

// Text
const (
	CurrentVersioN = "Current Version: " + ColorHex95b806 + CurrentVersion + ColorReset
	// Acknowledgment of the original author is appreciated as this project is developed in an open-source environment.
	Copyright = "Copyright (©️) 2024 @H0llyW00dzZ All rights reserved."
	TIP       = SingleAsterisk + ColorHex95b806 + " Use the commands " + ColorReset +
		BoldText + ColorYellow + ShortHelpCommand + ColorYellow +
		BoldText + ColorHex95b806 + " or " + ColorReset + BoldText + ColorYellow + HelpCommand + ColorReset +
		BoldText + ColorHex95b806 + " to display a list of available commands." + ColorReset
)

// Context RAM's labyrinth
const (
	ContextUserInvokeTranslateCommands = "Translating to %s: %s"
	SummaryPrefix                      = aiNerd + " 📝 📌 Summary of this discussion:\n\n"
)

// List RestfulAPI Error
const (
	Code500 = "500" // indicate that server so bad hahaha
)

// dotFiles
const (
	// a better way instead of stupid hardcoding
	// Reason: In compiled languages, "" denotes a string constant and '' denotes a rune, unlike in interpreted languages that sometimes is confusing.
	// Also Using constants improves readability and maintainability over stupid hardcoding values.
	dotMD          = ".md"
	dotTxt         = ".txt"
	dotPng         = ".png"
	dotJpg         = ".jpg"
	dotJpeg        = ".jpeg"
	dotWebp        = ".webp"
	dotHeic        = ".heic"
	dotHeif        = ".heif"
	dotString      = "."
	dotStringComma = ", "
	oRString       = " or "
)

// ANSI color codes
const (
	// Note: By replacing the ANSI escape sequence from "\033" to "\x1b", might can avoid a rare bug that sometimes occurs on different machines,
	// although the original code works fine on mine (Author: @H0llyW00dzZ).
	ColorRed    = "\x1b[31m"
	ColorGreen  = "\x1b[32m"
	ColorYellow = "\x1b[33m"
	ColorBlue   = "\x1b[34m"
	ColorPurple = "\x1b[35m"
	ColorCyan   = "\x1b[36m"
	// ColorHex95b806 represents the color #95b806 using an ANSI escape sequence for 24-bit color.
	ColorHex95b806 = "\x1b[38;2;149;184;6m"
	// ColorCyan24Bit represents the color #11F0F7 using an ANSI escape sequence for 24-bit color.
	ColorCyan24Bit   = "\x1b[38;2;17;240;247m"
	ColorPurple24Bit = "\x1b[38;2;255;0;255m"
	ColorReset       = "\x1b[0m"
)

// ANSI Text Formatting.
const (
	// bold text.
	BoldText = "\x1b[1m"
	// reset bold text formatting.
	ResetBoldText = "\x1b[22m"
	// italic text
	ItalicText = "\x1B[3m"
	// reset italic text formatting.
	ResetItalicText = "\x1B[23m"
)

const (
	// UserMessage indicates a message that originates from a human user.
	UserMessage MessageType = iota // magic
	// AIMessage indicates a message that originates from an AI or bot.
	AIMessage
	// SystemMessage indicates a message that provides system-level information.
	SystemMessage
)

// mime formatting
const (
	FormatJPEG = "jpeg"
	FormatPNG  = "png"
	FormatHEIC = "heic"
	FormatHEIF = "heif"
	FormatWEBP = "webp"
)

// model configuration
const (
	MinOutputTokens int32 = 20 // Define the minimum number of tokens as a constant
)
