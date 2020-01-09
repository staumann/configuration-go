# configuration-go

## What is this?
This is a simple configuration framework for yaml files

##How to use

###Init
````go
    config.Init("development","/",WARN)
````
###Get String
````go
    config.GetString("fancy.key")
    config.GetStringWithDefault("fancy.key","defaultValue")
````
###Get boolean
````go
    config.GetBoolean("fancy.boolean.key")
    config.GetBooleanWithDefault("fancy.boolean.key", false)
````
###Get integer
````go
    config.GetInteger("fancy.integer.key")
    config.GetIntegerWithDefault("fancy.integer.key", 55)
````
