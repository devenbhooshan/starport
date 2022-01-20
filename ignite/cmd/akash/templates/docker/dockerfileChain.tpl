FROM starport/cli

ENTRYPOINT []

WORKDIR /apps

USER root
RUN chown -R root:root /apps

COPY . .

EXPOSE 26656 1317 26657

RUN starport chain build --output .
RUN starport chain init
RUN ls -lha

CMD [ "/apps/hellod", "start"]