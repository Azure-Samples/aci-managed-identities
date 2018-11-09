FROM node:11.1.0
WORKDIR  /
COPY ./sample.js /
RUN npm install request-promise
RUN npm install request
RUN npm install azure-keyvault

CMD ["sh", "-c", "node sample.js"]
