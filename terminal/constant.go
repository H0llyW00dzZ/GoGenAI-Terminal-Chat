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
	CurrentVersion = "v0.3.8"
)

// Defined constants for the terminal package
const (
	SignalMessage               = "\nReceived an interrupt, shutting down gracefully..."
	RecoverGopher               = "%sRecovered from panic:%s %s%v%s"
	ObjectHighLevelString       = "%s %s"  // Catch High level string
	ObjectTripleHighLevelString = "%%%s%%" // Catch High level triple string
	// TimeFormat is tailored for AI responses, providing a layout conducive to formatting chat transcripts.
	TimeFormat   = "2006/01/02 15:04:05"
	StripChars   = "---"
	NewLineChars = '\n'
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
	YouNerd = "🤓 You: "
	AiNerd  = "🤖 AI: "
	// Pass Context to LLM's Google AI
	youNerd                = "🤓"
	aiNerd                 = "🤖"
	ContextPrompt          = "Hello! How can I assist you today?"
	ShutdownMessage        = "Shutting down gracefully..."
	UnknownCommand         = "Unknown command."
	ContextCancel          = "Context canceled, shutting down..." // sending a messages to gopher officer
	ANewVersionIsAvailable = "A newer version is available: %s\n\n"
	ReleaseName            = "- %s\n\n"
	FullChangeLog          = "**%s**\n"
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
		"%s" // Better Response for AI
	// Quit Prompt commands
	ContextPromptShutdown = "a user attempted an command: **%s** of **%s**\n" +
		"Please provide a shutdown message as you are AI."
	// Help Prompt commands
	HelpCommandPrompt = "**This a System messages**:**%s**\n\n" +
		"The user attempted an command: **%s**\n" +
		"Can you provide help information for the available commands?\n" +
		// Better Response for AI instead of "Hard Coded" hahaha
		"List Command Available:\n**%s** or **%s**\n**%s** or **%s**\n" +
		"**%s** or **%s**\n**%s** - **%s**, **%s**, **%s**" +
		"**%s** <text> **%s** <targetlanguage>\n\n**%s %s**"
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
	PingCommand        = ":ping" // Currently marked as TODO
	ClearCommand       = ":clear"
	PrefixChar         = ":"
	// List args
	ClearChatHistoryArgs = "chat history"
)

// Defined List error message
const (
	ErrorGettingShutdownMessage               = "Error getting shutdown message from AI: %v"
	ErrorHandlingCommand                      = "Error handling command: %v"
	ErrorCountingTokens                       = "Error counting tokens: %v"
	ErrorSendingMessage                       = "Error sending message to AI: %v"
	ErrorReadingUserInput                     = "Error reading user input: %v"
	ErrorFailedToFetchReleaseInfo             = "Failed to fetch the latest release info: %v"
	ErrorReceivedNon200StatusCode             = "Received non-200 status code: %v"
	ErrorFailedToReadTheResponseBody          = "Failed to read the response body: %v"
	ErrorFaileduUnmarshalTheReleaseData       = "Failed to unmarshal the release data: %v"
	ErrorFailedTagToFetchReleaseInfo          = "Failed to fetch release info for tag '%s': %v"
	ErrorFailedTagUnmarshalTheReleaseData     = "Failed to unmarshal release data for tag '%s': %v"
	ErrorFailedTosendmessagesToAI             = "Failed to send messages to AI: %v"
	ErrorFailedToCreateNewAiClient            = "Failed to create new AI client: %v"
	ErrorFailedToStartAIChatSessionAttempt    = "Failed to start AI chat session, attempt %d/%d"
	ErrorFailedtoStartAiChatSessionAfter      = "Failed to start AI chat session after %d attempts"
	ErrorChatSessionisnill                    = "chat session is nil"
	ErrorFailedtoStartAiChatSession           = "failed to start AI chat session"
	ErrorFailedToRenewSession                 = "Failed to renew session: %v"
	ErrorAiChatSessionStillNill               = "AI chat session is still nil after renewal attempt"
	ErrorLowLevelFailedtoStartAiChatSession   = "failed to start a new AI chat session: %w"
	ErrorUserAttemptUnrecognizedCommandPrompt = "**From System**:**%s**\n\nThe user attempted an unrecognized command: **%s**" // Better Response for AI
	ErrorFailedtoSendUnrecognizedCommandToAI  = "Failed to send unrecognized command to AI: %v"
	HumanErrorWhileTypingCommandArgs          = "Invalid Command Arguments"
	ErrorPingFailed                           = "Ping failed: %v"
)

// Defined List of characters
const (
	SingleAsterisk          = "*"
	DoubleAsterisk          = "**"
	SingleBacktick          = "`"
	TripleBacktick          = "```"
	StringNewLine           = "\n"
	BinaryAnsiChar          = '\x1b'
	BinaryLeftSquareBracket = '['
	BinaryAnsiSquenseChar   = 'm'
	BinaryAnsiSquenseString = "m"
	BinaryRegexAnsi         = `\x1b\[[0-9;]*m`
)

// Defined List of Environment variables
const (
	DEBUG_MODE  = "DEBUG_MODE"
	DEBUGPREFIX = "🔎 DEBUG:"
	// Note: Currently only executing CMD, will add more later
	DEBUGEXECUTINGCMD = "Executing " +
		// Better Readability use Custom HEX color
		ColorHex95b806 + "%s" + ColorReset +
		" command with parts: " +
		// Better Readability use Custom HEX color
		ColorHex95b806 + "%#v" + ColorReset
	SHOW_PROMPT_FEEDBACK = "SHOW_PROMPT_FEEDBACK"
	PROMPTFEEDBACK       = "Safety Rating for category %s: %s\n"
)

// Defined Prefix System
const (
	// Note: This is a prefix for the system
	SYSTEMPREFIX = "⚙️  SYSTEM: "
	SystemSafety = "Safety level set to " + ColorHex95b806 + "%s" + ColorReset + "."
	Low          = "low"
	Default      = "default"
	High         = "high"
)
