---
version: "2.0"

services:
  app:
    image: devenbhooshan/starport-app:4f13404ac06dd2ff77f913373078391de4e574ec
    expose:
        - port: 26656
          to:
            - global: true
        - port: 26657
          to:
            - global: true
        - port: 1317
          to:
            - global: true
        - port: 9090
          to:
            - global: true

profiles:
  compute:
    app:
      resources:
        cpu:
          units: 1
        memory:
          size: 512Mi
        storage:
          size: 512Mi
  placement:
    dcloud:
      attributes:
        host: akash
      signedBy:
        anyOf:
          - "akash1365yvmc4s7awdyj3n2sav7xfx76adc6dnmlx63"
      pricing:
        app:
          denom: uakt
          amount: 100

deployment:
  app:
    dcloud:
      profile: app
      count: 1