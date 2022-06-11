#!/bin/bash
k6 run --vus 10 --duration 30s test.script.js
