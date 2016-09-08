#ExecAgent#

**Caution: This is WIP.** Not all features described in documentation are implemented completely, and not all implemented features (currently main API) are guaranteed to work error-proof.

##Docs and links##
-	[ExecAgent] (https://github.com/evoevodin/machine-agent)
-	[ExecClientAPI] (docs/ExecClientAPI.md)
-	[ExecClientChannelAPI] (docs/ExecClientChannelAPI.md)

##Usage##
To use library - clone the repository and run **npm install**. This will install all dependencies of the library. Then to produce minified version of library run **npm run build**. The code is written in native ES6 (verified with [eslint] (http://eslint.org), it is compiled and transpiled using [Babel] (https://babeljs.io) and then combined into single library without any dependencies using [webpack] (https://webpack.github.io). Resulting universal library should work with any king of importing into both browser and Node.js (not tested for Node.js yet).

##Testing##
Testing is done using [Mocha] (https://mochajs.org)+[Chai] (http://chaijs.com)+[Chai-as-promised] (https://github.com/domenic/chai-as-promised). To test first produce normal version of library by running **npm run dev** (which is also used during development, this command automatically compile and bundle library on any changes in source code). Then in test/testrunner.html change execAgentUrl variable to address of your own ExecAgent, and open testrunner.html in any browser. If you have problems with CORS, you should use Google Chrome/Chromium and disable CORS policy by running executable with parameters --disable-web-security --user-data-dir=(path to folder where all data for unsecure Chrome session will be stored).

##License##
License file is coming soon