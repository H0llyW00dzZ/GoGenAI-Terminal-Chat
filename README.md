# GoGenAI-Terminal-Chat

[![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat)](https://goreportcard.com/report/github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat)
[![Go Reference](https://pkg.go.dev/badge/github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat.svg)](https://pkg.go.dev/github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat)

**Note:** This repository is a work in progress (WIP).

**Estimated Time of Arrival (ETA):** **Unknown**. *This project is developed on a personal basis during my free time and is not associated with any company or enterprise endeavors.*

<p align="center">
  <img src="https://i.imgur.com/j4OXbC8.png" alt="Go Picture" />
</p>

## 🤔 Why Continue Development?

Interesting to built it in terminal after touring journey through the Go programming language, plus exploring Google's AI capabilities (currently in beta with Gemini).

> [!NOTE]
> This repository is specifically designed to adhere to idiomatic Go principles. 🤪

***Fun Fact***: ***Did You Know?*** If your `Go` code resembles a jungle of `if` statements (think 10+ nested layers – a big no-no!), it's less `Go` and more `Stop-and-ask-for-directions`. Flatten those conditionals and let your code run as smoothly as a greased `gopher` on a slip 'n slide! 🤪

## ✨ Current Features

- **Terminal-Based Interaction**: Experience AI chatting within the comfort of your terminal, with a focus on simplicity and efficiency.
- **Session Chat History**: Maintain a transcript of your dialogue, capturing both queries and AI responses for a continuous conversational flow.
- **Intelligent Shutdown**: Benefit from built-in signal handling for a smooth exit process, ensuring your session ends without disruption and with proper resource cleanup.
- **Realistic Typing Animation**: Enjoy a more lifelike interaction with a simulated typing effect that mimics human conversation timing.
> [!NOTE]
> This `Realistic Typing Animation` specialized feature is economical in terms of resource consumption (e.g, `memory,cpu`), in contrast to `front-end languages` or `other languages` that tend to be more resource-intensive.
- **Ease of Deployment**: Quickly deploy the chat application using Docker, minimizing setup time and complexity.
- **Command Handling**: Integrate special chat commands naturally within the conversation. The `:quit` command, for example, allows users to end their session in an orderly fashion. When this command is executed, it prompts a cooperative shutdown sequence with the AI, which generates an appropriate goodbye message. This thoughtful design enhances the user experience by providing a conversational closure that is both natural and polite, ensuring that the session termination is as engaging as the conversation itself.
- **Concurrency**: Leverage the power of Go's concurrency model with goroutines.
- **Minimalist Package**: `DebugOrErrorLogger`
The `DebugOrErrorLogger` package offers a streamlined and efficient logging system, designed specifically for Go applications that require robust error and debug logging capabilities with minimal overhead.

### 🛠️ 📦 Minimalist Package `DebugOrErrorLogger` Key Features

- 🔎 **Conditional Debug Logging (Currently Marked as TODO 📝)**: The logger allows for debug messages to be conditionally output based on the `DEBUG_MODE` environment variable. When set to `true`, detailed debug information will be printed to `os.Stderr`, aiding in the development and troubleshooting process.

- 🎨 **Color-Coded Error Output**: Errors are distinctly colorized in red when logged, making them stand out in the terminal for immediate attention. This colorization helps in quickly identifying errors amidst other log outputs.

- 😱 🔋 **Panic Recovery**: A recovery function is provided to gracefully handle and log any panics that may occur during runtime. This function ensures that a panic message is clearly logged with colorized output, preventing the application from crashing unexpectedly and aiding in rapid diagnosis.

- ⚡ **Simple API**: The package exposes a simple and intuitive API, with methods for debug and error logging that accept format strings and variadic arguments, similar to the standard `Printf` and `Println` functions.

- 🔓 ⚙️ **Environment Variable Configuration**: The debug mode can be easily toggled on or off through an environment variable, allowing for flexible configuration without the need to recompile the application.


## ⚖️ Enhanced Scalability

- **Modular Command Framework**: The system now employs a modular command framework that maps commands to dedicated handler functions. This design simplifies the addition of new commands and the modification of existing ones, reducing complexity and improving maintainability. It's a scalable approach that can easily grow with the application's functionality, accommodating an expanding set of features without cluttering the command processing logic.
- **Streamlined Codebase**: By decoupling command logic from the main interaction flow, the codebase becomes more organized and easier to navigate. Developers can quickly identify where to add new command handlers and understand how commands are processed, leading to a more developer-friendly experience.
- **Flexible Command Extension**: The new structure allows for the easy integration of additional commands as the system evolves. Whether it's implementing administrative controls, user preferences, or new AI features, the command framework is designed to handle growth efficiently and logically.

> [!NOTE]
> This specialized feature, better than code resembles a jungle of `if if if if` statements has been successfully integrated.

By adopting this scalable command handling system, the chat application is well-positioned to evolve alongside advancements in AI and user expectations, ensuring a robust and future-proof user experience.

> [!NOTE]  
> Subject planning to continuously improve and add features, enhancing functionality without adding unnecessary complexity. Stay tuned for updates!

**`Go`** is designed to be straightforward and efficient, avoiding the unnecessary complexities (fuck complexities, this is `go` anti complexities) often encountered in other programming languages.

## ⚙️ Installation

To use GoGenAI Terminal Chat, you need to have Docker installed on your machine. If you don't have Docker installed, please follow the [official Docker installation guide](https://docs.docker.com/get-docker/).

Once Docker is set up, you can pull the image from GitHub Packages by running:

```sh
docker pull ghcr.io/h0llyw00dzz/gogenai-terminal-chat:latest
```

## Usage

To start a chat session with GoGenAI, run the following command in your terminal. Make sure to replace `YOUR_API_KEY` with the actual API key provided to you.

```sh
docker run -it --rm --name mychatapp -e API_KEY=YOUR_API_KEY ghcr.io/h0llyw00dzz/gogenai-terminal-chat:latest
```

This command will start the GoGenAI Terminal Chat application in interactive mode. You will be able to type your messages and receive responses from the AI.

### Environment Variables

- `API_KEY`: Your API key for accessing the generative AI model. This is a required environment variable.


## 📸 Screenshot

<p align="center">
  <img src="https://i.imgur.com/H9vDYjZ.png" alt="kubectl terminal" />
</p>

### 📊 Average Consumption (Tested on Production Docker):

<p align="center">
  <img src="https://i.imgur.com/3haszmI.png" alt="docker" />
</p>

#### Idle

<p align="center">
  <img src="https://i.imgur.com/nhcFuXJ.png" alt="docker" />
</p>

#### Gopher doing task

<p align="center">
  <img src="https://i.imgur.com/BGtQar4.png" alt="docker" />
</p>

> [!NOTE]
> The `Average Consumption` metrics are calculated without including the use of a storage system like a database and are based on the assumption that each function is relatively simple, with an average cyclomatic complexity of 5 as the maximum. However, consumption may increase with more complex functions (e.g., those with a cyclomatic complexity of 10 or more, which are not recommended).

### 🤖 Better LLM's:
<p align="center">
  <img src="https://i.imgur.com/ge884rb.png" alt="terminal" />
</p>

<p align="center">
  <img src="https://i.imgur.com/k8JTX64.png" alt="terminal" />
</p>

### Kubernetes in the wild's:
<p align="center">
  <img src="https://i.imgur.com/4iFId5w.png" alt="k8s" />
  <img src="https://i.imgur.com/qVEpOZy.png" alt="k8s" />
  <img src="https://i.imgur.com/Ox7CqMV.png" alt="k8s" />
  <img src="https://i.imgur.com/BH81O61.png" alt="k8s" />
</p>

## Example how it work

🤓 You: hello are you ?

🤖 AI: : I am a large language model, trained by Google.

🤓 You: I have question about k8s, How do I enable audit logging in Kubernetes, and what should I look for when analyzing the logs?

🤖 AI: To enable audit logging in Kubernetes, you can use the following steps:

1. Edit the Kubernetes API server configuration file, typically located at `/etc/kubernetes/manifests/kube-apiserver.yaml`.
2. Add the following configuration under the `audit-policy` section:

```
enabled: true
```

3. Optionally, you can configure the following settings:

   - `policyFile`: The path to a file containing the audit policy rules.
   - `logPath`: The path to the file where the audit logs will be stored.
   - `logFormat`: The format of the audit logs.
   - `maxBatchSize`: The maximum number of audit events to include in a single batch.
   - `maxBatchWait`: The maximum amount of time to wait for a batch of audit events to be filled before writing them to disk.

4. Save and restart the Kubernetes API server.

Once audit logging is enabled, you can analyze the logs to gain insights into the activity on your Kubernetes cluster. Some things you may want to look for include:

- Successful and failed authentication and authorization attempts
- Create, update, and delete operations on Kubernetes resources
- Changes to cluster configuration
- API calls that take a long time to complete
- Errors and warnings generated by the Kubernetes API server

By analyzing the audit logs, you can identify security risks, troubleshoot issues, and gain a better understanding of how your Kubernetes cluster is being used.

Here are some additional resources that you may find helpful:

- [Kubernetes Audit Logging Documentation](https://kubernetes.io/docs/tasks/debug-application-cluster/audit/)

---

### ⚡ Another Example

<p align="center">
  <img src="https://i.imgur.com/PxjZ0Dz.png" alt="gopher" />
</p>

```sh
$ go run cmd/main.go
```

🤖 AI: Hello! How can I assist you today?

🤓 You: I have question does this correct about "Strive for simplicity in each function (aim for a cyclomatic complexity under 5 to demonstrate expertise); simplicity is a sign of Go programming mastery." ?

🤖 AI: Yes, the statement "Strive for simplicity in each function (aim for a cyclomatic complexity under 5 to demonstrate expertise); simplicity is a sign of Go programming mastery" is generally correct.

**Cyclomatic complexity** is a measure of the complexity of a function. It is calculated by counting the number of independent paths through the function. A high cyclomatic complexity indicates that the function is difficult to understand and maintain.

**Simplicity** is a key principle of Go programming. Go code is typically concise and easy to read. This makes it easier to understand and maintain.

**Aiming for a cyclomatic complexity under 5** is a good way to ensure that your functions are simple and easy to understand. This will help you to write better code and avoid introducing bugs.

Of course, there are times when it is necessary to write functions that are more complex. However, you should always strive for simplicity whenever possible.

Here are some tips for writing simple functions:

* Use short, descriptive variable names.
* Avoid nesting loops and conditionals.
* Use functions to break down complex tasks into smaller, more manageable tasks.
* Document your functions clearly.

By following these tips, you can write simple, easy-to-understand functions that will make your code more maintainable and bug-free.

---

🤓 You: alright thx for info

🤖 AI: You're welcome! I'm glad I could be of assistance. If you have any more questions, feel free to ask.

---

🤓 You: :quit

🤖 AI: It's been a pleasure assisting you today. Remember, I'm always here if you have any more questions. Take care and have a wonderful day!

---

## 🙌 Contributing
Contributions is always welcome! If you're able to write code that adheres to a cyclomatic complexity of 5 or less per function, please feel free to submit a pull request. Pull requests that do not meet this criterion will be rejected. Alternatively, if you encounter any issues, do not hesitate to submit an issue report.

## 📝 TODO

### Front-End:
- [ ] **UI Front-End Terminal**
- [ ] **Reporting System**

> [!NOTE]
> The `Reporting System` is designed to capture and handle runtime panic events in the Go application, facilitating streamlined error reporting and analysis.

## 🕵️‍♂️ HACKING
Available [`here`](SECURITY.md)
