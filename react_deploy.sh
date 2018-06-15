#!/bin/bash
cd web
rm -rf *
cd ../react_app/my-app
npm run build
cd build
mv * ~/go/src/github.com/louiscarteron/WebApps2018/web/
