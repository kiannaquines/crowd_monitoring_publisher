## Crowd Monitoring MQTT Publisher

### Note
You must already have your wifi adapter install and also must be in monitor mode


### Guide
1. Rename the `env` file to `.env`
2. Run `go mod tidy`
3. Fill up the necesarry fields in `.env` you just renamed
4. Compile the go script to binary `go build -o <your_executable_name> main.go`
5. Run the executable using `sudo ./<your_executable_name>`

## Output

This should be the output of the script after successfully running it

`2024/09/24 20:27:01 Connected to MQTT broker`
`2024/09/24 20:27:01 Crowd Monitoring Publisher is starting at interface wlan1`

## Any bugs? 
Please submit a ticket if there is an issue with the publisher
