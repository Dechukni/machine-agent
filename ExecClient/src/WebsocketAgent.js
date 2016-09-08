'use strict';

import * as ws from 'websocket';

export default class WebsocketAgent {
  constructor(url) {
    this.handlers = [];
    this.tempHandler = [];
    this.shouldHandleEvents = true;

    this.promiseResolvers = [];
    this.promiseRejecters = [];

    this.socket = new ws.w3cwebsocket(url);
    this.socket.onmessage = (event) => {
      let parsedEvent = JSON.parse(event.data);

      if (parsedEvent.id !== undefined) {
        if (parsedEvent.body !== null) {
          if (parsedEvent.body.commandLine !== undefined && this.tempHandler[parsedEvent.id] !== undefined) {
            this.addEventHandler(parsedEvent.body.pid, this.tempHandler[parsedEvent.id]);
          }
          this.promiseResolvers[parsedEvent.id](parsedEvent.body);
        } else {
          this.promiseRejecters[parsedEvent.id](parsedEvent.error);
        }
      } else {
        if (this.shouldHandleEvents) {
          if (this.handlers[parsedEvent.body.pid] !== undefined) {
            for (let handler of this.handlers[parsedEvent.body.pid]) {
              handler(parsedEvent);
            }
          }
        }
      }
    };
  }
  waitUntilOpen() {
    return new Promise(
      (resolve) => {
        if (this.socket.readyState === 1) {
          resolve();
        } else {
          this.socket.onopen = () => {
            resolve();
          };
        }
      }
    );
  }
  invokeOperationCall(data) {
    return this.waitUntilOpen()
      .then(() => {
        this.socket.send(JSON.stringify(data));
        return new Promise(
          (resolve, reject) => {
            this.promiseResolvers[data.id] = resolve;
            this.promiseRejecters[data.id] = reject;
          });
      });
  }
  addEventHandler(pid, handler) {
    if (this.handlers[pid] === undefined) {
      this.handlers[pid] = [];
    }
    this.handlers[pid].push(handler);
  }
  addTempEventHandler(id, handler) {
    this.tempHandler[id] = handler;
  }
}
