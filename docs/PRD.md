# Product Requirements Document: ClaudePilot

## 1. Introduction

ClaudePilot is a terminal-based application designed for developers, researchers, and AI enthusiasts who need to manage and orchestrate multiple Anthropic Claude instances simultaneously. It provides a powerful Text-based User Interface (TUI) built with Golang that allows users to spawn, monitor, and facilitate complex interactions between multiple AI sessions from a single terminal window. This tool is born out of the need to move beyond single-threaded conversations with AI and enable parallelized, collaborative AI workflows.

## 2. Core Goals

* **Effortless Orchestration:** To provide a user-friendly TUI for managing the lifecycle of multiple Claude sessions, including spawning, starting, and stopping.
* **Advanced Monitoring:** To enable users to monitor the status, resource usage, and output of each Claude session in real-time.
* **Seamless Integration:** To facilitate deep interaction and data exchange between different Claude sessions, allowing them to work in concert.
* **Flexible Prompting:** To allow users to spawn new Claude sessions with specific, templated, or dynamically generated prompts.

## 3. Features

### 3.1. Session Lifecycle Management

* **Spawn New Sessions:** Users can create a new Claude session, giving it a name and an initial prompt. The session is created but not necessarily started immediately.
* **Start/Stop Sessions:** Users can individually or globally start and stop the execution of the AI logic within the sessions. A "stopped" session retains its history but is not actively processing.
* **Kill Sessions:** Users can terminate any running Claude session, freeing up resources. The session history can be optionally persisted.
* **List Sessions:** The TUI will display a list of all active and inactive Claude sessions, showing their status (e.g., `running`, `stopped`, `error`), a unique ID or name, and a summary of their current task.

### 3.2. Concurrent Session Monitoring

* **Real-time Updates:** The TUI will update in real-time to show the latest output and status of each session.
* **Resource Monitoring:** Display basic resource usage (CPU, memory) for each session's process.
* **Log Viewer:** A dedicated, scrollable view for each session's complete log history, including timestamps.
* **Status Bar:** A global status bar showing the total number of sessions, running sessions, and overall application status.

### 3.3. Inter-Session Communication & Integration

* **Piping:** Users can "pipe" the output of one Claude session directly as an input prompt to another. For example, `session-1 | session-2`.
* **Broadcasting:** A user can send a single prompt to multiple selected sessions at once.
* **Shared Context:** (Future Consideration) A mechanism for creating a shared "scratchpad" or context window that multiple sessions can read from and write to.

### 3.4. User Interaction & Prompting

* **Focused Interaction:** Users can select a specific session to interact with directly, sending prompts and viewing its output in a dedicated pane.
* **Prompt Templates:** Ability to save and load pre-defined prompts for spawning new sessions quickly.

## 4. User Interface (UI) & Experience (UX)

The TUI will be built using a modern Golang library (e.g., Bubble Tea) and will be divided into several resizeable panes:

* **Session List Pane (Left):** A vertically scrollable list of all Claude sessions. Each entry will show:
  * `[ID/Name]`
  * `[Status: Running/Stopped/Error]`
  * `[Last activity summary]`
  * The currently selected session will be highlighted.
* **Session Output Pane (Right, Top):** Displays the full, scrollable output of the currently selected session. This pane will support word wrapping and searching.
* **Input Pane (Right, Bottom):** A command-line input for users to interact with the selected session or the TUI itself. It will feature a command history.

### Key Commands / Hotkeys

* **Global:**
  * `Ctrl+C`: Quit the application.
  * `Tab` / `Shift+Tab`: Cycle between panes.
  * `?`: Show a help dialog with all keybindings.
* **Session List Pane:**
  * `n`: Create a **n**ew session.
  * `k`: **K**ill the selected session.
  * `s`: **S**tart/ **S**top the selected session.
  * `j`/`k` or `↓`/`↑`: Navigate the session list.
* **Input Pane:**
  * `Enter`: Send the prompt to the selected session.
  * `Ctrl+P`: Pipe output of selected session to another (will prompt for target).

## 5. Technical Stack

* **Language:** Golang
* **TUI Framework:** Bubble Tea (or a similar modern choice like tview).
* **Concurrency:** Go Routines and Channels for managing sessions.

## 6. Non-Functional Requirements

* **Performance:** The TUI should be responsive and handle at least 10-15 concurrent sessions without significant lag on a standard developer machine.
* **Reliability:** The tool should be stable, handle API errors gracefully (e.g., from the Claude API), and not crash unexpectedly.
* **Usability:** The TUI should be intuitive for users familiar with terminal multiplexers like `tmux` or `htop`.
* **Cross-Platform:** Should compile and run on macOS, Linux, and Windows.

## 7. Future Considerations

* **Session Snapshots:** Ability to save and load the state of all sessions and their histories.
* **Scripting:** Allow users to script interactions between sessions using a simple DSL or a scripting language (e.g., Lua) for automating complex workflows.
* **Plugin System:** Enable users to extend the functionality of the tool with custom plugins.
* **Cost Tracking:** If using official APIs, provide an estimated cost for the tokens used per session and globally.
