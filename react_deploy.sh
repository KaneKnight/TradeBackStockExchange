#!/bin/bash
cd $GOPATH/src/github.com/louiscarteron/WebApps2018/web
rm -rf *
cd $GOPATH/src/github.com/louiscarteron/WebApps2018/react_app/my-app
npm run build
cd $GOPATH/src/github.com/louiscarteron/WebApps2018/react_app/my-app/build
mv * $GOPATH/src/github.com/louiscarteron/WebApps2018/web/
