FROM cypress/included:13.8.0

WORKDIR /cypress

RUN npm init -y && npm install cypress-wait-until@3.0.1 @cypress/xpath@2.0.3 @4tw/cypress-drag-drop@2.2.5

#COPY cypress.config.js ./cypress/
#COPY plugins ./cypress/plugins
#COPY support ./cypress/support
#COPY e2e ./cypress/e2e

CMD ["cypress", "run", "--config-file", "cypress/cypress.config.js", "--browser", "chrome"]
