# notify-done-to-slack

`notify-done-to-slack` is a CLI tool that notifies the completion of command execution to Slack.

## Installation

To install this project, follow these steps:

```sh
git clone https://github.com/kashu-02/notify-done-to-slack.git
cd notify-done-to-slack
go mod tidy
```

## Usage
### Configure Webhook URL
First, you need to configure the Slack Webhook URL. Run the following command to set it up:
```
notify-done-to-slack configure --webhook-url YOUR_SLACK_WEBHOOK_URL
```

### Notify Command Completion
To notify the completion of a command execution to Slack, use the tool as follows:
```
notify-done-to-slack your-command with any args
```

You can also specify the number of lines to include from the start or end of the output:
```
notify-done-to-slack --head 10 your-command with any args
notify-done-to-slack -e 10 your-command with any args

notify-done-to-slack --tail 10 your-command with any args
notify-done-to-slack -t 10 your-command with any args
```

## License
This project is licensed under the MIT License. See the LICENSE file for details.
