package environment

var rootPath = "/var/lib/" + getAgentName()
var configFilePath = "/etc/" + getAgentName() + "/" + getAgentName() + ".conf"
