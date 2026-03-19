## Project: "Samovar"

### Generated with
 - Types for the network messaging: false
 - Enabled Observer (http://localhost:9911): true
 - Loggers: colored

#### Used command

This project has been generated with the `ergo` tool. To install this tool, use the following command:

`$ go install ergo.services/tools/ergo@latest`

Below the command that was used to generate this project:

```$ /home/one/go/bin/ergo -init Samovar{} -with-app Api -with-sup Api:Sup -with-actor Sup:Actor -with-logger colored -with-web "Sup:Web{host:localhost,port:4000}" -with-observer ```
