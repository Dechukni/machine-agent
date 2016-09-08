// import * as chai from 'chai';
// import * as chaiAsPromised from 'chaiAsPromised';
// import ExecClient from '../lib/ExecClient';

var expect = chai.expect;
var client;

chai.use(chaiAsPromised);

describe('Given an instance of my ExecClient', function () {
  before(function () {
    if (typeof (client) === 'undefined') {
      client = new ExecClient(execAgentUrl);
    }
  });
  describe('when I need a test that will pass 100%', function () {
    it('2 plus 2 should be equal to 4', function () {
      expect(2 + 2).to.equal(4);
    });
  });
  describe('when I need to start process', function () {
    it('should start process and return its name', function () {
      var resultName;

      resultName = client.process.start('ping', 'ping ya.ru')
        .then(function (res) {
          return res.body.name;
        });
      return expect(resultName).to.eventually.equal('ping');
    });
  });
  describe('when I need to start process and then kill it', function () {
    it('should start process by providing object as parameter, return pid of it and then kill it', function () {
      var obj = {name: 'ping', commandLine: 'ping google.com', type: 'YAFP'};
      var resultCode;
      var pid;

      resultCode = client.process.start(obj)
        .then(function (res) {
          pid = res.body.pid;
          return client.process.kill(pid);
        })
        .then(function (res) {
          return res.code;
        });
      return expect(resultCode).to.eventually.equal(200);
    });
  });
  describe('when I need to start process and then obtain info about it', function () {
    it('should start process by providing object as parameter, ' +
       'and return info about it including the same name', function () {
      var obj = {name: 'ping', commandLine: 'ping vk.com', type: 'YAFP'};
      var resultName;
      var pid;

      resultName = client.process.start(obj)
        .then(function (res) {
          pid = res.body.pid;
          return client.process.info(pid);
        })
        .then(function (res) {
          return res.body.name;
        });
      return expect(resultName).to.eventually.equal(obj.name);
    });
  });
  describe('when I need to obtain info about all processes', function () {
    it('should return info about all processes including our test one', function () {
      var resultElem;
      var resultElemName;

      resultElemName = client.process.start('noOneEverWillCreateProcessWithThatStupidName', 'ping ya.ru')
        .then(function () {
          return client.all.info(true);
        })
        .then(function (res) {
          resultElem = res.body.find(function (elem) {
            return (elem.name === 'noOneEverWillCreateProcessWithThatStupidName');
          });
          return resultElem.name;
        });
      return expect(resultElemName).to.eventually.equal('noOneEverWillCreateProcessWithThatStupidName');
    });
  });
  describe('when I need to start process and then obtain its logs', function () {
    it('should start process, and return logs skipping first 5 of them', function () {
      var getLogsParam = {'pid': 0, 'skip': 5};
      var resultCode;

      resultCode = client.process.start('ping', 'ping github.com')
        .then(function (res) {
          getLogsParam.pid = res.body.pid;
          return client.process.getLogs(getLogsParam);
        })
        .then(function (res) {
          return res.code;
        });
      return expect(resultCode).to.eventually.equal(200);
    });
  });
});

