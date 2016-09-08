'use strict';

import RESTcaller from './RESTcaller.js';
import ExecClientChannel from './ExecClientChannel.js';
import utils from './utils.js';

export default class ExecClient {
  constructor(address = 'localhost', port = '9000', saveHistory = true) {
    this.address = `http://${address}:${port}`;
    const caller = new RESTcaller(this.address);

    this.process = {
      start: function (nameObj, cmdLine) {
        if (cmdLine !== undefined) {
          return caller.start({name: nameObj, commandLine: cmdLine});
        }
        return caller.start(nameObj);
      },
      kill: function (pid) {
        return caller.kill(pid);
      },
      info: function (pid) {
        return caller.info(pid);
      },
      getLogs: function (pidObj) {
        if (utils.isInt(pidObj)) {
          return caller.getLogs(pidObj, {});
        }
        let pidReal = pidObj.pid;

        delete pidObj.pid;
        return caller.getLogs(pidReal, pidObj);
      },
      subscribe: function (pidObj, channel) {
        if (channel !== undefined) {
          return caller.subscribe({pid: pidObj, channel: channel});
        }
        let pidChannelReal = {pid: pidObj.pid, channel: pidObj.channel};

        delete pidObj.pid;
        delete pidObj.channel;
        return caller.subscribe(pidChannelReal, pidObj);
      },
      unsubscribe: function (pid, channel) {
        return caller.unsubscribe(pid, channel);
      },
      upsubscribe: function (pid, channel, types) {
        return caller.upsubscribe(pid, channel, types);
      }
    };
    this.all = {
      info: function (all) {
        return caller.allInfo(all);
      }
    };
  }
  channel() {
    let channelAddress = `${this.address.replace('http', 'ws')}/connect`;

    return new ExecClientChannel(channelAddress);
  }
}
