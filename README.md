flowchart
    MS[mqtt.server] -->|send msg| BSP(bus.mqtt.producer)
    BSP -->|mqtt_new_message| BDC([bus.device.consumer])
    BDC -->|send request| DRS[(device.raw_state)]
    BDC -->|send report| DCL[device.change_log]
    BDC -->|send report| DC{device.calc}
    DRS -->|send report| BDP(bus.device.producer)
    BDP -->|device_raw_state_changed| BDC
    BDP -->|device_state_changed| BCC([bus.client.consumer])
    DC <-->|call state| DRS
    DC -->|send requests| DS[(device.state)]
    DS -->|send reports| BDP
