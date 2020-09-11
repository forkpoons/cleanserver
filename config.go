package main

type Configs struct {
	System struct {
		DB       string
		Log      string
		LogLevel int
	}
}

func (c *Configs) GetConnString() string {
	return c.System.DB
}

func (c *Configs) GetLogWay() string {
	return "file"
}

func (c *Configs) GetLogDestination() string {
	return c.System.Log
}

func (c *Configs) GetLogApplicationName() string {
	return "App"
}
