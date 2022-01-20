import json

config = {}
with open("schema_config.json") as f:
  config = json.load(f)

props = {}
for key in config:
  props[key] = {
    "type": "string",
    "description": config[key],
  }
out = {
  "$schema": "http://json-schema.org/schema",
  "type": "object",
  "properties": props,
}

with open("schema.json", "w+") as f:
  json.dump(out, f, indent=2)