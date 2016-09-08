'use strict';

import WebsocketAgent from './WebsocketAgent.js';
import utils from './utils';

export default class ExecClientChannel {
  constructor(address = 'ws://localhost:9000/connect') {
    const websocket = new WebsocketAgent(address);

    this.process = {
      start: function (nameObj, cmdLine) {
        let data = {operation: 'process.start', id: utils.generateId()};

        if (cmdLine !== undefined) {
          data.body = {name: nameObj, commandLine: cmdLine};
          return websocket.invokeOperationCall(data);
        }
        if (nameObj.handler !== undefined) {
          websocket.addTempEventHandler(data.id, nameObj.handler);
          delete nameObj.handler;
        }
        data.body = nameObj;
        return websocket.invokeOperationCall(data);
      },
      kill: function (pid) {
        let data = {operation: 'process.kill', id: utils.generateId()};

        data.body = {pid: pid};
        return websocket.invokeOperationCall(data);
      },
      subscribe: function (pidObj) {
        let data = {operation: 'process.subscribe', id: utils.generateId()};

        if (utils.isInt(pidObj)) {
          data.body = {pid: pidObj};
          return websocket.invokeOperationCall(data);
        }
        if (pidObj.handler !== undefined) {
          websocket.addEventHandler(pidObj.pid, pidObj.handler);
          delete pidObj.handler;
        }
        data.body = pidObj;
        return websocket.invokeOperationCall(data);
      },
      unsubscribe: function (pid) {
        let data = {operation: 'process.unsubscribe', id: utils.generateId()};

        data.body = {pid: pid};
        return websocket.invokeOperationCall(data);
      },
      upsubscribe: function (pidObj, eventTypes) {
        let data = {operation: 'process.updateSubscriber', id: utils.generateId()};

        if (eventTypes !== undefined) {
          data.body = {pid: pidObj, eventTypes: eventTypes};
          return websocket.invokeOperationCall(data);
        }
        if (pidObj.handler !== undefined) {
          websocket.addEventHandler(pidObj.pid, pidObj.handler);
          delete pidObj.handler;
        }
        data.body = pidObj;
        return websocket.invokeOperationCall(data);
      },
      getLogs: function (pidObj) {
        let data = {operation: 'process.getLogs', id: utils.generateId()};

        if (utils.isInt(pidObj)) {
          data.body = {pid: pidObj};
          return websocket.invokeOperationCall(data);
        }
        data.body = pidObj;
        return websocket.invokeOperationCall(data);
      }
    };
    this.all = {
      startEventHandling: function () {
        websocket.shouldHandleEvents = true;
      },
      stopEventHandling: function () {
        websocket.shouldHandleEvents = false;
      }
    };
  }
}
