---
version: "2.0"

services:
  web:
    image: devenbhooshan/starport-web:a0793d14fa330478ea8d3699811b089b2aafed41
    expose:
        - port: 8080
          to:
            - global: true

profiles:
  compute:
    web:
      resources:
        cpu:
          units: 0.5
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
        web:
          denom: uakt
          amount: 100

deployment:
  web:
    dcloud:
      profile: web
      count: 1