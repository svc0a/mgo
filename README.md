# Purpose of project

the purpose of this project is to generate fields map from struct.

# Supported database

## mongodb(default)

```
	err := Define().Generate()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("success")
```

## postgreSQL

```
	err := Define(WithPostgre()).Generate()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("success")
```

# set dir

```
	err := Define(WithDir("../")).Generate()
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("success")
```