FROM cypress/base:12

WORKDIR /e2e

COPY package.json cypress.json ./

RUN npm i

RUN $(npm bin)/cypress verify

COPY ./wait /wait

RUN ["chmod", "+x", "/wait"]

CMD /wait && $(npm bin)/cypress run
