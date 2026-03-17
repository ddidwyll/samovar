## Project: "Samowar"

### Generated with
 - Types for the network messaging: false
 - Enabled Observer (http://localhost:9911): true
 - Loggers: colored
 

### Supervision Tree

Applications
 - `Device{}` samowar/apps/device/device.go
   - `DeviceSup{}` samowar/apps/device/devicesup.go
     - `DeviceActor{}` samowar/apps/device/deviceactor.go
 - `Mqtt{}` samowar/apps/mqtt/mqtt.go
   - `MqttSup{}` samowar/apps/mqtt/mqttsup.go
     - `MqttActor{}` samowar/apps/mqtt/mqttactor.go
 - `Script{}` samowar/apps/script/script.go
   - `ScriptSup{}` samowar/apps/script/scriptsup.go
     - `ScriptActor{}` samowar/apps/script/scriptactor.go
 - `Client{}` samowar/apps/client/client.go
   - `ClientSup{}` samowar/apps/client/clientsup.go
     - `ClientActor{}` samowar/apps/client/clientactor.go
     - `ClientService{}` samowar/apps/client/clientservice.go


#### Used command

This project has been generated with the `ergo` tool. To install this tool, use the following command:

`$ go install ergo.services/tools/ergo@latest`

Below the command that was used to generate this project:

```$ /home/one/go/bin/ergo -init Samowar{} -with-app Device -with-sup Device:DeviceSup -with-actor DeviceSup:DeviceActor -with-app Mqtt -with-sup Mqtt:MqttSup -with-actor MqttSup:MqttActor -with-app Script -with-sup Script:ScriptSup -with-actor ScriptSup:ScriptActor -with-app Client -with-sup Client:ClientSup -with-actor ClientSup:ClientActor -with-logger colored -with-web "ClientSup:ClientService{host:localhost,port:4000}" -with-observer ```
