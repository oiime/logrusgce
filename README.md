# logrusgce

A simple formatter for logrus that logs in gce format, a rework of `https://github.com/znly/logrus-gcp` as it's quite out of date


    go get github.com/oiime/logrusgce


## Usage

```golang
log := logrus.New()
log.SetFormatter(logrusgce.NewGCEFormatter())
```

## Usage with postprocessor (allowing adding additional fields)
```golang
log := logrus.New()
log.SetFormatter(logrusgce.NewGCEFormatter().WithPostProcess(logrusgce.PostprocessHttpRequest("http-request")))

```

Look at PostprocessHttpRequest for an example on how to add postprocessors

