import * as request from 'superagent';

export default class RESTcaller {
  constructor(url) {
    this.url = url;
  }
  start(data) {
    return new Promise(
      (resolve, reject) => {
        request
          .post(`${this.url}/process`)
          .send(data)
          .set('Accept', 'application/json')
          .end(function (err, res) {
            if (err || !res.ok) {
              let codeDescription;

              if (err.status === 400) codeDescription = 'Incoming data is not valid';
              if (err.status === 404) codeDescription = 'There is no such channel';
              if (err.status === 500) codeDescription = 'Something went wrong in ExecAgent';
              let result = {
                code: err.status,
                codeDescription: codeDescription
              };

              reject(result);
            } else {
              let result = {
                code: 200,
                codeDescription: 'Successfully started process',
                body: res.body
              };

              resolve(result);
            }
          });
      });
  }
  kill(pid) {
    return new Promise(
      (resolve, reject) => {
        request
          .del(`${this.url}/process/${pid}`)
          .set('Accept', 'application/json')
          .end(function (err, res) {
            if (err || !res.ok) {
              let codeDescription;

              if (err.status === 400) codeDescription = 'Incoming pid is not valid';
              if (err.status === 404) codeDescription = 'There is no such process';
              if (err.status === 500) codeDescription = 'Something went wrong in ExecAgent';
              let result = {
                code: err.status,
                codeDescription: codeDescription
              };

              reject(result);
            } else {
              let result = {
                code: 200,
                codeDescription: 'Successfully killed process',
                body: res.body
              };

              resolve(result);
            }
          });
      });
  }
  info(pid) {
    return new Promise(
      (resolve, reject) => {
        request
          .get(`${this.url}/process/${pid}`)
          .set('Accept', 'application/json')
          .end(function (err, res) {
            if (err || !res.ok) {
              let codeDescription;

              if (err.status === 400) codeDescription = 'Incoming pid is not valid';
              if (err.status === 404) codeDescription = 'There is no such process';
              if (err.status === 500) codeDescription = 'Something went wrong in ExecAgent';
              let result = {
                code: err.status,
                codeDescription: codeDescription
              };

              reject(result);
            } else {
              let result = {
                code: 200,
                codeDescription: 'Successfully obtained info about process',
                body: res.body
              };

              resolve(result);
            }
          });
      });
  }
  getLogs(pid, query) {
    return new Promise(
      (resolve, reject) => {
        request
          .get(`${this.url}/process/${pid}/logs`)
          .query(query)
          .set('Accept', 'application/json')
          .end(function (err, res) {
            if (err || !res.ok) {
              let codeDescription;

              if (err.status === 400) codeDescription = 'Incoming pid is not valid';
              if (err.status === 404) codeDescription = 'There is no such process';
              if (err.status === 500) codeDescription = 'Something went wrong in ExecAgent';
              let result = {
                code: err.status,
                codeDescription: codeDescription
              };

              reject(result);
            } else {
              let result = {
                code: 200,
                codeDescription: 'Successfully obtained logs of process',
                body: res.body
              };

              resolve(result);
            }
          });
      });
  }
  subscribe(pidChannel, query) {
    return new Promise(
      (resolve, reject) => {
        request
          .post(`${this.url}/process/${pidChannel.pid}/events/${pidChannel.channel}`)
          .query(query)
          .set('Accept', 'application/json')
          .end(function (err, res) {
            if (err || !res.ok) {
              let codeDescription;

              if (err.status === 400) codeDescription = 'Incoming parameter(s) are not valid';
              if (err.status === 404) codeDescription = 'There is no such process or channel';
              if (err.status === 500) codeDescription = 'Something went wrong in ExecAgent';
              let result = {
                code: err.status,
                codeDescription: codeDescription
              };

              reject(result);
            } else {
              let result = {
                code: 200,
                codeDescription: 'Successfully subscribed channel to process events',
                body: res.body
              };

              resolve(result);
            }
          });
      });
  }
  unsubscribe(pid, channel) {
    return new Promise(
      (resolve, reject) => {
        request
          .delete(`${this.url}/process/${pid}/events/${channel}`)
          .set('Accept', 'application/json')
          .end(function (err, res) {
            if (err || !res.ok) {
              let codeDescription;

              if (err.status === 400) codeDescription = 'Incoming parameter(s) are not valid';
              if (err.status === 404) codeDescription = 'There is no such process or channel';
              if (err.status === 500) codeDescription = 'Something went wrong in ExecAgent';
              let result = {
                code: err.status,
                codeDescription: codeDescription
              };

              reject(result);
            } else {
              let result = {
                code: 200,
                codeDescription: 'Successfully unsubscribed channel to process events',
                body: res.body
              };

              resolve(result);
            }
          });
      });
  }
  upsubscribe(pid, channel, types) {
    return new Promise(
      (resolve, reject) => {
        request
          .put(`${this.url}/process/${pid}/events/${channel}`)
          .query(types)
          .set('Accept', 'application/json')
          .end(function (err, res) {
            if (err || !res.ok) {
              let codeDescription;

              if (err.status === 400) codeDescription = 'Incoming parameter(s) are not valid';
              if (err.status === 404) codeDescription = 'There is no such process or channel';
              if (err.status === 500) codeDescription = 'Something went wrong in ExecAgent';
              let result = {
                code: err.status,
                codeDescription: codeDescription
              };

              reject(result);
            } else {
              let result = {
                code: 200,
                codeDescription: 'Successfully updated subscription of a channel to process events',
                body: res.body
              };

              resolve(result);
            }
          });
      });
  }
  allInfo(all) {
    return new Promise(
      (resolve, reject) => {
        request
          .get(`${this.url}/process`)
          .query(all)
          .set('Accept', 'application/json')
          .end(function (err, res) {
            if (err || !res.ok) {
              let codeDescription;

              if (err.status === 500) codeDescription = 'Something went wrong in ExecAgent';
              let result = {
                code: err.status,
                codeDescription: codeDescription
              };

              reject(result);
            } else {
              let result = {
                code: 200,
                codeDescription: 'Successfully obtained info about processes',
                body: res.body
              };

              resolve(result);
            }
          });
      });
  }
}
