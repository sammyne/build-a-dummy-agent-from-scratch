#!/bin/bash

API_KEY=AIzaSyDqVWgayGXe5eSQiF6d8N3kW2Bmwj3BcwI

curl "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $API_KEY" \
-d @request.json
