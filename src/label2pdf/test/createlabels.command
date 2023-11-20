#!/bin/sh

DIRNAME=`dirname "$0"`

pushd $DIRNAME

./label2pdf test/page.json test/label.json test/out.pdf

popd
