Some tricks and tips learned from mcoding (see [video](https://www.youtube.com/watch?v=9L77QExPmI0&ab_channel=mCoding) for more details).


Logging docs: https://docs.python.org/3/library/logging.html


- Your Python app will have a tree of loggers, with the root logger at the very top of the hierarchy.
- Do not instantiate your logger directly. Use logging instead:
`python
logger = logging.getLogger(name)  # you can namespace this using '.', e.g 'myapp.subApp'
`
- Avoid calling the root logger directly (e.g `logging.info("my message")`)
- Each logger has a Level, Filters, and Handlers. Level/Filters drops logrecords for the current logger and all upstream loggers.
- Each Handler has a Level, Filters, and Formatters. A handler's Level and Filters do not drop Log Records for the upstream Loggers/Handlers.
- Common setup that will be enough for most use cases: place filters and handlers on the Root logger. A single non-root logger should be enough.

A basic config:
```python
import logging

logger = logging.getLogger("app")

logging.basicConfig(level="INFO")
```

But this only logs to stdout. Normally you want to log to a file, send an email as a side effect, etc. That's when handlers and custom configs are useful. You can setup a dictionary with custom configs

```python

# this config logs WARNING and up to stderr and DEBUG an up to a set of rotating log files
# notice the customization of log format as well
my_config = {
  "version": 1,
  "disable_existing_loggers": false,
  "formatters": {
    "simple": {
      "format": "%(levelname)s: %(message)s"
    },
    "detailed": {
      "format": "[%(levelname)s|%(module)s|L%(lineno)d] %(asctime)s: %(message)s",
      "datefmt": "%Y-%m-%dT%H:%M:%S%z"
    }
  },
  "handlers": {
    "stderr": {
      "class": "logging.StreamHandler",
      "level": "WARNING",
      "formatter": "simple",
      "stream": "ext://sys.stderr"
    },
    "file": {
      "class": "logging.handlers.RotatingFileHandler",
      "level": "DEBUG",
      "formatter": "detailed",
      "filename": "logs/my_app.log",
      "maxBytes": 10000,
      "backupCount": 3
    }
  },
  "loggers": {
    "root": {
      "level": "DEBUG",
      "handlers": [
        "stderr",
        "file"
      ]
    }
  }
}

logging.config.dictConfig(my_config)
...
```

- In prod you probably want these configs in json/yaml files.
- mCoding also talked about two neat features: 
    1. custom formatters: https://github.com/mCodingLLC/VideosSampleCode/blob/c18a3573bf112a78663f2917da1652b0d31b1b11/videos/135_modern_logging/mylogger.py#L33
    1. queued/async logging to avoid blocking I/O while your app logs to a file before replying
    to the client. https://github.com/mCodingLLC/VideosSampleCode/blob/c18a3573bf112a78663f2917da1652b0d31b1b11/videos/135_modern_logging/logging_configs/4-queued-json-stderr.json#L5-L38