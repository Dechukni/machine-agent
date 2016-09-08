var expect = chai.expect;
var client;
var channel;

chai.use(chaiAsPromised);

describe('Given an instance of my ExecClientChannel', function () {
  before(function () {
    if (typeof (client) === 'undefined') {
      client = new ExecClient(execAgentUrl);
    }
    channel = client.channel();
  });
  describe('when I need a test that will pass 100%', function () {
    it('2 plus 2 should be equal to 4', function () {
      expect(2 + 2).to.equal(4);
    });
  });
  describe('when I need to start process', function () {
    it('should start process and return its name', function () {
      var resultName;

      resultName = channel.process.start('ping', 'ping ya.ru')
        .then(function (res) {
          return res.name;
        });
      return expect(resultName).to.eventually.equal('ping');
    });
  });
  describe('when I need to start process and then kill it', function () {
    it('should start process by providing object as parameter, return pid of it and then kill it', function () {
      var obj = {name: 'ping', commandLine: 'ping google.com', type: 'YAFP'};
      var resultDesc;
      var pid;

      resultDesc = channel.process.start(obj)
        .then(function (res) {
          pid = res.pid;
          return channel.process.kill(pid);
        })
        .then(function (res) {
          return res.text;
        });
      return expect(resultDesc).to.eventually.equal('Successfully killed');
    });
  });
  describe('when I need to start process and then obtain its logs', function () {
    it('should start process, and return logs skipping first 5 of them', function () {
      var getLogsParam = {'pid': 0, 'skip': 5};
      var result;

      result = channel.process.start('ping', 'ping github.com')
        .then(function (res) {
          getLogsParam.pid = res.pid;
          return channel.process.getLogs(getLogsParam);
        });
      return expect(result).to.be.fulfilled;
    });
  });
});

