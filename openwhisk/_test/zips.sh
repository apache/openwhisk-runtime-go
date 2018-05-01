#!/bin/bash
cd "$(dirname $0)"
rm action.zip 2>/dev/null
zip -r -q action.zip action
