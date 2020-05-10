# GraphQL / Angular branch

This branch broke off from v0.3.0 and worked at various times towards:

- API server improvements; these improvements added back to v0.3.2
  - CORS headers to allow all origins
  - Configurable address and port
- Angular UI (bedrock-ui/)
  - It has a data service, but I haven't refamiliarized myself with what it does
  - The only active component shows a "chunk map", just a black & white set of squares representing defined chunks in the world
  - It's pretty clunky to use; need to `npm install`, then `npm run ng serve` and also have the command line api running on localhost:8080
  - Eventually the app would have been bundled into the cli api
  - I don't expect to revisit this part of the project
- GraphQL API (mcpegql/)
  - It looks like I didn't get beyond get/put/delete of raw keys
  - But allowing several formats of data
  - It is working when running mcpetool as api, browse to http://127.0.0.1:8080/graphql
  - I like GraphQL, but I don't expect to continue this; currently I'm favoring using an embedded script instead of an api
- Miscellaneous - I probably want to look at these and pull or update them into a current branch
  - d3.js example - I think this is another chunk map, but it's not working for me right now
  - Some notes on experiments I did with deleting chunks
  - AngularJS example which works and seems to parse a lot of info

  I merged really just the fixed powershell example plus go.mod & go.sum from a recent branch, but no real work has been done in this branch since January 2019.